package lru

type lrucache interface {
	Existed(key interface{}) bool
	Add(value Item) error
	Remove(key interface{})
	Get(key interface{}) Item
	Size() int
	Iterate() error
}

type Item interface {
	GetKey() interface{}
}
