package main

import (
	"fmt"

	"github.com/lrucache/lrucache"
)

type Cache interface {
	Write(k, v string)
	Read(k string) (string, bool)
	Print()
}

func test(c Cache) {

}
func main() {
	c := lrucache.NewLRUCache(10)
	for i := 0; i < 100; i++ {
		c.Write(fmt.Sprintf("%d", i), fmt.Sprintf("%d", i))
		c.Read("0")
	}
	c.Read("94")
	c.Print() // Should have 94, 0, 99, 98, 97, 96, 95, 93, 92, 91
	c.Write("100", "100")
	c.Print() // Should have 100, 94, 0, 99, 98, 97, 96, 95, 93, 92
}
