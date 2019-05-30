package lru

import (
	"fmt"
)

type Entry struct {
	key string

	entry int
}

func (entry *Entry) GetKey() interface{} {
	return entry.key
}

func ExampleLruRemove() {
	lru := New(3)
	lru.Add(&Entry{key: "1", entry: 1})
	lru.Add(&Entry{key: "2", entry: 2})
	lru.Add(&Entry{key: "3", entry: 3})
	lru.Add(&Entry{key: "4", entry: 4})

	fmt.Println(lru.Size(), lru.Existed("1"), lru.Existed("2"),
		lru.Existed("3"), lru.Existed("4"))

	// Output:
	// 3 false true true true
}

func ExampleLruGet() {
	lru := New(3)
	lru.Add(&Entry{key: "1", entry: 1})
	lru.Add(&Entry{key: "2", entry: 2})
	lru.Add(&Entry{key: "3", entry: 3})
	lru.Get("1")
	lru.Add(&Entry{key: "4", entry: 4})

	fmt.Println(lru.Size(), lru.Existed("1"), lru.Existed("2"),
		lru.Existed("3"), lru.Existed("4"))

	// Output:
	// 3 true false true true
}
