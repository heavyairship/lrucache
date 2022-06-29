package lrucache

import "fmt"

type entry struct {
	k string
	v string
}

type node struct {
	prev  *node
	next  *node
	entry *entry
}

type dLinkedList struct {
	start *node
	end   *node
}

func (l *dLinkedList) push(k, v string) *node {
	n := &node{
		entry: &entry{
			k: k,
			v: v,
		},
	}
	if l.start == nil {
		l.end = n
	} else {
		l.start.prev = n
		n.next = l.start
	}
	l.start = n
	return n
}

func (l *dLinkedList) moveFront(n *node) {
	if n == l.start {
		return
	}
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
	if n == l.end {
		l.end = n.prev
	}
	n.prev = nil
	n.next = l.start
	l.start.prev = n
	l.start = n
}

type LRUCache struct {
	list  dLinkedList
	cache map[string]*node
	size  uint
}

func NewLRUCache(size uint) *LRUCache {
	return &LRUCache{
		cache: make(map[string]*node, size),
		size:  size,
	}
}

func (l *LRUCache) evict() {
	// Size 0 case
	if l.list.end == nil {
		return
	}

	n := l.list.end
	delete(l.cache, n.entry.k)

	// Size 1 case
	if l.list.start == l.list.end {
		l.list.start = nil
		l.list.end = nil
		return
	}

	// Size > 1 case
	n.prev.next = nil
	l.list.end = n.prev
	n.prev = nil
}

func (l *LRUCache) Write(k, v string) {
	if uint(len(l.cache)) == l.size {
		l.evict()
	}
	l.cache[k] = l.list.push(k, v)
}

func (l *LRUCache) Read(k string) (string, bool) {
	if n, ok := l.cache[k]; ok {
		l.list.moveFront(n)
		return n.entry.v, ok
	}
	return "", false
}

func (l *LRUCache) Print() {
	fmt.Print("cache:\n  ")
	for _, n := range l.cache {
		fmt.Printf("%v=>%v, ", n.entry.k, n.entry.v)
	}
	fmt.Print("\nlist:\n  ")
	curr := l.list.start
	for curr != nil {
		fmt.Printf("%v=>%v, ", curr.entry.k, curr.entry.v)
		curr = curr.next
	}
	fmt.Println()
}
