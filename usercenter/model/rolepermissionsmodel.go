package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolePermissionsModel = (*customRolePermissionsModel)(nil)

type (
	// RolePermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolePermissionsModel.
	RolePermissionsModel interface {
		rolePermissionsModel
		withSession(session sqlx.Session) RolePermissionsModel
		FindAllByRoleId(ctx context.Context, roleId uint64) ([]*RolePermissions, error)
		DeleteByRoleId(ctx context.Context, roleId uint64) error
	}

	customRolePermissionsModel struct {
		*defaultRolePermissionsModel
		c cache.CacheConf
	}
)

// NewRolePermissionsModel returns a model for the database table.
func NewRolePermissionsModel(conn sqlx.SqlConn, c cache.CacheConf) RolePermissionsModel {
	return &customRolePermissionsModel{
		defaultRolePermissionsModel: newRolePermissionsModel(conn, c),
	}
}

func (m *customRolePermissionsModel) withSession(session sqlx.Session) RolePermissionsModel {
	return NewRolePermissionsModel(sqlx.NewSqlConnFromSession(session), m.c)
}

func (m *customRolePermissionsModel) FindAllByRoleId(ctx context.Context, roleId uint64) ([]*RolePermissions, error) {
	var rolePermissions []*RolePermissions
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `role_id` = ?", rolePermissionsRows, m.table)
	err := m.defaultRolePermissionsModel.QueryRowsNoCacheCtx(ctx, &rolePermissions, query, roleId)
	if err != nil {
		return nil, err
	}
	return rolePermissions, nil
}

func (m *customRolePermissionsModel) DeleteByRoleId(ctx context.Context, roleId uint64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE `role_id` = ?", m.table)
	_, err := m.defaultRolePermissionsModel.ExecCtx(
		ctx,
		func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
			return conn.ExecCtx(ctx, query, roleId)
		},
		fmt.Sprintf("usercenter:role_permissions:%d", roleId),
	)
	return err
}
