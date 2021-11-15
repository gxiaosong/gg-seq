package service

import (
	"sync"

	"github.com/gouez/gg-seq/comm"
)

type defaultIdGeneratorFactory struct {
	cache map[string]comm.IdGenerator
	mu    sync.Mutex
}

func NewIdGeneratorFactory() comm.IdGeneratorFactory {
	return &defaultIdGeneratorFactory{
		cache: make(map[string]comm.IdGenerator),
	}
}

func (d *defaultIdGeneratorFactory) GetIdGenerator(bizType string) comm.IdGenerator {
	if v, ok := d.cache[bizType]; ok {
		return v
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	return nil
}
