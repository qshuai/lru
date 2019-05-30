package lru

type lrucache interface {
	Existed(key string) bool
	Add(value *Entry)
	Remove(key string)
	Get(key string) *Entry
	Size() int
}

type Item interface {
	GetKey() interface{}
}
