package authservicelogic

import (
	"context"
	"errors"

	"time"

	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *pb.RefreshTokenReq) (*pb.LoginResp, error) {
	// Verify refresh token
	token, err := jwt.Parse(in.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.JwtAuth.AccessSecret), nil
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	userId := int64(claims["userId"].(float64))

	// Find user
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(userId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.Unauthenticated, "user not found")
		}
		return nil, err
	}

	// Fetch Roles
	userRoles, err := l.svcCtx.UserRolesModel.FindAllByUserId(l.ctx, user.Id)
	if err != nil {
		return nil, err
	}
	var roleCodes []string
	for _, ur := range userRoles {
		role, _ := l.svcCtx.RolesModel.FindOne(l.ctx, ur.RoleId)
		if role != nil {
			roleCodes = append(roleCodes, role.Code)
		}
	}

	// Generate New Tokens
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, accessExpire, int64(user.Id))
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, accessExpire*72, int64(user.Id))
	if err != nil {
		return nil, err
	}

	return &pb.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		Expires:      now + accessExpire,
		Uid:          user.Id,
		Roles:        roleCodes,
	}, nil
}

func (l *RefreshTokenLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
