package cryptx

import (
	"fmt"

	"golang.org/x/crypto/scrypt"
)

// PasswordEncrypt 使用 scrypt 对密码进行加密
// salt: 盐值
// password: 原始密码
// 返回: 加密后的十六进制字符串
func PasswordEncrypt(salt, password string) string {
	// scrypt.Key(password, salt, N, r, p, keyLen)
	// N: CPU/memory cost parameter (must be a power of 2, e.g. 16384 or 32768)
	// r: block size parameter
	// p: parallelization parameter
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	return fmt.Sprintf("%x", string(dk))
}

// PasswordVerify 验证密码是否正确
// salt: 盐值
// password: 原始密码
// hashedPassword: 加密后的密码
// 返回: 是否匹配
func PasswordVerify(salt, password, hashedPassword string) bool {
	return PasswordEncrypt(salt, password) == hashedPassword
}
