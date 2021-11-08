package biz

import (
	"sync"
	"sync/atomic"

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
	bizType       string
	Current       *Segment
	Next          *Segment
	mu            sync.Mutex
	serivce       SegmentService
	isLoadingNext atomic.Value
}

func NewIdGeneratorFactory(logger log.Logger, service *SeqUsecase) *IdGeneratorFactory {
	return &IdGeneratorFactory{
		service: service,
		log:     log.NewHelper(log.With(logger, "biz/IdGeneratorFactory")),
	}
}

func NewCachedIdGenerator(bizType string, service SegmentService) *CachedIdGenerator {
	c := &CachedIdGenerator{
		bizType: bizType,
		serivce: service,
	}
	c.isLoadingNext.Store(false)
	return c
}

func (c *CachedIdGenerator) loadCurrent() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Current == nil || !c.Current.Useful() {
		if c.Next == nil {
			s, err := c.querySegment()
			if err != nil {
				panic(err)
			}
			c.Current = s
		} else {
			c.Current = c.Next
			c.Next = nil
		}
	}

}

func (c *CachedIdGenerator) querySegment() (*Segment, error) {
	var (
		s   *Segment
		err error
	)
	s, err = c.serivce.GetNextSegment(c.bizType)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *CachedIdGenerator) loadNext() {
	if c.Next == nil && !c.isLoadingNext.Load().(bool) {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.Next == nil {
			c.isLoadingNext.Store(true)
			go func() {
				defer func() {
					if e := recover(); e != nil {
						c.isLoadingNext.Store(false)
					}
				}()
				s, err := c.querySegment()
				if err != nil {
					panic(err)
				}
				c.Next = s
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
	defer f.mu.Unlock()
	if v, ok := f.generators.Load(bizType); ok {
		return v.(IdGenerator)
	}
	g := idCreator(bizType)
	f.generators.Store(bizType, g)
	return g
}

func (f *IdGeneratorFactory) GetGenerator(bizType string) IdGenerator {
	return f.getGenerator(bizType, func(bizType string) IdGenerator {
		return NewCachedIdGenerator(bizType, f.service)
	})
}
