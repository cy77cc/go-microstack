package userservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/common/cryptx"
	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *pb.CreateUserReq) (*pb.UserResp, error) {
	_, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, in.Username)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}
	if !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}

	passwordHash := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password)

	newUser := &model.Users{
		Username:     in.Username,
		PasswordHash: passwordHash,
		Email:        in.Email,
		Phone:        in.Phone,
		Avatar:       in.Avatar,
		Status:       1,
	}

	res, err := l.svcCtx.UsersModel.Insert(l.ctx, newUser)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	newUser.Id = uint64(id)

	return &pb.UserResp{
		User: &pb.User{
			Id:         newUser.Id,
			Username:   newUser.Username,
			Email:      newUser.Email,
			Phone:      newUser.Phone,
			Avatar:     newUser.Avatar,
			Status:     int32(newUser.Status),
			CreateTime: newUser.CreateTime,
			UpdateTime: newUser.UpdateTime,
		},
	}, nil
}
