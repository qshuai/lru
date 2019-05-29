package lru

import (
	"fmt"
)

func ExampleLruRemove() {
	lru := New(3)
	lru.Add("1", &Entry{key: "1", num: 1})
	lru.Add("2", &Entry{key: "2", num: 2})
	lru.Add("3", &Entry{key: "3", num: 3})
	lru.Add("4", &Entry{key: "4", num: 4})

	fmt.Println(lru.Size(), lru.Existed("1"), lru.Existed("2"),
		lru.Existed("3"), lru.Existed("4"))

	// Output:
	// 3 false true true true
}

func ExampleLruGet() {
	lru := New(3)
	lru.Add("1", &Entry{key: "1", num: 1})
	lru.Add("2", &Entry{key: "2", num: 2})
	lru.Add("3", &Entry{key: "3", num: 3})
	lru.Get("1")
	lru.Add("4", &Entry{key: "4", num: 4})

	fmt.Println(lru.Size(), lru.Existed("1"), lru.Existed("2"),
		lru.Existed("3"), lru.Existed("4"))

	// Output:
	// 3 true false true true
}
