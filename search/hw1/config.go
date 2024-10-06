package main

import "fmt"

type GeneratorNames struct {
	iter_memtable int
	iter_fiter int
}

func NewGeneratorNames() *GeneratorNames {
	return &GeneratorNames{
		iter_memtable: 0,
		iter_fiter: 0,
	}
}

func (gen *GeneratorNames) GenerateNextNameMemtable() string {
	gen.iter_memtable += 1
	return fmt.Sprintf("data/memtable_%d.txt", gen.iter_memtable)
}

func (gen *GeneratorNames) GenerateNextNameFilter() string {
	gen.iter_fiter += 1
	return fmt.Sprintf("data/filter_%d.txt", gen.iter_fiter)
}

type Config struct {
	ram_limit int
	block_size int
	max_sstables int
	capacity_filter int
	num_filters int
}

func NewConfig(ram_limit, block_size,max_sstables, capacity_filter, num_filters int) *Config {
	return &Config{
		ram_limit: ram_limit,
		block_size: block_size,
		max_sstables: max_sstables,
		capacity_filter: capacity_filter,
		num_filters: num_filters,
	}
}
