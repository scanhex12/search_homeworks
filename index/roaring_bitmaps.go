package index

import (
    "errors"
	"fmt"
	"strconv"
	"strings"
)

type Bitmap struct {
	mask []uint64
	cap int
	start int
}

func NewBitmap(cap, start int) *Bitmap {
	return &Bitmap{mask: make([]uint64, cap), cap: cap, start: start};
}

func (bm *Bitmap) Or(other *Bitmap) error {
	if (bm.cap != other.cap || bm.start != other.start) {
		return errors.New("");
	}
	
	for i := 0; i < bm.cap; i++ {
		bm.mask[i] |= other.mask[i];
	}
	return nil;
}

func (bm *Bitmap) And(other *Bitmap) error {
	if (bm.cap != other.cap || bm.start != other.start) {
		return errors.New("");
	}
	
	for i := 0; i < bm.cap; i++ {
		bm.mask[i] |= other.mask[i];
	}
	return nil;
}

func (bm *Bitmap) Contains(query int) bool {
	if query < bm.start || query >= bm.start + 64 * bm.cap {
		return false;
	}

	index := (query - bm.start) / 64;
	offset := (query - bm.start) % 64;
	remined := (bm.mask[index] >> uint64(offset)) % 2;
	return remined == 1;
}

func (bm* Bitmap) Insert(query int) bool {
	if bm.Contains(query) {
		return false;
	}

	if query < bm.start || query >= bm.start + 64 * bm.cap {
		return false;
	}

	index := (query - bm.start) / 64;
	offset := (query - bm.start) % 64;

	bm.mask[index] += (uint64(1) << offset);
	return true;
}

func (bm* Bitmap) Erase(query int) bool {
	if !bm.Contains(query) {
		return false;
	}

	index := (query - bm.start) / 64;
	offset := (query - bm.start) % 64;

	bm.mask[index] -= (uint64(1) << offset);
	return true;
}

func (b *Bitmap) Serialize() string {
	maskStr := make([]string, len(b.mask))
	for i, val := range b.mask {
		maskStr[i] = strconv.FormatUint(val, 10)
	}

	return strings.Join(maskStr, ",") + " " + strconv.Itoa(b.cap) + " " + strconv.Itoa(b.start)
}

func (b *Bitmap) Deserialize(data string) error {
	parts := strings.Split(data, " ")
	if len(parts) < 3 {
		return fmt.Errorf("invalid input string")
	}

	maskStr := strings.Split(parts[0], ",")
	b.mask = make([]uint64, len(maskStr))
	for i, val := range maskStr {
		num, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		b.mask[i] = num
	}

	capacity, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	b.cap = capacity

	start, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}
	b.start = start

	return nil
}

func (b* Bitmap) Enumerate() []int {
	var answer []int
	for i := 0; i < b.cap; i++ {
		remined := b.mask[i]
		for j := 0; j < 64; j++ {
			if remined % 2 == 1 {
				answer = append(answer, b.start + 64 * i + j)
			}
			remined >>= 1;
		}
	}
	return answer;
}

type RoaringBitmaps struct {
	blocks map[int]*Bitmap
	block_cap int
}

func NewRoaringBitmaps(block_cap int) *RoaringBitmaps {
	return &RoaringBitmaps{
		blocks: make(map[int]*Bitmap, 0),
		block_cap: block_cap,
	};
}

func (rbm *RoaringBitmaps) Insert(query int) bool {
	bitmap_index := query / rbm.block_cap;
	bitmap, contains := rbm.blocks[bitmap_index];
	if !contains {
		rbm.blocks[bitmap_index] = NewBitmap(rbm.block_cap, bitmap_index * rbm.block_cap);
		bitmap = rbm.blocks[bitmap_index];
	}
	return bitmap.Insert(query);
}

func (rbm *RoaringBitmaps) Contains(query int) bool {
	bitmap_index := query / rbm.block_cap;
	bitmap, contains := rbm.blocks[bitmap_index];
	if !contains {
		return false;
	}
	return bitmap.Contains(query);
}

func (r *RoaringBitmaps) Serialize() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(r.block_cap))
	sb.WriteString(" ")

	for key, bitmap := range r.blocks {
		sb.WriteString(strconv.Itoa(key))
		sb.WriteString(":")
		sb.WriteString(bitmap.Serialize())
		sb.WriteString(";")
	}

	return sb.String()
}

func (r *RoaringBitmaps) Deserialize(data string) error {
	parts := strings.SplitN(data, " ", 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid input string")
	}

	blockCap, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	r.block_cap = blockCap

	r.blocks = make(map[int]*Bitmap)

	blockParts := strings.Split(parts[1], ";")
	for _, block := range blockParts {
		if block == "" {
			continue
		}

		keyVal := strings.SplitN(block, ":", 2)
		if len(keyVal) != 2 {
			return fmt.Errorf("invalid block format")
		}

		key, err := strconv.Atoi(keyVal[0])
		if err != nil {
			return err
		}

		bitmap := &Bitmap{}
		err = bitmap.Deserialize(keyVal[1])
		if err != nil {
			return err
		}

		r.blocks[key] = bitmap
	}

	return nil
}

func (r *RoaringBitmaps) Enumerate() []int {
	var result []int

	for _, bitmap := range r.blocks {
		result = append(result, bitmap.Enumerate()...)
	}

	return result
}
