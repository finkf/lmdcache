package main

import (
	"log"
	"sync"

	"github.com/golang/groupcache/lru"
)

type cache struct {
	cache *lru.Cache
	lock  sync.Mutex
}

func newCache(n int) *cache {
	return &cache{cache: lru.New(n)}
}

func (c *cache) get(k interface{}, f func(interface{}) (interface{}, error)) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if v, found := c.cache.Get(k); found {
		log.Printf("lookup %v [cached]", k)
		return v, nil
	}
	log.Printf("lookup %v", k)
	v, err := f(k)
	if err != nil {
		return nil, err
	}
	c.cache.Add(k, v)
	return v, nil
}
