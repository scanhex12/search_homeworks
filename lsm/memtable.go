package lsm

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type KeyValue struct {
	Key   string
	Value string
}

type Memtable struct {
	filePath string
	blockSize int
}

func NewMemtable(filePath string, blockSize int) *Memtable {
	return &Memtable{
		filePath: filePath,
		blockSize: blockSize,
	}
}

func (m *Memtable) LoadBlock(startIndex int) ([]KeyValue, error) {
	file, err := os.Open(m.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data []KeyValue
	currentIndex := 0

	for scanner.Scan() {
		if currentIndex >= startIndex && currentIndex < startIndex+m.blockSize {
			line := scanner.Text()
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				data = append(data, KeyValue{Key: parts[0], Value: parts[1]})
			}
		}
		currentIndex++

		if currentIndex >= startIndex+m.blockSize {
			break
		}
	}

	return data, scanner.Err()
}

func (m *Memtable) BinarySearch(key string) (string, bool) {
	blockSize := m.blockSize
	startIndex := 0
	var result string
	found := false

	for {
		block, err := m.LoadBlock(startIndex)
		if err != nil {
			fmt.Println("Ошибка загрузки блока:", err)
			return "", false
		}

		if len(block) == 0 {
			break
		}

		sort.Slice(block, func(i, j int) bool {
			return block[i].Key < block[j].Key
		})

		index := sort.Search(len(block), func(i int) bool {
			return block[i].Key >= key
		})

		for index < len(block) && block[index].Key == key {
			result = block[index].Value
			found = true
			index++
		}

		startIndex += blockSize
	}

	return result, found
}

func (m *Memtable) PrintFileContents() error {
	file, err := os.Open(m.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return scanner.Err()
}

func (m *Memtable) Merge(other *Memtable, outputFilePath string) error {
	block2 := make([]KeyValue, 0)

	i := 0
	j := 0

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	for {
		block1, err := m.LoadBlock(i)
		if err != nil {
			return err
		}
		block2, err = other.LoadBlock(j)
		if err != nil {
			return err
		}

		if len(block1) == 0 && len(block2) == 0 {
			break
		}

		for len(block1) > 0 && len(block2) > 0 {
			if block1[0].Key <= block2[0].Key {
				writer.WriteString(fmt.Sprintf("%s:%s\n", block1[0].Key, block1[0].Value))
				block1 = block1[1:]
			} else {
				writer.WriteString(fmt.Sprintf("%s:%s\n", block2[0].Key, block2[0].Value))
				block2 = block2[1:]
			}
		}

		for len(block1) > 0 {
			writer.WriteString(fmt.Sprintf("%s:%s\n", block1[0].Key, block1[0].Value))
			block1 = block1[1:]
		}

		for len(block2) > 0 {
			writer.WriteString(fmt.Sprintf("%s:%s\n", block2[0].Key, block2[0].Value))
			block2 = block2[1:]
		}

		i += m.blockSize
		j += other.blockSize 
	}

	writer.Flush() 
	return nil
}

func (m *Memtable) AddKeyValue(key, value string) error {
	file, err := os.OpenFile(m.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s:%s\n", key, value))
	if err != nil {
		return err
	}

	return nil
}
