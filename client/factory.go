package client

import (
	"sync"

	"github.com/gouez/gg-seq/comm"
)

type clientIdGeneratorFactory struct {
	cache   map[string]comm.IdGenerator
	mu      sync.Mutex
	service comm.SegmentService
}

func NewClientIdGeneratorFactory(url string) comm.IdGeneratorFactory {
	return &clientIdGeneratorFactory{
		cache:   make(map[string]comm.IdGenerator),
		service: NewHttpSegmentService(url),
	}
}
func (d *clientIdGeneratorFactory) GetIdGenerator(bizType string) comm.IdGenerator {
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

func (d *clientIdGeneratorFactory) createIdGenerator(bizType string) comm.IdGenerator {
	return comm.NewCacheIdGenerator(bizType, d.service)
}
