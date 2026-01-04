package authservicelogic

import (
	"context"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"
	"github.com/golang-jwt/jwt/v4"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyTokenLogic {
	return &VerifyTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyTokenLogic) VerifyToken(in *pb.VerifyTokenReq) (*pb.VerifyTokenResp, error) {

	get, err := l.svcCtx.Rdb.Get("usercenter:balcklist:" + in.AccessToken)
	if err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
	}

	if get != "" {
		return nil, xcode.NewErrCodeMsg(xcode.Unauthorized, "invalid access token")
	}

	token, err := jwt.Parse(in.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.JwtAuth.AccessSecret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	uid := uint64(claims["user_id"].(float64))
	if !ok || !token.Valid {
		return nil, xcode.NewErrCodeMsg(xcode.TokenInvalid, "invalid refresh token")
	}

	_, err = l.svcCtx.UsersModel.FindOne(l.ctx, uid)

	if err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.UserNotExist, "user not exist")
	}

	roles, err := l.svcCtx.UserRolesModel.FindAllByUserId(l.ctx, uid)

	if err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
	}

	var roleCodes []string

	for _, ur := range roles {
		role, _ := l.svcCtx.RolesModel.FindOne(l.ctx, ur.RoleId)
		roleCodes = append(roleCodes, role.Code)
	}

	return &pb.VerifyTokenResp{
		Uid:   uint64(claims["user_id"].(float64)),
		Valid: true,
		Roles: roleCodes,
	}, nil
}
