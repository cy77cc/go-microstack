package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

// Instance Nacos实例结构
type Instance struct {
	clientConfig *constant.ClientConfig
	serverConfig []constant.ServerConfig
	ConfigClient config_client.IConfigClient
	NamingClient naming_client.INamingClient
}

// Config Nacos配置
type Config struct {
	Endpoints   []string `yaml:"endpoints"`
	Port        uint64   `yaml:"port"`
	Namespace   string   `yaml:"namespace"`
	Group       string   `yaml:"group"`
	TimeoutMs   uint64   `yaml:"timeout_ms"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	ContextPath string   `yaml:"context_path"`
	IdentityKey string   `yaml:"identity_key"`
	IdentityVal string   `yaml:"identity_val"`
	Token       string   `yaml:"token"`
}
