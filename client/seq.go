package client

import (
	"sync"

	"github.com/gouez/gg-seq/comm"
)

var (
	idGeneratorFactory comm.IdGeneratorFactory
	once               sync.Once
)

func SetUrl(url string) {
	once.Do(func() {
		idGeneratorFactory = NewClientIdGeneratorFactory(url)
	})
}

func GetNextId(bizType string) uint64 {

	return idGeneratorFactory.GetIdGenerator(bizType).GetId()
}

func GetNextIds(bizType string, size int) []uint64 {

	return idGeneratorFactory.GetIdGenerator(bizType).GetIds(size)
}
