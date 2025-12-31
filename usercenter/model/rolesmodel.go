package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolesModel = (*customRolesModel)(nil)

type (
	// RolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolesModel.
	RolesModel interface {
		rolesModel
		withSession(session sqlx.Session) RolesModel
		FindAll(ctx context.Context, page, pageSize int64) ([]*Roles, int64, error)
	}

	customRolesModel struct {
		*defaultRolesModel
		c cache.CacheConf
	}
)

// NewRolesModel returns a model for the database table.
func NewRolesModel(conn sqlx.SqlConn, c cache.CacheConf) RolesModel {
	return &customRolesModel{
		defaultRolesModel: newRolesModel(conn, c),
	}
}

func (m *customRolesModel) withSession(session sqlx.Session) RolesModel {
	return NewRolesModel(sqlx.NewSqlConnFromSession(session), m.c)
}

func (m *customRolesModel) FindAll(ctx context.Context, page, pageSize int64) ([]*Roles, int64, error) {
	var roles []*Roles
	var count int64
	offset := (page - 1) * pageSize
	countQuery := fmt.Sprintf("SELECT count(*) FROM %s", m.table)
	err := m.defaultRolesModel.QueryRowNoCacheCtx(ctx, &count, countQuery)
	if err != nil {
		return nil, 0, err
	}
	query := fmt.Sprintf("SELECT %s FROM %s LIMIT ?, ?", rolesRows, m.table)
	err = m.defaultRolesModel.QueryRowsNoCacheCtx(ctx, &roles, query, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}
