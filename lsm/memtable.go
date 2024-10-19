package lsm

import (
	"bufio"
	//"errors"
	"fmt"
	"os"
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

func (m *Memtable) GetKeyValueStr(key, value string) string {
	res := fmt.Sprintf("%s:%s", key, value)
	for len(res) < m.blockSize - 1 {
		res = fmt.Sprintf("%s^", res)
	}
	return fmt.Sprintf("%s\n", res)
}

func (m *Memtable) ParseStr(line string) KeyValue {
	parts := strings.SplitN(line, ":", 2)
	key := parts[0]
	value := strings.SplitN(parts[1], "^", 2)
	data := KeyValue{Key: key, Value: value[0]}
	return data
}

func (m *Memtable) LoadBlock(startIndex int) (*KeyValue, error) {
	file, err := os.Open(m.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	offset := int64(startIndex) * int64(m.blockSize)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)

	var data KeyValue

	scanner.Scan()
	line := scanner.Text()
	if len(line) != m.blockSize - 1 {
		return nil, nil
	}
	data = m.ParseStr(line)

	return &data, scanner.Err()
}

func (m *Memtable) BinarySearch(key string) (string, bool) {
	fileInfo, err := os.Stat(m.filePath)
	if err != nil {
		panic(err)
	}
	fileSize := fileInfo.Size()
	result, found := "", false
	left, right := int64(0), (fileSize/int64(m.blockSize))-1

	for left <= right {
		mid := left + (right-left)/2

		block, err := m.LoadBlock(int(mid))
		if err != nil {
			panic(err)
		}
		key_mid, value_mid := block.Key, block.Value
		if key_mid == key {
			result = value_mid
			found = true
			left = mid + 1
		} else if key_mid < key {
			left = mid + 1
		} else {
			right = mid - 1
		}
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
	var block2 *KeyValue 

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

		if block1 == nil && block2 == nil {
			break
		}
		
		if block1 != nil && block2 != nil {
			if block1.Key <= block2.Key {
				writer.WriteString(m.GetKeyValueStr(block1.Key, block1.Value))
				i += 1
				block1 = nil
			} else {
				writer.WriteString(m.GetKeyValueStr(block2.Key, block2.Value))
				j += 1
				block2 = nil
			}
		}

		if block1 != nil {
			writer.WriteString(m.GetKeyValueStr(block1.Key, block1.Value))
			i += 1
			block1 = nil
		}

		if block2 != nil {
			writer.WriteString(m.GetKeyValueStr(block2.Key, block2.Value))
			j += 1
			block2 = nil
		}
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


	_, err = file.WriteString(m.GetKeyValueStr(key, value))
	if err != nil {
		return err
	}

	return nil
}
