package timed_index

import (
	"math"
	"time"
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

func NewIndexConfig() *IndexerConfig {
	return &IndexerConfig{
		preproc_type : Stem,
	}
}

type Indexer struct {
	bitmaps []*roaring.Bitmap
	inverted_bitmaps []*roaring.Bitmap
	total_docs *roaring.Bitmap
}

func NewIndexer() *Indexer {
	bitmaps := make([]*roaring.Bitmap, 64)
	for i := 0; i < 64; i += 1 {
		bitmaps[i] = roaring.NewBitmap()
	}

	inverted_bitmaps := make([]*roaring.Bitmap, 64)
	for i := 0; i < 64; i += 1 {
		inverted_bitmaps[i] = roaring.NewBitmap()
	}

	return &Indexer{
		bitmaps: bitmaps,
		inverted_bitmaps: inverted_bitmaps,
		total_docs: roaring.NewBitmap(),
	};
}

func ConvertTimeToUnix(dateStr string) (int64, error) {
	layout := "2006-01-02 15:04:05"

	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return 0, err
	}

	unixTime := t.Unix()
	return unixTime, nil
}

func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

func max(a, b int64) int64 {
    if a > b {
        return a
    }
    return b
}

func TraverseTree(v, tl, tr, l, r int64, current_representation *[]bool, results *[][]bool) {
	if tl > tr || l > r {
		return
	}
	if tl == l && tr == r {
		newRepresentation := make([]bool, len(*current_representation))
		copy(newRepresentation, *current_representation)
		*results = append(*results, newRepresentation)
		return
	}
	tm := (tl + tr) / 2;
	*current_representation = append(*current_representation, false)
	TraverseTree(2 * v, tl, tm, l, min(r, tm), current_representation, results)
	(*current_representation)[len(*current_representation) - 1] = true
	TraverseTree(2 * v + 1, tm + 1, tr, max(l, tm + 1), r, current_representation, results)
	*current_representation = (*current_representation)[0:len(*current_representation) - 1]
} 

func ToBitRepresentation(value int64) []bool {
	result := make([]bool, 0) 
	for i := 0; i < 63; i += 1 {
		result = append(result, value % 2 == 1)
		value >>= 1
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }
	return result
}

func (indexer *Indexer) AddDocumentByIntegers(unix_timestamp_start int64, doc_index int) {
	indexer.total_docs.Add(uint32(doc_index))
	bit_representation := ToBitRepresentation(unix_timestamp_start)
	for i, value := range bit_representation {
		if value {
			indexer.bitmaps[i].AddInt(doc_index)
		} else {
			indexer.inverted_bitmaps[i].AddInt(doc_index)
		}
	}
}

func (indexer *Indexer) AddDocument(timestamp string, doc_index int) error {
	unix_timestamp_start, err := ConvertTimeToUnix(timestamp)
	if err != nil {
		return err
	}
	indexer.AddDocumentByIntegers(unix_timestamp_start, doc_index)
	return nil
}

func (indexer *Indexer) SearchByIntegersRange(unix_timestamp_start, unix_timestamp_end int64) []uint32 {
	available_bit_representation := make([][]bool, 0)
	current_representation := make([]bool, 0)
	TraverseTree(1, 0, math.MaxInt64, unix_timestamp_start, unix_timestamp_end, &current_representation, &available_bit_representation)
	
	result_bitmap := roaring.NewBitmap()
	
	for _, good_bit_representation := range available_bit_representation {
		current_bitmap := indexer.total_docs
		for i, bit := range good_bit_representation {
			if bit {
				current_bitmap = roaring.And(current_bitmap, indexer.bitmaps[i])
			} else {
				current_bitmap = roaring.And(current_bitmap, indexer.inverted_bitmaps[i])
			}
		}
		result_bitmap = roaring.Or(result_bitmap, current_bitmap)
	}
	return result_bitmap.ToArray()
}


func (indexer *Indexer) SearchByTimestamps(timestamp_start, timestamp_end string) ([]uint32, error) {
	unix_timestamp_start, err := ConvertTimeToUnix(timestamp_start)
	if err != nil {
		return []uint32{}, err
	}
	unix_timestamp_end, err := ConvertTimeToUnix(timestamp_end)
	if err != nil {
		return []uint32{}, err
	}

	result_bitmap := indexer.SearchByIntegersRange(unix_timestamp_start, unix_timestamp_end)
	return result_bitmap, nil	
}

type IntervalIndexer struct {
	start_bitmaps []*roaring.Bitmap
	inverted_start_bitmaps []*roaring.Bitmap
	end_bitmaps []*roaring.Bitmap
	inverted_end_bitmaps []*roaring.Bitmap
	total_docs *roaring.Bitmap
}

func NewIntervalIndexer() *IntervalIndexer {
	bitmaps := make([]*roaring.Bitmap, 64)
	for i := 0; i < 64; i += 1 {
		bitmaps[i] = roaring.NewBitmap()
	}

	inverted_bitmaps := make([]*roaring.Bitmap, 64)
	for i := 0; i < 64; i += 1 {
		inverted_bitmaps[i] = roaring.NewBitmap()
	}

	end_bitmaps := make([]*roaring.Bitmap, 64)
	for i := 0; i < 64; i += 1 {
		end_bitmaps[i] = roaring.NewBitmap()
	}

	inverted_end_bitmaps := make([]*roaring.Bitmap, 64)
	for i := 0; i < 64; i += 1 {
		inverted_end_bitmaps[i] = roaring.NewBitmap()
	}

	return &IntervalIndexer{
		start_bitmaps: bitmaps,
		inverted_start_bitmaps: inverted_bitmaps,
		end_bitmaps: end_bitmaps,
		inverted_end_bitmaps: inverted_end_bitmaps,
		total_docs: roaring.NewBitmap(),
	};
}

func (indexer *IntervalIndexer) AddIntegerDocument(unix_timestamp_start, unix_timestamp_end int64, doc_index uint32) {
	indexer.total_docs.Add(doc_index)
	bit_start_representation := ToBitRepresentation(unix_timestamp_start)
	bit_end_representation := ToBitRepresentation(unix_timestamp_end)
	
	for i, value := range bit_start_representation {
		if value {
			indexer.start_bitmaps[i].Add(doc_index)
		} else {
			indexer.inverted_start_bitmaps[i].Add(doc_index)
		}
	}

	for i, value := range bit_end_representation {
		if value {
			indexer.end_bitmaps[i].Add(doc_index)
		} else {
			indexer.inverted_end_bitmaps[i].Add(doc_index)
		}
	}
}

func (indexer *IntervalIndexer) AddDocument(timestamp_start, timestamp_end string, doc_index uint32) error {
	unix_timestamp_start, err := ConvertTimeToUnix(timestamp_start)
	if err != nil {
		return err
	}
	unix_timestamp_end, err := ConvertTimeToUnix(timestamp_end)
	if err != nil {
		return err
	}
	indexer.AddIntegerDocument(unix_timestamp_start, unix_timestamp_end, doc_index)
	return nil
}


func (indexer *IntervalIndexer) SearchByIntegersRange(unix_timestamp_start, unix_timestamp_end int64) []uint32 {
	available_start_bit_representation := make([][]bool, 0)
	current_start_representation := make([]bool, 0)
	TraverseTree(1, 0, math.MaxInt64, 0, unix_timestamp_end, &current_start_representation, &available_start_bit_representation)
	
	start_result_bitmap := roaring.NewBitmap()
	
	for _, good_bit_representation := range available_start_bit_representation {
		current_bitmap := indexer.total_docs
		for i, bit := range good_bit_representation {
			if bit {
				current_bitmap = roaring.And(current_bitmap, indexer.start_bitmaps[i])
			} else {
				current_bitmap = roaring.And(current_bitmap, indexer.inverted_start_bitmaps[i])
			}
		}
		start_result_bitmap = roaring.Or(start_result_bitmap, current_bitmap)
	}

	available_end_bit_representation := make([][]bool, 0)
	current_end_representation := make([]bool, 0)
	TraverseTree(1, 0, math.MaxInt64, unix_timestamp_start, math.MaxInt64, &current_end_representation, &available_end_bit_representation)

	end_result_bitmap := roaring.NewBitmap()
	
	for _, good_bit_representation := range available_end_bit_representation {
		current_bitmap := indexer.total_docs
		for i, bit := range good_bit_representation {
			if bit {
				current_bitmap = roaring.And(current_bitmap, indexer.end_bitmaps[i])
			} else {
				current_bitmap = roaring.And(current_bitmap, indexer.inverted_end_bitmaps[i])
			}
		}
		end_result_bitmap = roaring.Or(end_result_bitmap, current_bitmap)
	}

	final_segments := roaring.And(start_result_bitmap, end_result_bitmap)
	return final_segments.ToArray()
}

func (indexer *IntervalIndexer) SearchByTimestamps(timestamp_start, timestamp_end string) ([]uint32, error) {
	unix_timestamp_start, err := ConvertTimeToUnix(timestamp_start)
	if err != nil {
		return []uint32{}, err
	}
	unix_timestamp_end, err := ConvertTimeToUnix(timestamp_end)
	if err != nil {
		return []uint32{}, err
	}
	return indexer.SearchByIntegersRange(unix_timestamp_start, unix_timestamp_end), nil	
}

