package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AuthRefreshTokensModel = (*customAuthRefreshTokensModel)(nil)

type (
	// AuthRefreshTokensModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuthRefreshTokensModel.
	AuthRefreshTokensModel interface {
		authRefreshTokensModel
		withSession(session sqlx.Session) AuthRefreshTokensModel
	}

	customAuthRefreshTokensModel struct {
		*defaultAuthRefreshTokensModel
		c cache.CacheConf
	}
)

// NewAuthRefreshTokensModel returns a model for the database table.
func NewAuthRefreshTokensModel(conn sqlx.SqlConn, c cache.CacheConf) AuthRefreshTokensModel {
	return &customAuthRefreshTokensModel{
		defaultAuthRefreshTokensModel: newAuthRefreshTokensModel(conn, c),
	}
}

func (m *customAuthRefreshTokensModel) withSession(session sqlx.Session) AuthRefreshTokensModel {
	return NewAuthRefreshTokensModel(sqlx.NewSqlConnFromSession(session), m.c)
}
