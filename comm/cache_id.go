package comm

import "sync"

type cacheIdGenerator struct {
	bizType       string
	service       SegmentService
	current       *Segment
	next          *Segment
	isLoadingNext bool
	mu            sync.Mutex
}

func NewCacheIdGenerator(bizType string, service SegmentService) IdGenerator {

	g := &cacheIdGenerator{
		service: service,
		bizType: bizType,
	}
	g.loadCurrent()
	return g
}

func (c *cacheIdGenerator) loadCurrent() {
	if c.current == nil || !c.current.Useful() {
		if c.next == nil {
			s := c.querySegment()
			c.current = s
		} else {
			c.current = c.next
			c.next = nil
		}
	}
}
func (c *cacheIdGenerator) loadNext() {
	if c.next == nil && !c.isLoadingNext {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.next == nil && !c.isLoadingNext {
			c.isLoadingNext = true
			go func() {
				defer func() {
					c.isLoadingNext = false
				}()
				c.next = c.querySegment()
			}()
		}
	}
}

func (c *cacheIdGenerator) querySegment() *Segment {
	s := c.service.GetNextSegment(c.bizType)
	return s
}

func (c *cacheIdGenerator) GetId() uint64 {
	for {
		if c.current == nil {
			c.loadCurrent()
			continue
		}
		result := c.current.NextId()
		if result.Code == OVER {
			c.loadCurrent()
		} else {
			if result.Code == Loading {
				c.loadNext()
			}
			return result.Id
		}
	}
}

func (c *cacheIdGenerator) GetIds(size int) []uint64 {
	var ids = []uint64{}
	for i := 0; i < size; i++ {
		ids = append(ids, c.GetId())
	}
	return ids
}
