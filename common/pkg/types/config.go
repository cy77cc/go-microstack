package types

type Mysql struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Endpoint string `yaml:"endpoint" json:"endpoint,optional"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
}

type Minio struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Endpoint string `yaml:"endpoint" json:"endpoint,optional"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	SSL      bool   `yaml:"ssl" json:"ssl,optional"`
}
