package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func MustNewRedisClient(cfg *Config) redis.UniversalClient {
	switch cfg.Type {

	case STANDALONE:
		return redis.NewClient(&redis.Options{
			Addr:            cfg.Common.Addrs[0],
			Username:        cfg.Common.Username,
			Password:        cfg.Common.Password,
			DB:              cfg.Common.DB,
			DialTimeout:     cfg.Common.DialTimeout,
			ReadTimeout:     cfg.Common.ReadTimeout,
			WriteTimeout:    cfg.Common.WriteTimeout,
			PoolSize:        cfg.Common.PoolSize,
			MinIdleConns:    cfg.Common.MinIdleConns,
			PoolTimeout:     cfg.Common.PoolTimeout,
			MaxRetries:      cfg.Common.MaxRetries,
			MinRetryBackoff: cfg.Common.MinRetryBackoff,
			MaxRetryBackoff: cfg.Common.MaxRetryBackoff,
			TLSConfig:       cfg.Common.TLS,
		})

	case SENTINEL:
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:      cfg.Sentinel.MasterName,
			SentinelAddrs:   cfg.Common.Addrs,
			Username:        cfg.Common.Username,
			Password:        cfg.Common.Password,
			DB:              cfg.Common.DB,
			DialTimeout:     cfg.Common.DialTimeout,
			ReadTimeout:     cfg.Common.ReadTimeout,
			WriteTimeout:    cfg.Common.WriteTimeout,
			PoolSize:        cfg.Common.PoolSize,
			MinIdleConns:    cfg.Common.MinIdleConns,
			PoolTimeout:     cfg.Common.PoolTimeout,
			MaxRetries:      cfg.Common.MaxRetries,
			MinRetryBackoff: cfg.Common.MinRetryBackoff,
			MaxRetryBackoff: cfg.Common.MaxRetryBackoff,
			TLSConfig:       cfg.Common.TLS,
			ReplicaOnly:     true,
		})

	case CLUSTER:
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:           cfg.Common.Addrs,
			Username:        cfg.Common.Username,
			Password:        cfg.Common.Password,
			DialTimeout:     cfg.Common.DialTimeout,
			ReadTimeout:     cfg.Common.ReadTimeout,
			WriteTimeout:    cfg.Common.WriteTimeout,
			PoolSize:        cfg.Common.PoolSize,
			MinIdleConns:    cfg.Common.MinIdleConns,
			PoolTimeout:     cfg.Common.PoolTimeout,
			MaxRetries:      cfg.Common.MaxRetries,
			MinRetryBackoff: cfg.Common.MinRetryBackoff,
			MaxRetryBackoff: cfg.Common.MaxRetryBackoff,
			RouteByLatency:  cfg.Cluster.RouteByLatency,
			RouteRandomly:   cfg.Cluster.RouteRandomly,
			TLSConfig:       cfg.Common.TLS,
		})

	default:
		return nil
	}
}

func DefaultCommonOptions() CommonOptions {
	return CommonOptions{
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolSize:        20,
		MinIdleConns:    5,
		PoolTimeout:     4 * time.Second,
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	}
}
