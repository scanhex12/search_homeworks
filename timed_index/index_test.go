package timed_index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndex(t *testing.T) {
	indexer := NewIndexer()
	indexer.AddDocumentByIntegers(0, 6)
	indexer.AddDocumentByIntegers(1, 7)
	indexer.AddDocumentByIntegers(2, 8)
	indexer.AddDocumentByIntegers(3, 1)
	indexer.AddDocumentByIntegers(4, 2)

	assert.Equal(t, []uint32{1,7,8}, indexer.SearchByIntegersRange(int64(1), int64(3)))
	assert.Equal(t, []uint32{}, indexer.SearchByIntegersRange(int64(5), int64(10)))
	assert.Equal(t, []uint32{1,6,7,8}, indexer.SearchByIntegersRange(int64(0), int64(3)))
	assert.Equal(t, []uint32{7}, indexer.SearchByIntegersRange(int64(1), int64(1)))
}

func TestIntervalIndex(t *testing.T) {
	indexer := NewIntervalIndexer()
	indexer.AddIntegerDocument(0, 1, 6)
	indexer.AddIntegerDocument(1, 5, 7)
	indexer.AddIntegerDocument(2, 3, 8)
	indexer.AddIntegerDocument(3, 8, 1)
	indexer.AddIntegerDocument(4, 6, 2)

	assert.Equal(t, []uint32{1,6,7,8}, indexer.SearchByIntegersRange(int64(1), int64(3)))
	assert.Equal(t, []uint32{1,2,7}, indexer.SearchByIntegersRange(int64(5), int64(10)))
	assert.Equal(t, []uint32{1,6,7,8}, indexer.SearchByIntegersRange(int64(0), int64(3)))
	assert.Equal(t, []uint32{6,7}, indexer.SearchByIntegersRange(int64(1), int64(1)))
}
