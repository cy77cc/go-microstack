package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BucketConfigModel = (*customBucketConfigModel)(nil)

type (
	// BucketConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBucketConfigModel.
	BucketConfigModel interface {
		bucketConfigModel
		WithSession(session sqlx.Session) BucketConfigModel
	}

	customBucketConfigModel struct {
		*defaultBucketConfigModel
		c cache.CacheConf
	}
)

// NewBucketConfigModel returns a model for the database table.
func NewBucketConfigModel(conn sqlx.SqlConn, c cache.CacheConf) BucketConfigModel {
	return &customBucketConfigModel{
		defaultBucketConfigModel: newBucketConfigModel(conn, c),
	}
}

func (m *customBucketConfigModel) WithSession(session sqlx.Session) BucketConfigModel {
	return NewBucketConfigModel(sqlx.NewSqlConnFromSession(session), m.c)
}
