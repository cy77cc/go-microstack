package jwtx

import "github.com/golang-jwt/jwt/v4"

// GetToken 生成 JWT Token
// secretKey: 密钥
// iat: 签发时间戳 (秒)
// seconds: 过期时间 (秒)
// uid: 用户ID
// 返回: token字符串, 错误信息
func GetToken(secretKey string, iat, seconds int64, uid uint64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["uid"] = uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

// GetTokenWithClaims 生成带自定义载荷的 JWT Token
func GetTokenWithClaims(secretKey string, iat, seconds int64, payloads map[string]interface{}) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
