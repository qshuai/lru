package lru

import (
	"bytes"
	"container/list"
	"fmt"
	"sync"
)

type Cache struct {
	mtx      sync.Mutex
	m        map[interface{}]*list.Element
	l        *list.List
	capacity int
}

type Entry struct {
	// the map key
	Key interface{}

	Item interface{}
}

var _ lrucache = &Cache{}
var _ Item = &Entry{}

func (entry *Entry) GetKey() interface{} {
	return entry.Key
}

func (lru *Cache) String() string {
	lru.mtx.Lock()
	defer lru.mtx.Unlock()

	lastEntryNum := len(lru.m) - 1
	curEntry := 0
	buf := bytes.NewBufferString("[")
	for k := range lru.m {
		buf.WriteString(fmt.Sprintf("%v", k))
		if curEntry < lastEntryNum {
			buf.WriteString(", ")
		}
		curEntry++
	}

	buf.WriteString("]")

	return fmt.Sprintf("<%d/%d>%s", len(lru.m), lru.capacity, buf.String())
}

func (lru *Cache) Existed(key string) bool {
	lru.mtx.Lock()
	_, ok := lru.m[key]
	lru.mtx.Unlock()

	return ok
}

func (lru *Cache) Get(key string) *Entry {
	lru.mtx.Lock()
	defer lru.mtx.Unlock()

	ele, ok := lru.m[key]
	if !ok {
		return nil
	}

	// update the node's sequence at the same time
	lru.l.MoveToFront(ele)

	return ele.Value.(*Entry)
}

func (lru *Cache) Add(value *Entry) {
	if lru.capacity == 0 {
		return
	}

	if len(lru.m)+1 > lru.capacity {
		ele := lru.l.Back()
		entry := ele.Value.(*Entry)
		delete(lru.m, entry.Key)

		ele.Value = value
		lru.l.MoveToFront(ele)
		lru.m[value.Key] = ele

		return
	}

	node := lru.l.PushFront(value)
	lru.m[value.Key] = node
}

func (lru *Cache) Remove(key string) {
	lru.mtx.Lock()
	defer lru.mtx.Unlock()

	ele, ok := lru.m[key]
	if !ok {
		return
	}

	delete(lru.m, key)
	lru.l.Remove(ele)
}

func (lru *Cache) Size() int {
	lru.mtx.Lock()
	defer lru.mtx.Unlock()

	return len(lru.m)
}

func New(capacity int) *Cache {
	return &Cache{
		m:        make(map[interface{}]*list.Element, capacity),
		l:        list.New(),
		capacity: capacity,
	}
}
