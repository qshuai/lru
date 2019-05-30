package lru

import (
	"bytes"
	"container/list"
	"fmt"
	"sync"
)

type Cache struct {
	mtx         sync.Mutex
	m           map[interface{}]*list.Element
	l           *list.List
	capacity    int
	accessOrder bool // keep access order or insert order
}

var _ lrucache = &Cache{}

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

// Existed will not influence the entry order
func (lru *Cache) Existed(key interface{}) bool {
	lru.mtx.Lock()
	_, ok := lru.m[key]
	lru.mtx.Unlock()

	return ok
}

func (lru *Cache) Get(key interface{}) Item {
	lru.mtx.Lock()
	defer lru.mtx.Unlock()

	ele, ok := lru.m[key]
	if !ok {
		return nil
	}

	// update the node's sequence at the same time
	if lru.accessOrder {
		lru.l.MoveToFront(ele)
	}

	return ele.Value.(Item)
}

func (lru *Cache) Add(value Item) {
	if lru.capacity == 0 {
		return
	}

	if len(lru.m)+1 > lru.capacity {
		ele := lru.l.Back()
		entry := ele.Value.(Item)
		delete(lru.m, entry.GetKey())

		ele.Value = value
		if lru.accessOrder {
			lru.l.MoveToFront(ele)
		}
		lru.m[value.GetKey()] = ele

		return
	}

	node := lru.l.PushFront(value)
	lru.m[value.GetKey()] = node
}

func (lru *Cache) Remove(key interface{}) {
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

func New(capacity int, accessOrder bool) *Cache {
	return &Cache{
		m:           make(map[interface{}]*list.Element, capacity),
		l:           list.New(),
		capacity:    capacity,
		accessOrder: accessOrder,
	}
}
