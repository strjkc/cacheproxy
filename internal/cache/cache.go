package cache

import "time"

type Cache []map[string]CacheEntry

type CacheEntry struct {
	Reply   []byte
	Fetched time.Time
}

func NewCache() *Cache {
	shards := make(Cache, 0, 0)
	for i := 0; i < 32; i++ {
		shards = append(shards, make(map[string]CacheEntry, 0))
	}
	return &shards
}

func (c *Cache) Cache(item CacheEntry, key string) {
	item, ok := c.GetItem(key)
	if ok {
		c.UpdateCache(item, key)
	}
	//TODO: create a hash of key and get shard index based of of that
	(*c)[0][key] = item

}

func (c *Cache) GetItem(key string) (CacheEntry, bool) {
	item, ok := (*c)[0][key]
	if ok {
		return item, true
	}
	return CacheEntry{}, false
}

func (c *Cache) UpdateCache(item CacheEntry, key string) {
	(*c)[0][key] = item
}

func (c *Cache) RemoveCacheItem(key string) {

}

func (c *Cache) IsCacheItemExpired(key string) bool {
	item, ok := c.GetItem(key)
	if ok {
		if time.Since(item.Fetched) > time.Duration(time.Second*60) {
			return true
		}
		return false
	}
	return true
}
