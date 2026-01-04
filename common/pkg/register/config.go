package register

// Config 注册中心配置
type Config struct {
	Type      string   `json:"type" yaml:"type"`
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
	Username  string   `json:"username" yaml:"username"`
	Password  string   `json:"password" yaml:"password"`
	Namespace string   `json:"namespace" yaml:"namespace"`
	Timeout   int64    `json:"timeout" yaml:"timeout"` // milliseconds
}
