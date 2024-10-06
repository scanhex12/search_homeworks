package main

import (
	"encoding/gob"
	"errors"
	"hash/fnv"
	"math"
	"os"
)

type BloomFilter struct {
	Bitset []bool
	Size   uint
	K      uint
}

func NewBloomFilter(size uint, numHashFunction uint) *BloomFilter {
	return &BloomFilter{
		Bitset: make([]bool, size),
		Size: size,
		K: numHashFunction,
	}
}

func (bf *BloomFilter) Add(item string) {
	for i := uint(0); i < bf.K; i++ {
		hash := bf.hash(item, i)
		index := hash % bf.Size
		bf.Bitset[index] = true
	}
}

func (bf *BloomFilter) Contains(item string) bool {
	for i := uint(0); i < bf.K; i++ {
		hash := bf.hash(item, i)
		index := hash % bf.Size
		if !bf.Bitset[index] {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) hash(item string, salt uint) uint {
	h := fnv.New64a()
	h.Write([]byte(item))
	h.Write([]byte{byte(salt)})
	return uint(h.Sum64())
}

func (bf *BloomFilter) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(bf)
	if err != nil {
		return err
	}

	return nil
}

func LoadFromFile(filename string) (*BloomFilter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var bf BloomFilter
	err = decoder.Decode(&bf)
	if err != nil {
		return nil, err
	}

	return &bf, nil
}

func (bf *BloomFilter) Unite(other *BloomFilter) (*BloomFilter, error) {
	if bf.Size != other.Size {
		return nil, errors.New("Non equals sizes fo bloom filters");
	}
	if bf.K != other.K {
		return nil, errors.New("Non equals K of bloom filters");
	}

	new_filter := NewBloomFilter(bf.Size, bf.K)
	for i := 0; i < int(bf.Size); i++ {
		new_filter.Bitset[i] = bf.Bitset[i] || other.Bitset[i]
	}
	return new_filter, nil
}

func OptimalNumHashFuncs(size, numItems uint) uint {
	return uint(math.Ceil((float64(size) / float64(numItems)) * math.Ln2))
}
