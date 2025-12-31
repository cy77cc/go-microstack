package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MultipartUploadModel = (*customMultipartUploadModel)(nil)

type (
	// MultipartUploadModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMultipartUploadModel.
	MultipartUploadModel interface {
		multipartUploadModel
		WithSession(session sqlx.Session) MultipartUploadModel
		UpdateStatusByUploadId(ctx context.Context, uploadId string, status int64, completeTime int64) error
	}

	customMultipartUploadModel struct {
		*defaultMultipartUploadModel
		c cache.CacheConf
	}
)

// NewMultipartUploadModel returns a model for the database table.
func NewMultipartUploadModel(conn sqlx.SqlConn, c cache.CacheConf) MultipartUploadModel {
	return &customMultipartUploadModel{
		defaultMultipartUploadModel: newMultipartUploadModel(conn, c),
	}
}

func (m *customMultipartUploadModel) WithSession(session sqlx.Session) MultipartUploadModel {
	return NewMultipartUploadModel(sqlx.NewSqlConnFromSession(session), m.c)
}

func (m *customMultipartUploadModel) UpdateStatusByUploadId(ctx context.Context, uploadId string, status int64, completeTime int64) error {
	query := fmt.Sprintf("update %s set `status` = ?, `complete_time` = ? where `upload_id` = ?", m.table)
	_, err := m.defaultMultipartUploadModel.ExecCtx(
		ctx,
		func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
			return conn.ExecCtx(ctx, query, status, completeTime, uploadId)
		},
		fmt.Sprintf("fileserver:uploadId:%s", uploadId),
	)
	return err
}
