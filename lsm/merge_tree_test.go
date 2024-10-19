package lsm

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
	"os"
)

func Rebuild() {
	os.RemoveAll("./data")
	os.MkdirAll("./data", 0755)
}

func TestMergeTree1(t *testing.T) {
	Rebuild()
	defer Rebuild()

	config := NewConfig(1, 100, 8, 5, 1)

	tree := NewMergeTable(config)
	err := tree.Insert("key1", "value1")
	assert.Equal(t, err, nil)

	err = tree.Insert("key2", "value2")
	assert.Equal(t, err, nil)

	value1, exists := tree.Search("key1")
	assert.True(t, exists)
	assert.Equal(t, value1, "value1")

	value2, exists := tree.Search("key2")
	assert.True(t, exists)
	assert.Equal(t, value2, "value2")

	err = tree.Insert("key1", "value3")
	assert.Equal(t, err, nil)

	value3, exists := tree.Search("key1")
	assert.True(t, exists)
	assert.Equal(t, value3, "value3")
}

func TestMergeTree2(t *testing.T) {
	Rebuild()
	defer Rebuild()

	num_check_keys := 4

	config := NewConfig(1, 100, 8, 5, 1)
	tree := NewMergeTable(config)

	for i := 0; i < 2 * num_check_keys; i++ {
		key := fmt.Sprintf("key_%d", i % num_check_keys)
		value := fmt.Sprintf("value_%d", i)
		tree.Insert(key, value)
	}

	for i := 0; i < num_check_keys; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i + num_check_keys)
		result, exists := tree.Search(key)
		assert.True(t, exists)
		assert.Equal(t, value, result)
	}

	for i := num_check_keys; i < 2 * num_check_keys; i++ {
		key := fmt.Sprintf("key_%d", i)
		_, exists := tree.Search(key)
		assert.False(t, exists)
	}
}