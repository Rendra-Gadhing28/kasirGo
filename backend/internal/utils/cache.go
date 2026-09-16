package utils

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var MemoryCache = cache.New(5*time.Minute, 10*time.Minute)

func SetCache(key string, value interface{}, duration time.Duration) {
	MemoryCache.Set(key, value, duration)
}

func GetCache(key string) (interface{}, bool) {
	return MemoryCache.Get(key)
}

func DeleteCache(key string) {
	MemoryCache.Delete(key)
}

func FlushCache() {
	MemoryCache.Flush()
}

// ComputeETag creates an MD5 hash string for HTTP ETag
func ComputeETag(data []byte) string {
	return fmt.Sprintf(`"%x"`, md5.Sum(data))
}
