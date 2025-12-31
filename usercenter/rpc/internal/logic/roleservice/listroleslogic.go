package roleservicelogic

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListRolesLogic) ListRoles(in *pb.ListRolesReq) (*pb.ListRolesResp, error) {
	roles, total, err := l.svcCtx.RolesModel.FindAll(l.ctx, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, err
	}

	var pbRoles []*pb.Role
	for _, role := range roles {
		pbRoles = append(pbRoles, &pb.Role{
			Id:          role.Id,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Status:      int32(role.Status),
			CreateTime:  role.CreateTime,
			UpdateTime:  role.UpdateTime,
		})
	}

	return &pb.ListRolesResp{
		Roles: pbRoles,
		Total: total,
	}, nil
}
