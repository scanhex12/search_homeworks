package patternindex

import (
	"search/index"
	"search/lsm"
	"strconv"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

type PatternIndexConfig struct {
	index_document_lsm *lsm.Config
	index_document *index.IndexerConfig
	documents_config *lsm.Config
	ngram_coef int
}

type PatternIndex struct {
	index *index.Indexer
	documents *lsm.MergeTree
	ngram_coef int
}

func NewPatternIndex(config *PatternIndexConfig) *PatternIndex {
	return &PatternIndex{
		index : index.NewIndexer(config.index_document_lsm, *config.index_document),
		documents: lsm.NewMergeTable(config.documents_config),
		ngram_coef: config.ngram_coef,
	}
}

func (indexer *PatternIndex) InsertPrefixDocuments(text string, index_doc int) {
	lower_text := strings.ToLower(text)
	prefixes := SplitTextToListPrefixes(lower_text)
	indices := make([]int, len(prefixes))
	for i := 0; i < len(prefixes); i++ {
		indices[i] = index_doc
	}
	indexer.index.AddBatchDocument(prefixes, indices)
	indexer.documents.Insert(strconv.Itoa(index_doc), text)
}


func (indexer *PatternIndex) InsertPatternDocuments(text string, index_doc int) {
	lower_text := strings.ToLower(text)
	prefixes := SplitTextToListNGrams(lower_text, indexer.ngram_coef)
	indices := make([]int, len(prefixes))
	for i := 0; i < len(prefixes); i++ {
		indices[i] = index_doc
	}
	indexer.index.AddBatchDocument(prefixes, indices)
	indexer.documents.Insert(strconv.Itoa(index_doc), text)
}

func (indexer *PatternIndex) SearchByPrefix(prefix string, limit int) ([]string, error) {
	lower_prefix := strings.ToLower(prefix)
	indices, err := indexer.index.GetListDocuments(lower_prefix, limit)
	if err != nil {
		return []string{}, err
	}
	docs := make([]string, 0)
	for _, ind := range indices {
		doc, contains := indexer.documents.Search(strconv.Itoa(ind))
		if !contains {
			continue
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (indexer *PatternIndex) SearchByPattern(pattern string, limit int) ([]string, error) {
	lower_pattern := strings.ToLower(pattern)

	important_words := GetImportantSubwordFromPattern(lower_pattern, limit)

	bitmap := roaring.NewBitmap()
	for _, word := range important_words {
		cur_bitmap, err := indexer.index.GetMergedBitmapDocuments(word, limit)
		if err != nil {
			return []string{}, err
		}
		bitmap.Or(cur_bitmap)
	}

	docs := make([]string, 0)
	for _, ind := range bitmap.ToArray() {
		doc, contains := indexer.documents.Search(strconv.Itoa(int(ind)))
		if !contains {
			continue
		}
		if !MatchPatternToText(doc, pattern) {
			continue
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

