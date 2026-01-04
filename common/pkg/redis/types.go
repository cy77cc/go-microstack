package redis

import (
	"crypto/tls"
	"time"
)

type TYPE string

const (
	CLUSTER    TYPE = "cluster"
	SENTINEL   TYPE = "sentinel"
	STANDALONE TYPE = "standalone"
)

type CommonOptions struct {
	// 基础
	Addrs    []string `yaml:"addrs" json:"addrs"`
	Username string   `yaml:"username" json:"username"`
	Password string   `yaml:"password" json:"password"`
	DB       int      `yaml:"db" json:"db"`

	// 超时控制
	DialTimeout  time.Duration `yaml:"dial_timeout" json:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`

	// 连接池
	PoolSize     int           `yaml:"pool_size" json:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns" json:"min_idle_conns"`
	PoolTimeout  time.Duration `yaml:"pool_timeout" json:"pool_timeout"`

	// 重试
	MaxRetries      int           `yaml:"max_retries" json:"max_retries"`
	MinRetryBackoff time.Duration `yaml:"min_retry_backoff" json:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `yaml:"max_retry_backoff" json:"max_retry_backoff"`

	// TLS
	TLS *tls.Config `yaml:"-" json:"-"`
}

type SentinelOptions struct {
	MasterName string `yaml:"master_name" json:"master_name"`
}

type ClusterOptions struct {
	RouteByLatency bool `yaml:"route_by_latency" json:"route_by_latency"`
	RouteRandomly  bool `yaml:"route_randomly" json:"route_randomly"`
}

type Config struct {
	// 是集群还是哨兵还是普通的
	Type     TYPE            `yaml:"type" json:"type"`
	Common   CommonOptions   `yaml:"common" json:"common"`
	Sentinel SentinelOptions `yaml:"sentinel" json:"sentinel"`
	Cluster  ClusterOptions  `yaml:"cluster" json:"cluster"`
}
