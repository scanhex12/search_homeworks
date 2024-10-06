package main

type MergeTree struct {
	ram_memory map[string]string
	sstables []*SSTable
	config *Config
	generator_names *GeneratorNames
}

func NewMergeTable(config *Config) *MergeTree {
	return &MergeTree{
		ram_memory: make(map[string]string, 0),
		sstables : make([]*SSTable, config.max_sstables),
		config: config,
		generator_names: NewGeneratorNames(),
	}
}

func (tree *MergeTree) Search(key string) (string, bool) {
	val, exists := tree.ram_memory[key]
	if exists {
		return val, exists
	}
	
	for _, sstable := range tree.sstables {
		if sstable == nil {
			continue
		}
		value_ptr, err := sstable.Search(key)
		if err != nil {
			panic("Some errors in sstale: " + err.Error())
		}
		if value_ptr != nil {
			return *value_ptr, true
		}
	}
	return "", false
}

func (tree *MergeTree) Insert(key, value string) error {
	if len(tree.ram_memory) < tree.config.ram_limit {
		tree.ram_memory[key] = value;
		return nil
	}

	_, exists := tree.ram_memory[key]
	if exists {
		tree.ram_memory[key] = value;
		return nil
	}

	filter_file_name := tree.generator_names.GenerateNextNameFilter()
	memtable_file_name := tree.generator_names.GenerateNextNameMemtable()
	new_sstable := NewSSTable(filter_file_name, memtable_file_name, 1, tree.config.block_size, tree.config.capacity_filter, tree.config.num_filters)
	new_sstable.Insert(key, value)
	var err error
	err = nil

	for i := 0; i < len(tree.sstables); i++ {
		if tree.sstables[i] == nil {
			tree.sstables[i] = new_sstable
			return nil
		} else {
			new_sstable, err = Unite(tree.generator_names, tree.sstables[i], new_sstable)
			if err != nil {
				return err
			}
			tree.sstables[i] = nil
		}
	}
	return nil
}
