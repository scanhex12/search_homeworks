package main

import (
	"fmt"
	"time"
	"os"
)

func RebuildState() {
	os.RemoveAll("./data")
	os.MkdirAll("./data", 0755)
}

func BenchmarkInsert(tree *MergeTree, numCheckKeys int) {
	start := time.Now()

	for i := 0; i < 2*numCheckKeys; i++ {
		key := fmt.Sprintf("key_%d", i%numCheckKeys)
		value := fmt.Sprintf("value_%d", i)
		tree.Insert(key, value)
	}

	elapsed := time.Since(start)
	fmt.Printf("Insert Time: %s\n", elapsed)
}

func BenchmarkSearch(tree *MergeTree, numCheckKeys int) {
	start := time.Now()

	for i := 0; i < numCheckKeys; i++ {
		key := fmt.Sprintf("key_%d", i)
		_, exists := tree.Search(key)
		if !exists {
			fmt.Printf("Key %s not found\n", key)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Search Time: %s\n", elapsed)
}

func BenchmarkRangeSearch(tree *MergeTree, numCheckKeys, rangeSize int) {
	start := time.Now()

	for i := 0; i < numCheckKeys; i++ {
		for j := i; j < i+rangeSize && j < numCheckKeys; j++ {
			key := fmt.Sprintf("key_%d", j)
			_, exists := tree.Search(key)
			if !exists {
				fmt.Printf("Key %s not found in range search\n", key)
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Range Search Time: %s\n", elapsed)
}

func main() {
	RebuildState()
	defer RebuildState()

	numCheckKeys := 100
	rangeSize := 5

	config := NewConfig(1, 4, 8, 5, 1)
	tree := NewMergeTable(config)

	BenchmarkInsert(tree, numCheckKeys)

	BenchmarkSearch(tree, numCheckKeys)

	BenchmarkRangeSearch(tree, numCheckKeys, rangeSize)
}
