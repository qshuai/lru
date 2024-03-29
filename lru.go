package lru

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type Cache struct {
	sync.RWMutex

	m           map[interface{}]*list.Element
	l           *list.List
	capacity    int
	accessOrder bool // keep access order or insert order

	Callback func(item Item) error
	ForEach  func(item Item) error
}

var CallbackNotFound = errors.New("the necessary callback function not found")

var _ lrucache = &Cache{}

func (lru *Cache) String() string {
	lru.RLock()
	defer lru.RUnlock()

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
	lru.RLock()
	_, ok := lru.m[key]
	lru.RUnlock()

	return ok
}

func (lru *Cache) Get(key interface{}) Item {
	lru.RLock()
	defer lru.RUnlock()

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

func (lru *Cache) Add(value Item) error {
	lru.Lock()
	defer lru.Unlock()

	if lru.capacity == 0 {
		if lru.Callback == nil {
			return CallbackNotFound
		}

		return lru.Callback(value)
	}

	if len(lru.m)+1 > lru.capacity {
		ele := lru.l.Back()
		entry := ele.Value.(Item)
		delete(lru.m, entry.GetKey())

		// execution the callback function firstly, to avoid
		// disrupt the origin element value
		if lru.Callback != nil {
			err := lru.Callback(entry)
			if err != nil {
				return err
			}
		}

		ele.Value = value
		if lru.accessOrder {
			lru.l.MoveToFront(ele)
		}
		lru.m[value.GetKey()] = ele

		return nil
	}

	node := lru.l.PushFront(value)
	lru.m[value.GetKey()] = node

	return nil
}

func (lru *Cache) Remove(key interface{}) {
	lru.Lock()
	defer lru.Unlock()

	ele, ok := lru.m[key]
	if !ok {
		return
	}

	delete(lru.m, key)
	lru.l.Remove(ele)
}

func (lru *Cache) Size() int {
	lru.RLock()
	defer lru.RUnlock()

	return lru.l.Len()
}

func (lru *Cache) Iterate() error {
	lru.RLock()
	defer lru.RUnlock()

	for _, v := range lru.m {
		err := lru.ForEach(v.Value.(Item))
		if err != nil {
			return err
		}
	}

	return nil
}

func New(capacity int, accessOrder bool) *Cache {
	return &Cache{
		m:           make(map[interface{}]*list.Element, capacity),
		l:           list.New(),
		capacity:    capacity,
		accessOrder: accessOrder,
	}
}
