package userservicelogic

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUsersLogic) ListUsers(in *pb.ListUsersReq) (*pb.ListUsersResp, error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 {
		in.PageSize = 10
	}

	users, total, err := l.svcCtx.UsersModel.FindAll(l.ctx, in.Query, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, err
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:            user.Id,
			Username:      user.Username,
			Email:         user.Email,
			Phone:         user.Phone,
			Avatar:        user.Avatar,
			Status:        int32(user.Status),
			CreateTime:    user.CreateTime,
			UpdateTime:    user.UpdateTime,
			LastLoginTime: user.LastLoginTime,
		})
	}

	return &pb.ListUsersResp{
		Users: pbUsers,
		Total: total,
	}, nil
}
