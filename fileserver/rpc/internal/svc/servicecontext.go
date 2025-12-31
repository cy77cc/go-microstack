package svc

import (
	"fmt"

	"github.com/cy77cc/go-microstack/fileserver/model"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	Conn   sqlx.SqlConn
	// models
	BucketModel   model.BucketConfigModel
	FileModel     model.FileInfoModel
	UploadModel   model.MultipartUploadModel
	UploadPartMod model.MultipartPartModel
	// storage
	Storage StorageRouter
	// tools
	Tools *Tools

	Rdb *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=10s",
		c.Mysql.Username, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Database))

	c.Mysql.Endpoint = fmt.Sprintf("%s:%d", c.Mysql.Host, c.Mysql.Port)
	c.Minio.Endpoint = fmt.Sprintf("%s:%d", c.Minio.Host, c.Minio.Port)

	rdb := redis.MustNewRedis(c.Redis.RedisConf)

	ctx := &ServiceContext{
		Config:        c,
		Conn:          conn,
		BucketModel:   model.NewBucketConfigModel(conn, c.CacheRedis),
		FileModel:     model.NewFileInfoModel(conn, c.CacheRedis),
		UploadModel:   model.NewMultipartUploadModel(conn, c.CacheRedis),
		UploadPartMod: model.NewMultipartPartModel(conn, c.CacheRedis),
		Tools:         NewTools(c.Upload.MaxFileSize, c.Upload.AllowedExts),
		Rdb:           rdb,
	}
	baseDir := c.Local.BaseDir
	if baseDir == "" {
		baseDir = "data"
	}
	router, err := NewStorageRouter(c, baseDir)
	if err != nil {
		logx.Errorf("init storage router error: %v", err)
	}
	if rr, ok := router.(*storageRouter); ok {
		rr.bm = ctx.BucketModel
	}
	ctx.Storage = router
	return ctx
}
