package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FileInfoModel = (*customFileInfoModel)(nil)

type (
	// FileInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFileInfoModel.
	FileInfoModel interface {
		fileInfoModel
		WithSession(session sqlx.Session) FileInfoModel
		FindOneByHash(ctx context.Context, hash string) (*FileInfo, error)
	}

	customFileInfoModel struct {
		*defaultFileInfoModel
		c cache.CacheConf
	}
)

// NewFileInfoModel returns a model for the database table.
func NewFileInfoModel(conn sqlx.SqlConn, c cache.CacheConf) FileInfoModel {
	return &customFileInfoModel{
		defaultFileInfoModel: newFileInfoModel(conn, c),
	}
}

func (m *customFileInfoModel) WithSession(session sqlx.Session) FileInfoModel {
	return NewFileInfoModel(sqlx.NewSqlConnFromSession(session), m.c)
}

func (m *customFileInfoModel) FindOneByHash(ctx context.Context, hash string) (*FileInfo, error) {
	var resp FileInfo
	query := fmt.Sprintf("select %s from %s where `hash` = ? and `status` = 1 limit 1", fileInfoRows, m.table)
	err := m.defaultFileInfoModel.QueryRowCtx(ctx,
		&resp,
		fmt.Sprintf("fileserver:file:%s", hash),
		func(ctx context.Context, conn sqlx.SqlConn, v any) error {
			return conn.QueryRowCtx(ctx, v, query, hash)
		})
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
