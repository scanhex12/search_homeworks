package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: should run lematize service before run
func TestLemmatize(t *testing.T) {
	processor := NewPreprocessor()
	result, err := processor.Lemmatize("running better feet")

	assert.Equal(t, err, nil)
	assert.Equal(t, result, []string{"run", "well", "foot"})
}

func TestStemmer(t *testing.T) {
	processor := NewPreprocessor()
	result, err := processor.Stem("running better feet")

	assert.Equal(t, err, nil)
	assert.Equal(t, result, []string{"run", "better", "feet"})
}

func TestStopWords(t *testing.T) {
	processor := NewPreprocessor()
	result, err := processor.ClassifyStopWords("the quick brown fox jumps")

	assert.Equal(t, err, nil)
	assert.Equal(t, result, StopWordsResponse{"brown":false, "fox":false, "jumps":false, "quick":false, "the":true})
}

func TestLemProcessing(t *testing.T) {
	processor := NewPreprocessor()
	result, err := processor.LemmatizeAndRemoveStopWords("the quick brown fox jumps")

	assert.Equal(t, err, nil)
	assert.Equal(t, result, []string{"quick", "brown", "fox", "jump"})
}

func TestStemProcessing(t *testing.T) {
	processor := NewPreprocessor()
	result, err := processor.StemAndRemoveStopWords("the quick brown fox jumps")

	assert.Equal(t, err, nil)
	assert.Equal(t, result, []string{"quick", "brown", "fox", "jump"})
}

