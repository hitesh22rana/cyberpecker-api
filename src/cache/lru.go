package cache

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type LRUCache struct {
	maxSize     int
	expiration  int
	cache       map[string]*list.Element
	expireTimes map[string]time.Time
	list        *list.List
	mu          sync.Mutex
	timer       *time.Timer
}

func NewLRUCache(maxSize int, expiration int) *LRUCache {
	if maxSize < 1 || expiration < 1 {
		log.Fatalln("maxSize and expiration should be positive integers")
	}

	return &LRUCache{
		maxSize:     maxSize,
		expiration:  expiration,
		cache:       make(map[string]*list.Element),
		expireTimes: make(map[string]time.Time),
		list:        list.New(),
	}
}

func (c *LRUCache) hasExpired(key string) bool {
	expireTime, ok := c.expireTimes[key]
	if !ok {
		return false
	}

	return time.Now().UTC().After(expireTime)
}

func (c *LRUCache) removeExpiredItems() {
	currentTime := time.Now().UTC()
	expiredKeys := make([]string, 0)

	for key, expireTime := range c.expireTimes {
		if currentTime.After(expireTime) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		c.remove(key)
	}

	c.startTimer()
}

func (c *LRUCache) remove(key string) {
	if ele, ok := c.cache[key]; ok {
		c.list.Remove(ele)
		delete(c.cache, key)
		delete(c.expireTimes, key)
	}
}

func (c *LRUCache) startTimer() {
	if c.timer == nil {
		c.timer = time.AfterFunc(time.Duration(c.expiration)*time.Second, c.removeExpiredItems)
	}
}

func (c *LRUCache) StopTimer() {
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
}

func (c *LRUCache) GetNews(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.cache[key]; ok {
		if c.hasExpired(key) {
			c.remove(key)
			return nil
		}
		c.list.MoveToFront(ele)
		return ele.Value
	}
	return nil
}

func (c *LRUCache) SetNews(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.cache[key]; ok {
		c.remove(key)
	}

	if len(c.cache) >= c.maxSize {
		oldestEle := c.list.Back()
		if oldestEle != nil {
			oldestKey := oldestEle.Value.(string)
			c.remove(oldestKey)
		}
	}

	expireTime := time.Now().UTC().Add(time.Duration(c.expiration) * time.Second)
	c.expireTimes[key] = expireTime
	c.cache[key] = c.list.PushFront(value)

	c.startTimer()
}
