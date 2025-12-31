package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PermissionsModel = (*customPermissionsModel)(nil)

type (
	// PermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPermissionsModel.
	PermissionsModel interface {
		permissionsModel
		withSession(session sqlx.Session) PermissionsModel
		FindAll(ctx context.Context, page, pageSize int64) ([]*Permissions, int64, error)
	}

	customPermissionsModel struct {
		*defaultPermissionsModel
		c cache.CacheConf
	}
)

// NewPermissionsModel returns a model for the database table.
func NewPermissionsModel(conn sqlx.SqlConn, c cache.CacheConf) PermissionsModel {
	return &customPermissionsModel{
		defaultPermissionsModel: newPermissionsModel(conn, c),
	}
}

func (m *customPermissionsModel) withSession(session sqlx.Session) PermissionsModel {
	return NewPermissionsModel(sqlx.NewSqlConnFromSession(session), m.c)
}

func (m *customPermissionsModel) FindAll(ctx context.Context, page, pageSize int64) ([]*Permissions, int64, error) {
	var permissions []*Permissions
	var count int64
	offset := (page - 1) * pageSize
	countQuery := fmt.Sprintf("SELECT count(*) FROM %s", m.table)
	err := m.defaultPermissionsModel.QueryRowCtx(
		ctx,
		&count,
		"usercenter:all_permissions_count",
		func(ctx context.Context, conn sqlx.SqlConn, v any) error {
			return conn.QueryRowCtx(ctx, v, countQuery)
		},
	)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT %s FROM %s LIMIT ?, ?", permissionsRows, m.table)
	err = m.defaultPermissionsModel.QueryRowsNoCacheCtx(
		ctx,
		&permissions,
		query,
		offset,
		pageSize,
	)
	if err != nil {
		return nil, 0, err
	}
	return permissions, count, nil
}
