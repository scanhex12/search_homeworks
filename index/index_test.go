package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"search/lsm"
	"os"
)

func Rebuild() {
	os.RemoveAll("./data")
	os.MkdirAll("./data", 0755)
}

func TestIndex(t *testing.T) {
	Rebuild()
	defer Rebuild()

	lsm_config := lsm.NewConfig(40, 4, 8, 5, 1)
	index_config := IndexerConfig{
		preproc_type: Stem,
	}

	indexer := NewIndexer(lsm_config, index_config)
	indexer.AddDocument("You are a conversational AI assistant that is provided a list of documents and a user query to answer based on information from the documents. The user also provides an answer mode which can be 'Grounded' or 'Mixed'. For answer mode Grounded only respond with exact facts from documents, for answer mode Mixed answer using facts from documents and your own knowledge. Cite all facts from the documents using <co: doc_id></co> tags.", 0);
	indexer.AddDocument("Considering the recent surge in negative document sentiment regarding platform performance, how has the text analytics platform addressed this issue, and what has been the impact on user engagement and customer satisfaction?.", 1);

	result, err := indexer.GetListDocuments("you", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, []int{0})

	result, err = indexer.GetListDocuments("document", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, []int{0, 1})
}

func TestBatchIndex(t *testing.T) {
	Rebuild()
	defer Rebuild()

	lsm_config := lsm.NewConfig(40, 4, 8, 5, 1)
	index_config := IndexerConfig{
		preproc_type: Stem,
	}

	indexer := NewIndexer(lsm_config, index_config)
	indexer.AddBatchDocument([]string{
		"You are a conversational AI assistant that is provided a list of documents and a user query to answer based on information from the documents. The user also provides an answer mode which can be 'Grounded' or 'Mixed'. For answer mode Grounded only respond with exact facts from documents, for answer mode Mixed answer using facts from documents and your own knowledge. Cite all facts from the documents using <co: doc_id></co> tags.",
		"Considering the recent surge in negative document sentiment regarding platform performance, how has the text analytics platform addressed this issue, and what has been the impact on user engagement and customer satisfaction?.",
	}, []int{0, 1})

	result, err := indexer.GetListDocuments("you", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, []int{0})

	result, err = indexer.GetListDocuments("document", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, []int{0, 1})
}
