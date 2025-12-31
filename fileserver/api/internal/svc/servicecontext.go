package svc

import (
	"github.com/cy77cc/go-microstack/fileserver/api/internal/config"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/middleware"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	FilesRpc            pb.FileserverClient
	SignatureMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.FileRpc)
	return &ServiceContext{
		Config:   c,
		FilesRpc: pb.NewFileserverClient(client.Conn()),
		SignatureMiddleware: middleware.NewSignatureMiddleware(middleware.SignatureConfig{
			Secret:  c.Sign.Secret,
			SkewSec: c.Sign.SkewSec,
		}).Handle,
	}
}
