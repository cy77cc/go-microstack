package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		withSession(session sqlx.Session) UsersModel
		FindAll(ctx context.Context, query string, page, pageSize int64) ([]*Users, int64, error)
	}

	customUsersModel struct {
		*defaultUsersModel
		c cache.CacheConf
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c),
	}
}

func (m *customUsersModel) withSession(session sqlx.Session) UsersModel {
	return NewUsersModel(sqlx.NewSqlConnFromSession(session), m.c)
}

func (m *customUsersModel) FindAll(ctx context.Context, q string, page, pageSize int64) ([]*Users, int64, error) {
	var users []*Users
	var count int64

	offset := (page - 1) * pageSize
	where := ""
	var args []interface{}
	if q != "" {
		where = "WHERE `username` LIKE ?"
		args = append(args, "%"+q+"%")
	}

	countQuery := fmt.Sprintf("SELECT count(*) FROM %s %s", m.table, where)
	err := m.defaultUsersModel.QueryRowNoCacheCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s LIMIT ?, ?", usersRows, m.table, where)
	args = append(args, offset, pageSize)
	err = m.defaultUsersModel.QueryRowsNoCacheCtx(ctx, &users, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}
