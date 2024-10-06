package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
)
func TestMemtable1(t *testing.T) {
	os.Remove("a.txt")
	defer os.Remove("a.txt")

	memtable1 := NewMemtable("a.txt", 4)
	memtable1.AddKeyValue("key1", "value1")
	memtable1.AddKeyValue("key2", "value2")
	memtable1.AddKeyValue("key3", "value3")

	val1, exists := memtable1.BinarySearch("key1")
	assert.True(t, exists)
	assert.Equal(t, val1, "value1")

	val2, exists := memtable1.BinarySearch("key2")
	assert.True(t, exists)
	assert.Equal(t, val2, "value2")

	val3, exists := memtable1.BinarySearch("key3")
	assert.True(t, exists)
	assert.Equal(t, val3, "value3")

	_, exists = memtable1.BinarySearch("key4")
	assert.False(t, exists)
}

func TestMemtable2(t *testing.T) {
	os.Remove("b.txt")
	os.Remove("c.txt")
	os.Remove("d.txt")

	defer os.Remove("b.txt")
	defer os.Remove("c.txt")
	defer os.Remove("d.txt")

	memtable1 := NewMemtable("b.txt", 4)
	memtable2 := NewMemtable("c.txt", 4)

	memtable1.AddKeyValue("key1", "value1")
	memtable1.AddKeyValue("key3", "value3")
	memtable1.AddKeyValue("key5", "value5")

	memtable2.AddKeyValue("key2", "value2")
	memtable2.AddKeyValue("key4", "value4")
	memtable2.AddKeyValue("key6", "value6")
	memtable2.AddKeyValue("key8", "value8")
	memtable2.AddKeyValue("key10", "value10")
	memtable2.AddKeyValue("key12", "value12")

	val1, exists := memtable2.BinarySearch("key2")
	assert.True(t, exists)
	assert.Equal(t, val1, "value2")

	val2, exists := memtable2.BinarySearch("key4")
	assert.True(t, exists)
	assert.Equal(t, val2, "value4")
	
	memtable1.Merge(memtable2, "d.txt")
	memtable3 := NewMemtable("d.txt", 4)

	val_memtable1, exists := memtable3.BinarySearch("key1")
	assert.True(t, exists)
	assert.Equal(t, val_memtable1, "value1")

	val_memtable2, exists := memtable3.BinarySearch("key2")
	assert.True(t, exists)
	assert.Equal(t, val_memtable2, "value2")

	val_memtable6, exists := memtable3.BinarySearch("key6")
	assert.True(t, exists)
	assert.Equal(t, val_memtable6, "value6")

	val_memtable12, exists := memtable3.BinarySearch("key12")
	assert.True(t, exists)
	assert.Equal(t, val_memtable12, "value12")

	_, exists = memtable1.BinarySearch("key14")
	assert.False(t, exists)

	_, exists = memtable1.BinarySearch("key9")
	assert.False(t, exists)
}

func TestMemtable3(t *testing.T) {
	os.Remove("b.txt")
	os.Remove("c.txt")
	os.Remove("d.txt")

	defer os.Remove("b.txt")
	defer os.Remove("c.txt")
	defer os.Remove("d.txt")

	memtable1 := NewMemtable("b.txt", 4)
	memtable2 := NewMemtable("c.txt", 4)

	memtable1.AddKeyValue("key1", "value1")
	memtable1.AddKeyValue("key3", "value3")
	memtable1.AddKeyValue("key5", "value5")

	memtable2.AddKeyValue("key1", "value2")
	memtable2.AddKeyValue("key2", "value4")
	memtable2.AddKeyValue("key5", "value6")
	memtable2.AddKeyValue("key8", "value8")
	memtable2.AddKeyValue("key10", "value10")
	memtable2.AddKeyValue("key12", "value12")

	memtable1.Merge(memtable2, "d.txt")
	memtable3 := NewMemtable("d.txt", 4)

	val_memtable1, exists := memtable3.BinarySearch("key1")
	assert.True(t, exists)
	assert.Equal(t, val_memtable1, "value2")

	val_memtable2, exists := memtable3.BinarySearch("key2")
	assert.True(t, exists)
	assert.Equal(t, val_memtable2, "value4")

	val_memtable6, exists := memtable3.BinarySearch("key5")
	assert.True(t, exists)
	assert.Equal(t, val_memtable6, "value6")

	val_memtable12, exists := memtable3.BinarySearch("key12")
	assert.True(t, exists)
	assert.Equal(t, val_memtable12, "value12")

	_, exists = memtable1.BinarySearch("key14")
	assert.False(t, exists)

	_, exists = memtable1.BinarySearch("key9")
	assert.False(t, exists)
}

func TestMemtable4(t *testing.T) {
	os.Remove("a.txt")
	defer os.Remove("a.txt")

	memtable1 := NewMemtable("a.txt", 4)
	memtable1.AddKeyValue("key1", "value1")
	memtable1.AddKeyValue("key1", "value2")
	memtable1.AddKeyValue("key1", "value3")

	val1, exists := memtable1.BinarySearch("key1")
	assert.True(t, exists)
	assert.Equal(t, val1, "value3")
}
