package cache

import (
	"crypto/sha512"
	"encoding/hex"
	"time"
)

var listCache = make(map[string]cache)
var limiter = make(chan int, 1)

type cache struct {
	hash string
	data interface{}
	time time.Time
}

func NewCache(search string, data interface{}, duration int64) *cache {
	newCache := &cache{
		hash: makeHash(search),
		data: data,
		time: time.Now().Add(time.Duration(duration) * time.Second),
	}
	listCache[newCache.hash] = *newCache
	go verifyCaches()
	return newCache
}

func verifyCaches() {
	if len(limiter) > 0 {
		return
	}
	limiter <- 1
	for len(listCache) > 0 {
		for _, r := range listCache {
			if time.Now().After(r.time) {
				removeFromSlice(r)
			}
			time.Sleep(1 * time.Second)
		}
	}
	<-limiter

}

func Get(search string) interface{} {
	return listCache[makeHash(search)]
}

func removeFromSlice(cache cache) {
	delete(listCache, cache.hash)
}

func makeHash(data string) string {

	hasher := sha512.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)
	hashCode := hex.EncodeToString(hashBytes)
	return hashCode

}
