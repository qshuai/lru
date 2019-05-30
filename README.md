lru
===

##### Usage:

```go
package main

import (
	"fmt"

	"github.com/qshuai/lru"
)

type user struct {
	key string

	score int
}

func (u *user) GetKey() interface{} {
	return u.key
}

func main() {
  // set capacity to 3 for this cache as a sample
  cache := lru.New(3, true)
  // add 4 entries, the first entry will be removed automaticlly
  cache.Add(&user{"Andy", 23})
  cache.Add(&user{"Tom", 24})
  cache.Add(&user{"Angel", 25})
  cache.Add(&user{"Nike", 26})
  fmt.Println(cache.Size(), cache.Existed("Andy"), cache.Existed("Tom"),
    cache.Existed("Angel"), cache.Existed("Nike"))

  // access the user Tom, the tom will be move to front
  cache.Get("Tom")
  // now add a entry, the tail entry user named Angel will be removed
  cache.Add(&user{"Jone", 27})
  fmt.Println(cache.Size(), cache.Existed("Tom"), cache.Existed("Angel"),
    cache.Existed("Nike"), cache.Existed("Jone"))

  // remove a enry by hand, we will discard this and update cache size
  cache.Remove("Jone")
  fmt.Println(cache.Size(), cache.Existed("Tom"), cache.Existed("Angel"),
    cache.Existed("Nike"), cache.Existed("Jone"))
}
```

