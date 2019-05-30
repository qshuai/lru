package lru

type lrucache interface {
	Existed(key string) bool
	Add(value Item)
	Remove(key string)
	Get(key string) Item
	Size() int
}

type Item interface {
	GetKey() interface{}
}
