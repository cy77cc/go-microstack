package utils

import (
	"encoding/json"
	"net"
	"path"
	"time"

	"github.com/google/uuid"
)

// GetFileExt 获取文件后缀名
func GetFileExt(fileName string) string {
	return path.Ext(fileName)
}

// GetTimestamp 获取当前时间戳
func GetTimestamp() int64 {
	return time.Now().Unix()
}

// GetMachineIP 获取本机IP
func GetMachineIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// GenUUID 生成 UUID
func GenUUID() string {
	return uuid.New().String()
}

// DeepCopy 深拷贝
func DeepCopy(src, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
