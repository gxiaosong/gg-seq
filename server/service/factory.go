package service

import (
	"sync"

	"github.com/gouez/gg-seq/comm"
)

type defaultIdGeneratorFactory struct {
	cache   map[string]comm.IdGenerator
	mu      sync.Mutex
	service comm.SegmentService
}

func NewIdGeneratorFactory(service comm.SegmentService) comm.IdGeneratorFactory {
	return &defaultIdGeneratorFactory{
		cache:   make(map[string]comm.IdGenerator),
		service: service,
	}
}

func (d *defaultIdGeneratorFactory) GetIdGenerator(bizType string) comm.IdGenerator {
	if v, ok := d.cache[bizType]; ok {
		return v
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if v, ok := d.cache[bizType]; ok {
		return v
	}
	idGen := d.createIdGenerator(bizType)
	d.cache[bizType] = idGen
	return idGen
}

func (d *defaultIdGeneratorFactory) createIdGenerator(bizType string) comm.IdGenerator {
	return comm.NewCacheIdGenerator(bizType, d.service)
}
