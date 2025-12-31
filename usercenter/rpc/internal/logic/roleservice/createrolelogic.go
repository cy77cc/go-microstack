package roleservicelogic

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

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRoleLogic) CreateRole(in *pb.CreateRoleReq) (*pb.RoleResp, error) {
	_, err := l.svcCtx.RolesModel.FindOneByCode(l.ctx, in.Code)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "role code already exists")
	}
	if !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}

	newRole := &model.Roles{
		Name:        in.Name,
		Code:        in.Code,
		Description: in.Description,
		Status:      1,
	}

	res, err := l.svcCtx.RolesModel.Insert(l.ctx, newRole)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	newRole.Id = uint64(id)

	return &pb.RoleResp{
		Role: &pb.Role{
			Id:          newRole.Id,
			Name:        newRole.Name,
			Code:        newRole.Code,
			Description: newRole.Description,
			Status:      int32(newRole.Status),
			CreateTime:  newRole.CreateTime,
			UpdateTime:  newRole.UpdateTime,
		},
	}, nil
}
