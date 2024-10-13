package lsm

type SSTable struct {
	filter_file_name string 
	memtable_file_name string
	level int
	block_size int
	capacity_filter int
	num_filters int
}

func NewSSTable(filter_file_name, memtable_file_name string, level, block_size, capacity_filter, num_filters  int) *SSTable {
	return &SSTable{
		filter_file_name : filter_file_name,
		memtable_file_name : memtable_file_name,
		level : level,
		block_size : block_size,
		capacity_filter : capacity_filter,
		num_filters: num_filters,
	}
}

func (sst *SSTable) Search(key string) (*string, error) {
	filter, err := LoadFromFile(sst.filter_file_name)
	if err != nil {
		return nil, err
	}
	
	if !filter.Contains(key) {
		return nil, nil
	}

	memtable := NewMemtable(sst.memtable_file_name, sst.block_size)
	res, contains := memtable.BinarySearch(key)
	if !contains {
		return nil, nil
	}
	return &res, nil
}

func (sst *SSTable) Insert(key, value string) error {
	filter, err := LoadFromFile(sst.filter_file_name)
	if err != nil {
		filter = NewBloomFilter(uint(sst.capacity_filter), uint(sst.num_filters))
	}
	filter.Add(key)
	filter.SaveToFile(sst.filter_file_name)

	memtable := NewMemtable(sst.memtable_file_name, sst.block_size)
	err = memtable.AddKeyValue(key, value)
	return err
}

func Unite(generator *GeneratorNames, sst1, sst2 *SSTable) (*SSTable, error) {
	file_name := generator.GenerateNextNameFilter()

	filter1, err := LoadFromFile(sst1.filter_file_name)
	if err != nil {
		return nil, err
	}
	filter2, err := LoadFromFile(sst2.filter_file_name)
	if err != nil {
		return nil, err
	}

	result_filter, err := filter1.Unite(filter2)
	err = result_filter.SaveToFile(file_name)
	if err != nil {
		return nil, err
	}
	
	memtable_1 := NewMemtable(sst1.memtable_file_name, sst1.block_size)
	memtable_2 := NewMemtable(sst2.memtable_file_name, sst2.block_size)

	new_memtable_file_name := generator.GenerateNextNameMemtable() 
	err = memtable_1.Merge(memtable_2, new_memtable_file_name)
	if err != nil {
		return nil, err
	}
	return &SSTable{
		filter_file_name : file_name,
		memtable_file_name : new_memtable_file_name,
		level : sst1.level + 1,
		block_size : sst1.block_size,
	}, nil
}
