package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
)
func TestBloomFilter(t *testing.T) {
	filter := NewBloomFilter(5, 1)
	assert.False(t, filter.Contains("key5"))
	filter.Add("key1")
	assert.True(t, filter.Contains("key1"))
}

func TestBloomFilter2(t *testing.T) {
	os.Remove("filter.txt")
	defer os.Remove("filter.txt")

	filter := NewBloomFilter(5, 1)
	filter.Add("key1")
	filter.Add("key2")
	filter.Add("key3")
	filter.SaveToFile("filter.txt")
	
	filter2, err := LoadFromFile("filter.txt")
	assert.Equal(t, err, nil)
	assert.True(t, filter2.Contains("key1"))
	assert.True(t, filter2.Contains("key2"))
	assert.True(t, filter2.Contains("key3"))
}
