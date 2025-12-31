package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cy77cc/go-microstack/common/register"
	"github.com/cy77cc/go-microstack/common/register/types"
	"github.com/cy77cc/go-microstack/common/utils"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/config"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/server"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/fileserver.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config

	if err := conf.Load(*configFile, &c, conf.UseEnv()); err != nil {
		logx.Errorf("Failed to load config: %v\n", err)
	}

	if c.Register.Type != "" {
		reg, err := register.NewRegister(context.Background(), c.Register.Type,
			register.WithEndpoints(c.Register.Endpoints...),
			register.WithAuth(c.Register.Username, c.Register.Password),
			register.WithNamespace(c.Register.Namespace),
			register.WithTimeout(time.Duration(c.Register.Timeout)*time.Millisecond),
		)
		if err == nil {
			// // 1. Load config
			// if item, err := reg.GetConfig(context.Background(), "fileserver", "DEFAULT_GROUP"); err == nil {
			// 	if err := yaml.Unmarshal([]byte(item.Value), &c); err != nil {
			// 		logx.Errorf("Failed to unmarshal remote config: %v\n", err)
			// 	} else {
			// 		logx.Infof("Loaded config from register center")
			// 	}
			// }

			// 2. Register service (Async)
			go func() {
				time.Sleep(3 * time.Second)
				parts := strings.Split(c.ListenOn, ":")
				port := 8080
				if len(parts) > 1 {
					p, _ := strconv.Atoi(parts[1])
					port = p
				}
				ip := utils.GetMachineIP()

				inst := &types.ServiceInstance{
					ID:          fmt.Sprintf("%s-%s-%d", c.Name, ip, port),
					ServiceName: c.Name,
					Host:        ip,
					Port:        port,
					GroupName:   "DEFAULT_GROUP",
					ClusterName: "DEFAULT",
					Weight:      10,
					Metadata:    map[string]string{"gRPC_port": strconv.Itoa(port)},
				}
				if err := reg.Register(context.Background(), inst); err != nil {
					logx.Errorf("Failed to register service: %v\n", err)
				} else {
					logx.Infof("Registered service %s to %s\n", c.Name, c.Register.Type)
				}
			}()
		}
	}

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterFileserverServer(grpcServer, server.NewFileserverServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
