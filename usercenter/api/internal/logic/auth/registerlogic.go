// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/authservice"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.UserCreateReq) (resp *types.TokenResp, err error) {
	// 1. 调用 RPC 创建用户
	_, err = l.svcCtx.UserService.CreateUser(l.ctx, &userservice.CreateUserReq{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	// 2. 注册成功后自动登录，生成 Token
	res, err := l.svcCtx.AuthService.Login(l.ctx, &authservice.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &types.TokenResp{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		Expires:      res.Expires,
		Uid:          res.Uid,
		Roles:        res.Roles,
	}, nil
}
