package lru

import (
	"bytes"
	"container/list"
	"fmt"
	"strconv"
	"sync"
)

type Cache struct {
	mtx      sync.Mutex
	m        map[string]*list.Element
	l        *list.List
	capacity int
}

type Entry struct {
	// the map key
	key string

	num int
}

func (entry *Entry) String() string {
	return strconv.Itoa(entry.num)
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

	return fmt.Sprintf("<%d>%s", lru.capacity, buf.String())
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

func (lru *Cache) Add(key string, value *Entry) {
	if lru.capacity == 0 {
		return
	}

	if len(lru.m)+1 > lru.capacity {
		ele := lru.l.Back()
		entry := ele.Value.(*Entry)
		delete(lru.m, entry.key)

		ele.Value = value
		lru.l.MoveToFront(ele)
		lru.m[key] = ele

		return
	}

	node := lru.l.PushFront(value)
	lru.m[key] = node
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
		m:        make(map[string]*list.Element, capacity),
		l:        list.New(),
		capacity: capacity,
	}
}
