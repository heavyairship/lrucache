package lrucache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	tests := map[string]struct {
		preload     func() *LRUCache
		key         string
		expectedOk  bool
		expectedVal string
	}{
		"missing key": {
			preload: func() *LRUCache {
				c := NewLRUCache(2)
				c.Write("2", "hi")
				return c
			},
			key:         "1",
			expectedOk:  false,
			expectedVal: "",
		},
		"present key": {
			preload: func() *LRUCache {
				c := NewLRUCache(2)
				c.Write("2", "hi")
				return c
			},
			key:         "2",
			expectedOk:  true,
			expectedVal: "hi",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cache := test.preload()
			val, ok := cache.Read(test.key)
			assert.Equal(t, test.expectedOk, ok)
			assert.Equal(t, test.expectedVal, val)
		})
	}
}

func TestWriteSimpleEviction(t *testing.T) {
	val, ok := "", false

	c := NewLRUCache(2)
	c.Write("1", "1")
	c.Read("1")
	c.Write("2", "2")
	c.Write("3", "3") // Evicts 1 => 1

	val, ok = c.Read("1")
	assert.Equal(t, false, ok)
	assert.Equal(t, "", val)

	val, ok = c.Read("3")
	assert.Equal(t, true, ok)
	assert.Equal(t, "3", val)

	val, ok = c.Read("2")
	assert.Equal(t, true, ok)
	assert.Equal(t, "2", val)
}

func TestWriteLRUEviction(t *testing.T) {
	val, ok := "", false

	c := NewLRUCache(2)
	c.Write("1", "1")
	c.Read("1")
	c.Write("2", "2")
	c.Write("1", "1") // Move 1 to the top of the LRU
	c.Write("3", "3") // Evicts 2 => 2

	val, ok = c.Read("1")
	assert.Equal(t, true, ok)
	assert.Equal(t, "1", val)

	val, ok = c.Read("3")
	assert.Equal(t, true, ok)
	assert.Equal(t, "3", val)

	val, ok = c.Read("2")
	assert.Equal(t, false, ok)
	assert.Equal(t, "", val)
}

func TestREadLRUEviction(t *testing.T) {
	val, ok := "", false

	c := NewLRUCache(2)
	c.Write("1", "1")
	c.Read("1")
	c.Write("2", "2")
	c.Read("1")       // Move 1 to the top of the LRU
	c.Write("3", "3") // Evicts 2 => 2

	val, ok = c.Read("1")
	assert.Equal(t, true, ok)
	assert.Equal(t, "1", val)

	val, ok = c.Read("3")
	assert.Equal(t, true, ok)
	assert.Equal(t, "3", val)

	val, ok = c.Read("2")
	assert.Equal(t, false, ok)
	assert.Equal(t, "", val)
}

func TestWriteSameKey(t *testing.T) {
	val, ok := "", false

	c := NewLRUCache(2)
	c.Write("1", "1")
	c.Write("2", "2'")

	// These overwrites shoud not cause any evictions
	c.Write("2", "2''")
	c.Write("2", "2'''")
	c.Write("1", "1'")

	val, ok = c.Read("1")
	assert.Equal(t, true, ok)
	assert.Equal(t, "1'", val)

	val, ok = c.Read("2")
	assert.Equal(t, true, ok)
	assert.Equal(t, "2'''", val)
}
