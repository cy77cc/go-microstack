package users

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.UserCreateReq) (resp *types.UserResp, err error) {
	res, err := l.svcCtx.UserService.CreateUser(l.ctx, &userservice.CreateUserReq{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserResp{
		Id:            res.User.Id,
		Username:      res.User.Username,
		Email:         res.User.Email,
		Phone:         res.User.Phone,
		Avatar:        res.User.Avatar,
		Status:        res.User.Status,
		CreateTime:    res.User.CreateTime,
		UpdateTime:    res.User.UpdateTime,
		LastLoginTime: res.User.LastLoginTime,
	}, nil
}
