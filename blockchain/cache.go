package blockchain

import (
	"sync"
)

type cacheable[K any, V any] struct {
	cache sync.Map
}

func (c *cacheable[K, V]) getOrQueryFunc(key K, queryFunc func(key K) (V, error)) (V, error) {
	if cached, ok := c.cache.Load(key); ok {
		return cached.(V), nil
	}

	val, err := queryFunc(key)
	if err != nil {
		return val, err
	}

	c.cache.Store(key, val)

	return val, nil
}
