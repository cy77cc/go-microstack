package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MultipartPartModel = (*customMultipartPartModel)(nil)

type (
	// MultipartPartModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMultipartPartModel.
	MultipartPartModel interface {
		multipartPartModel
		WithSession(session sqlx.Session) MultipartPartModel
		DeleteByUploadId(ctx context.Context, uploadId string) error
		FindAllByUploadId(ctx context.Context, uploadId string) ([]*MultipartPart, error)
	}

	customMultipartPartModel struct {
		*defaultMultipartPartModel
		c cache.CacheConf
	}
)

// NewMultipartPartModel returns a model for the database table.
func NewMultipartPartModel(conn sqlx.SqlConn, c cache.CacheConf) MultipartPartModel {
	return &customMultipartPartModel{
		defaultMultipartPartModel: newMultipartPartModel(conn, c),
	}
}

func (m *customMultipartPartModel) WithSession(session sqlx.Session) MultipartPartModel {
	return NewMultipartPartModel(sqlx.NewSqlConnFromSession(session), m.c)
}

func (m *customMultipartPartModel) DeleteByUploadId(ctx context.Context, uploadId string) error {
	query := fmt.Sprintf("delete from %s where `upload_id` = ?", m.table)
	_, err := m.defaultMultipartPartModel.ExecCtx(
		ctx,
		func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
			return conn.ExecCtx(ctx, query, uploadId)
		},
		fmt.Sprintf("fileserver:uploadId:%s", uploadId),
	)
	return err
}

func (m *customMultipartPartModel) FindAllByUploadId(ctx context.Context, uploadId string) ([]*MultipartPart, error) {
	query := fmt.Sprintf("select %s from %s where `upload_id` = ? order by `part_number` asc", multipartPartRows, m.table)
	var parts []*MultipartPart
	err := m.defaultMultipartPartModel.QueryRowsNoCacheCtx(ctx, &parts, query, uploadId)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
