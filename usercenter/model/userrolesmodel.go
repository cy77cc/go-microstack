package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserRolesModel = (*customUserRolesModel)(nil)

type (
	// UserRolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserRolesModel.
	UserRolesModel interface {
		userRolesModel
		withSession(session sqlx.Session) UserRolesModel
		FindAllByUserId(ctx context.Context, userId uint64) ([]*UserRoles, error)
		DeleteByRoleId(ctx context.Context, roleId uint64) error
	}

	customUserRolesModel struct {
		*defaultUserRolesModel
		c cache.CacheConf
	}
)

// NewUserRolesModel returns a model for the database table.
func NewUserRolesModel(conn sqlx.SqlConn, c cache.CacheConf) UserRolesModel {
	return &customUserRolesModel{
		defaultUserRolesModel: newUserRolesModel(conn, c),
	}
}

func (m *customUserRolesModel) withSession(session sqlx.Session) UserRolesModel {
	return NewUserRolesModel(sqlx.NewSqlConnFromSession(session), m.c)
}

func (m *customUserRolesModel) FindAllByUserId(ctx context.Context, userId uint64) ([]*UserRoles, error) {
	var userRoles []*UserRoles
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `user_id` = ?", userRolesRows, m.table)
	err := m.defaultUserRolesModel.QueryRowsNoCacheCtx(ctx, &userRoles, query, userId)
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (m *customUserRolesModel) DeleteByRoleId(ctx context.Context, roleId uint64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE `role_id` = ?", m.table)
	_, err := m.defaultUserRolesModel.ExecCtx(
		ctx,
		func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
			return conn.ExecCtx(ctx, query, roleId)
		},
		fmt.Sprintf("usercenter:user_roles:%d", roleId),
	)
	return err
}
