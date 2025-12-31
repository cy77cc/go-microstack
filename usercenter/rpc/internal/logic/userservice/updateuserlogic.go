package userservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.UserResp, error) {
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, err
	}

	if in.Email != "" {
		user.Email = in.Email
	}
	if in.Phone != "" {
		user.Phone = in.Phone
	}
	if in.Avatar != "" {
		user.Avatar = in.Avatar
	}
	if in.Status != 0 {
		user.Status = int64(in.Status)
	}

	err = l.svcCtx.UsersModel.Update(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.UserResp{
		User: &pb.User{
			Id:            user.Id,
			Username:      user.Username,
			Email:         user.Email,
			Phone:         user.Phone,
			Avatar:        user.Avatar,
			Status:        int32(user.Status),
			CreateTime:    user.CreateTime,
			UpdateTime:    user.UpdateTime,
			LastLoginTime: user.LastLoginTime,
		},
	}, nil
}
