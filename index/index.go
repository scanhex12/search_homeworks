package index

import (
	"search/lsm"
)

type PreprocessingType int

const (
	Stem PreprocessingType = iota
	Lemmatize
)

type IndexerConfig struct {
	preproc_type PreprocessingType
}

type Indexer struct {
	tree *lsm.MergeTree
	preprocessor PreProcessor
	indexer_config IndexerConfig
}

func NewIndexer(config *lsm.Config, indexer_config IndexerConfig) *Indexer {
	return &Indexer{
		tree: lsm.NewMergeTable(config),
		preprocessor: *NewPreprocessor(),
		indexer_config: indexer_config,
	};
}

func (indexer *Indexer) AddDocument(text string, index_doc int) error {
	var processed_text []string
	var err error

	if indexer.indexer_config.preproc_type == Stem {
		processed_text, err = indexer.preprocessor.StemAndRemoveStopWords(text)
		if err != nil {
			return err
		}
	} else {
		processed_text, err = indexer.preprocessor.LemmatizeAndRemoveStopWords(text)
		if err != nil {
			return err
		}
	}

	for _, word := range processed_text {
		val, contains := indexer.tree.Search(word)
		var bitmap *RoaringBitmaps
		if !contains {
			val = ""
			bitmap = NewRoaringBitmaps(65000)
		} else {
			bitmap = NewRoaringBitmaps(1)
			err = bitmap.Deserialize(val)
		}
		if err != nil {
			return err
		}

		bitmap.Insert(index_doc)
		indexer.tree.Insert(word, bitmap.Serialize());
	}
	return nil
}

func (indexer *Indexer) GetListDocuments(word string, limit int) ([]int, error) {
	val, contains := indexer.tree.Search(word)
	if !contains {
		return []int{}, nil
	}
	bitmap := NewRoaringBitmaps(1)
	err := bitmap.Deserialize(val)

	if err != nil {
		return []int{}, err
	}
	return bitmap.Enumerate(), nil
}

func MergeSortedSlices(slices [][]int) []int {
	merged := []int{}
	indices := make([]int, len(slices))

	for {
		minIndex := -1
		minValue := 0

		for i := 0; i < len(slices); i++ {
			if indices[i] < len(slices[i]) {
				if minIndex == -1 || slices[i][indices[i]] < minValue {
					minIndex = i
					minValue = slices[i][indices[i]]
				}
			}
		}

		if minIndex == -1 {
			break
		}

		merged = append(merged, minValue)
		indices[minIndex]++
	}

	return merged
}

func (indexer *Indexer) GetMergedListsDocuments(text string, limit int) ([]int, error) {
	var processed_text []string
	var err error

	if indexer.indexer_config.preproc_type == Stem {
		processed_text, err = indexer.preprocessor.StemAndRemoveStopWords(text)
		if err != nil {
			return []int{}, err
		}
	} else {
		processed_text, err = indexer.preprocessor.LemmatizeAndRemoveStopWords(text)
		if err != nil {
			return []int{}, err
		}
	}

	merged := [][]int{}
	total_size := 0
	for _, word := range processed_text {
		if (total_size >= limit) {
			break
		}
		word_result, err := indexer.GetListDocuments(word, limit - total_size)
		if err != nil {
			return []int{}, err
		}
		merged = append(merged, word_result)
		total_size += len(word_result)
	}

	return MergeSortedSlices(merged), nil	
}
