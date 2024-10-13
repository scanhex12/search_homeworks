package index

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestBitmap(t *testing.T) {
	bitmap := NewBitmap(10, 0)
	assert.False(t, bitmap.Contains(0));

	bitmap.Insert(10);
	assert.False(t, bitmap.Contains(9));
	assert.False(t, bitmap.Contains(11));
	assert.True(t, bitmap.Contains(10));

	assert.False(t, bitmap.Insert(640));
	assert.False(t, bitmap.Contains(640));

	assert.True(t, bitmap.Insert(639));
	assert.True(t, bitmap.Contains(639));
}

func TestBitmapSerialization(t *testing.T) {
	bitmap := NewBitmap(10, 0)

	bitmap.Insert(23);
	bitmap.Insert(42);
	bitmap.Insert(123);

	serialized := bitmap.Serialize();
	new_bitmap := NewBitmap(1, 0)
	new_bitmap.Deserialize(serialized);
	
	assert.True(t, new_bitmap.Contains(23));
	assert.True(t, new_bitmap.Contains(42));
	assert.True(t, new_bitmap.Contains(123));
	assert.False(t, new_bitmap.Contains(100));
}

func TestBitmapEnumerate(t *testing.T) {
	bitmap := NewBitmap(10, 0)

	bitmap.Insert(23);
	bitmap.Insert(42);
	bitmap.Insert(123);

	assert.Equal(t, bitmap.Enumerate(), []int{23, 42,123});
}

func TestRoaringBitmaps(t *testing.T) {
	bitmap := NewRoaringBitmaps(3);
	
	assert.False(t, bitmap.Contains(0));

	bitmap.Insert(10);
	assert.False(t, bitmap.Contains(9));
	assert.False(t, bitmap.Contains(11));
	assert.True(t, bitmap.Contains(10));

	assert.True(t, bitmap.Insert(640));
	assert.True(t, bitmap.Contains(640));

	assert.True(t, bitmap.Insert(639));
	assert.True(t, bitmap.Contains(639));
	assert.True(t, bitmap.Contains(640));
}

func TestRoaringBitmapSerialization(t *testing.T) {
	bitmap := NewRoaringBitmaps(10)

	bitmap.Insert(23);
	bitmap.Insert(42);
	bitmap.Insert(123);

	serialized := bitmap.Serialize();
	new_bitmap := NewRoaringBitmaps(1)
	new_bitmap.Deserialize(serialized);
	
	assert.True(t, new_bitmap.Contains(23));
	assert.True(t, new_bitmap.Contains(42));
	assert.True(t, new_bitmap.Contains(123));
	assert.False(t, new_bitmap.Contains(100));
}

func TestRoaringBitmapEnumerate(t *testing.T) {
	bitmap := NewRoaringBitmaps(3)

	bitmap.Insert(23);
	bitmap.Insert(42);
	bitmap.Insert(123);

	result := bitmap.Enumerate()
	sort.Ints(result)

	assert.Equal(t, result, []int{23, 42,123});
}