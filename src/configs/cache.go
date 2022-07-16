package configs

import "time"

var (
	CacheHost     string
	CachePort     string
	CachePassword string
	CacheDb       int
)

const (
	CACHE_CONNECTION_TIMEOUT time.Duration = time.Second * 2
)
