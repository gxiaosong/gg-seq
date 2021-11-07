package biz

import (
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

type IdGenerator interface {
	NextId() uint64
	NextIds(batchSize int) []uint64
}

type IdGeneratorCreator func(bizType string) IdGenerator

type IdGeneratorFactory struct {
	generators sync.Map
	mu         sync.Mutex
	log        *log.Helper
	service    *SeqUsecase
}

type CachedIdGenerator struct {
	bizType string
	Current *Segment
	Next    *Segment
	mu      sync.Mutex
	serivce SegmentService
}

func NewIdGeneratorFactory(logger log.Logger, service *SeqUsecase) *IdGeneratorFactory {
	return &IdGeneratorFactory{
		service: service,
		log:     log.NewHelper(log.With(logger, "module", "data/IdGeneratorFactory")),
	}
}

func NewCachedIdGenerator(bizType string, service SegmentService) *CachedIdGenerator {
	return &CachedIdGenerator{
		bizType: bizType,
		serivce: service,
	}
}

func (c *CachedIdGenerator) loadCurrent() {
	c.mu.Lock()
	if c.Current == nil {
		if c.Next == nil {
			c.Current = c.querySegment()
		} else {
			c.Current = c.Next
			c.Next = nil
		}
	}
	defer c.mu.Unlock()
}

func (c *CachedIdGenerator) querySegment() *Segment {
	var (
		s   *Segment
		err error
	)
	s, err = c.serivce.GetNextSegment(c.bizType)
	if err != nil {
		panic(err)
	}

	return s
}

func (c *CachedIdGenerator) loadNext() {
	if c.Next == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.Next == nil {
			go func() {
				c.Next = c.querySegment()
			}()
		}
	}
}
func (c *CachedIdGenerator) NextId() uint64 {
	for {
		if c.Current == nil {
			c.loadCurrent()
			continue
		}
		result := c.Current.NextId()
		if result.Code == Over {
			c.loadCurrent()
		} else {
			if result.Code == Loading {
				c.loadNext()
			}

			return result.Id
		}
	}
}

func (c *CachedIdGenerator) NextIds(batchSize int) []uint64 {
	var ids []uint64
	for i := 0; i < int(batchSize); i++ {
		ids = append(ids, c.NextId())
	}
	return ids
}

func (f *IdGeneratorFactory) getGenerator(bizType string, idCreator IdGeneratorCreator) IdGenerator {
	if v, ok := f.generators.Load(bizType); ok {
		return v.(IdGenerator)
	}
	f.mu.Lock()
	if v, ok := f.generators.Load(bizType); ok {
		f.mu.Unlock()
		return v.(IdGenerator)
	}
	g := idCreator(bizType)
	f.generators.Store(bizType, g)
	f.mu.Unlock()
	return g
}

func (f *IdGeneratorFactory) GetGenerator(bizType string) IdGenerator {
	return f.getGenerator(bizType, func(bizType string) IdGenerator {
		return NewCachedIdGenerator(bizType, f.service)
	})
}
