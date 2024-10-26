package index

import (
	"encoding/base64"
	"search/lsm"
	"github.com/RoaringBitmap/roaring"
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
		var bitmap *roaring.Bitmap
		if !contains {
			val = ""
			bitmap = roaring.NewBitmap()
		} else {
			bitmap = roaring.New()
			data, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return err
			}
			err = bitmap.UnmarshalBinary(data)
		}
		if err != nil {
			return err
		}

		bitmap.Add(uint32(index_doc))
		data, err := bitmap.MarshalBinary()
		if err != nil {
			return err
		}
		base64String := base64.StdEncoding.EncodeToString(data)
	
		indexer.tree.Insert(word, base64String);
	}
	return nil
}


func (indexer *Indexer) AddBatchDocument(texts []string, index_docs []int) error {
	wordBitmaps := make(map[string]*roaring.Bitmap)

	for i, text := range texts {
		var processed_text []string
		var err error

		// Выбор метода обработки: стемминг или лемматизация
		if indexer.indexer_config.preproc_type == Stem {
			processed_text, err = indexer.preprocessor.StemAndRemoveStopWords(text)
		} else {
			processed_text, err = indexer.preprocessor.LemmatizeAndRemoveStopWords(text)
		}
		if err != nil {
			return err
		}

		// Группировка слов и документов в битмапы
		for _, word := range processed_text {
			if _, exists := wordBitmaps[word]; !exists {
				wordBitmaps[word] = roaring.NewBitmap()
			}
			wordBitmaps[word].Add(uint32(index_docs[i]))
		}
	}

	// Обновление дерева одним проходом для всех слов
	for word, newBitmap := range wordBitmaps {
		val, contains := indexer.tree.Search(word)
		var existingBitmap *roaring.Bitmap
		if contains {
			existingBitmap = roaring.New()
			data, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return err
			}
			if err := existingBitmap.UnmarshalBinary(data); err != nil {
				return err
			}
			existingBitmap.Or(newBitmap)
		} else {
			existingBitmap = newBitmap
		}

		data, err := existingBitmap.MarshalBinary()
		if err != nil {
			return err
		}
		base64String := base64.StdEncoding.EncodeToString(data)
		indexer.tree.Insert(word, base64String)
	}
	return nil
}

func (indexer *Indexer) GetListDocuments(word string, limit int) ([]int, error) {
	bitmap, err := indexer.GetBitmapDocuments(word, limit)
	if err != nil {
		return []int{}, err
	}
	uint32Array := bitmap.ToArray()
    intArray := make([]int, len(uint32Array))
    for i, v := range uint32Array {
        intArray[i] = int(v)
    }
	return intArray, nil
}

func (indexer *Indexer) GetBitmapDocuments(word string, limit int) (*roaring.Bitmap, error) {
	val, contains := indexer.tree.Search(word)
	bitmap := roaring.NewBitmap()
	if !contains {
		return bitmap, nil
	}
	data, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil, err
	}
	err = bitmap.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return bitmap, nil
}

func (indexer *Indexer) GetMergedBitmapDocuments(text string, limit int) (*roaring.Bitmap, error) {
	var processed_text []string
	var err error

	if indexer.indexer_config.preproc_type == Stem {
		processed_text, err = indexer.preprocessor.StemAndRemoveStopWords(text)
		if err != nil {
			return nil, err
		}
	} else {
		processed_text, err = indexer.preprocessor.LemmatizeAndRemoveStopWords(text)
		if err != nil {
			return nil, err
		}
	}

	merged := roaring.NewBitmap()
	total_size := 0
	for _, word := range processed_text {
		if (total_size >= limit) {
			break
		}
		word_result, err := indexer.GetBitmapDocuments(word, limit - total_size)
		if err != nil {
			return nil, err
		}
		merged.Or(word_result)
		total_size += int(word_result.DenseSize())
	}

	return merged, nil	
}


func (indexer *Indexer) GetMergedListsDocuments(text string, limit int) ([]int, error) {
	merged, err := indexer.GetBitmapDocuments(text, limit)
	if err != nil {
		return []int{}, err
	}

	uint32Array := merged.ToArray()
    intArray := make([]int, len(uint32Array))
    for i, v := range uint32Array {
        intArray[i] = int(v)
    }

	return intArray, nil	
}
