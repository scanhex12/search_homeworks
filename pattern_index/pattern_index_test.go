package patternindex

import (
	"testing"
	"search/lsm"
	"search/index"
	"github.com/stretchr/testify/assert"
)

func TestPrefixIndex(t *testing.T) {
	lsm_config := lsm.NewConfig(4, 4, 8, 5, 1)
	index_config := index.NewIndexConfig()

	indexer_config := &PatternIndexConfig{
		index_document_lsm: lsm_config,
		index_document: index_config,
		documents_config: lsm_config,
		ngram_coef: -1,
	}

	indexer := NewPatternIndex(indexer_config)
	indexer.InsertPrefixDocuments("some strange text", 0)
	indexer.InsertPrefixDocuments("another strange text", 1)
	res, err := indexer.SearchByPrefix("ano", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"another strange text"})

	res, err = indexer.SearchByPrefix("str", 100)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"some strange text", "another strange text"})

}

func TestPatternIndex(t *testing.T) {
	lsm_config := lsm.NewConfig(4, 4, 8, 5, 1)
	index_config := index.NewIndexConfig()

	indexer_config := &PatternIndexConfig{
		index_document_lsm: lsm_config,
		index_document: index_config,
		documents_config: lsm_config,
		ngram_coef: 3,
	}

	indexer := NewPatternIndex(indexer_config)
	indexer.InsertPrefixDocuments("some strange text", 0)
	indexer.InsertPrefixDocuments("another strange text", 1)
	res, err := indexer.SearchByPattern("a*ther", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"another strange text"})
}
