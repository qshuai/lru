package lru

import (
	"fmt"
)

func ExampleLruRemove() {
	lru := New(3)
	lru.Add(&Entry{Key: "1", Item: 1})
	lru.Add(&Entry{Key: "2", Item: 2})
	lru.Add(&Entry{Key: "3", Item: 3})
	lru.Add(&Entry{Key: "4", Item: 4})

	fmt.Println(lru.Size(), lru.Existed("1"), lru.Existed("2"),
		lru.Existed("3"), lru.Existed("4"))

	// Output:
	// 3 false true true true
}

func ExampleLruGet() {
	lru := New(3)
	lru.Add(&Entry{Key: "1", Item: 1})
	lru.Add(&Entry{Key: "2", Item: 2})
	lru.Add(&Entry{Key: "3", Item: 3})
	lru.Get("1")
	lru.Add(&Entry{Key: "4", Item: 4})

	fmt.Println(lru.Size(), lru.Existed("1"), lru.Existed("2"),
		lru.Existed("3"), lru.Existed("4"))

	// Output:
	// 3 true false true true
}
