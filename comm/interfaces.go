package comm

type IdGenerator interface {
	GetId(bizType string) uint64
	GetIds(bizType string) []uint64
}
