package authservicelogic

import (
	"context"
	"errors"

	"time"

	"github.com/cy77cc/go-microstack/common/cryptx"
	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	// 1. Find user by username
	user, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.Unauthenticated, "invalid username or password")
		}
		return nil, err
	}

	// 2. Check password
	password := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password)
	if password != user.PasswordHash {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	// 3. Fetch Roles
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

	// 4. Generate Token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, accessExpire, int64(user.Id))
	if err != nil {
		return nil, err
	}
	refreshToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, accessExpire*72, int64(user.Id)) // 3 days? Or just longer.
	if err != nil {
		return nil, err
	}

	// 5. Update LastLoginTime
	user.LastLoginTime = time.Now().Unix()
	l.svcCtx.UsersModel.Update(l.ctx, user)

	return &pb.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expires:      now + accessExpire,
		Uid:          user.Id,
		Roles:        roleCodes,
	}, nil
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
