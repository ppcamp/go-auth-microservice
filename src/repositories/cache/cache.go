package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ppcamp/go-xtendlib/errors"
	log "github.com/sirupsen/logrus"
)

const (
	pingTime          time.Duration = 1 * time.Second
	CONNECTION_STRING string        = "%s:%s"
	ID_STRING_FORMAT  string        = "%s-%s"
)

type CacheConfig redis.Options

type Cache interface {
	Close() error

	Auth
	UserData
}

type cache struct {
	*redis.Client
	*auth
	*user
}

func NewCacheRepository(options CacheConfig, identifier string) (Cache, error) {
	ops := redis.Options(options)
	client := redis.NewClient(&ops)

	if resp := client.Ping(context.Background()); resp.Err() != nil {
		log.Warnf("redis connection failed %v", resp.Err())
		return nil, errors.Wraps("fail when ping the server", resp.Err())
	}

	return &cache{client, &auth{client, identifier}, &user{client, identifier}}, nil
}
