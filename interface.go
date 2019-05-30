package lru

type lrucache interface {
	Existed(key interface{}) bool
	Add(value Item)
	Remove(key interface{})
	Get(key interface{}) Item
	Size() int
}

type Item interface {
	GetKey() interface{}
}
