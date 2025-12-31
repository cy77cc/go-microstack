package users

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) ListUsers(req *types.UserListReq) (resp *types.UserListResp, err error) {
	res, err := l.svcCtx.UserService.ListUsers(l.ctx, &userservice.ListUsersReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Query:    req.Query,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.UserResp, 0, len(res.Users))
	for _, user := range res.Users {
		list = append(list, types.UserResp{
			Id:            user.Id,
			Username:      user.Username,
			Email:         user.Email,
			Phone:         user.Phone,
			Avatar:        user.Avatar,
			Status:        user.Status,
			CreateTime:    user.CreateTime,
			UpdateTime:    user.UpdateTime,
			LastLoginTime: user.LastLoginTime,
		})
	}

	return &types.UserListResp{
		Total: res.Total,
		List:  list,
	}, nil
}
