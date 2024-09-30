package pokecache

import (
	"log/slog"
	"sync"
	"time"
)

const cullInterval = 5

type PokeCache struct {
	Cache  map[string]cacheEntry
	ticker *time.Ticker
	mutex  *sync.Mutex
}

type cacheEntry struct {
	val       []byte
	timestamp time.Time
}

func NewCache(cullIntervalInSeconds int) PokeCache {
	newCache := PokeCache{
		Cache:  map[string]cacheEntry{},
		ticker: time.NewTicker(time.Duration(cullInterval) * time.Second),
		mutex:  &sync.Mutex{},
	}
	go newCache.timerTrigger()
	return newCache
}

func (pCache *PokeCache) timerTrigger() {
	for tick := range pCache.ticker.C {
		tllTime := tick.Add(-cullInterval * time.Second)
		pCache.reapLoop(tllTime)
	}
}

func (pCache *PokeCache) Get(key string) ([]byte, bool) {
	pCache.mutex.Lock()
	defer pCache.mutex.Unlock()
	data, found := pCache.Cache[key]
	if !found {
		return nil, false
	}
	slog.Debug("Pokecache - Cache hit: " + key)
	return data.val, true
}

func (pCache *PokeCache) Add(key string, val []byte) {
	pCache.mutex.Lock()
	defer pCache.mutex.Unlock()
	pCache.Cache[key] = cacheEntry{
		val:       val,
		timestamp: time.Now(),
	}
	slog.Debug("Pokecache - Added to cache: " + key)
}

func (pCache *PokeCache) reapLoop(ttl time.Time) {
	pCache.mutex.Lock()
	defer pCache.mutex.Unlock()

	for key, entry := range pCache.Cache {
		if entry.timestamp.Before(ttl) {
			delete(pCache.Cache, key)
		}
	}
}
