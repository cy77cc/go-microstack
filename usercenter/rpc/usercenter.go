package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cy77cc/go-microstack/common/pkg/register"
	"github.com/cy77cc/go-microstack/common/pkg/register/types"
	"github.com/cy77cc/go-microstack/common/pkg/utils"
	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/config"
	authserviceServer "github.com/cy77cc/go-microstack/usercenter/rpc/internal/server/authservice"
	permissionserviceServer "github.com/cy77cc/go-microstack/usercenter/rpc/internal/server/permissionservice"
	roleserviceServer "github.com/cy77cc/go-microstack/usercenter/rpc/internal/server/roleservice"
	userserviceServer "github.com/cy77cc/go-microstack/usercenter/rpc/internal/server/userservice"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/usercenter.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	// Register to Nacos/Etcd/Consul if configured
	if c.Register.Type != "" {
		reg, err := register.NewRegister(context.Background(), c.Register.Type,
			register.WithEndpoints(c.Register.Endpoints...),
			register.WithAuth(c.Register.Username, c.Register.Password),
			register.WithNamespace(c.Register.Namespace),
			register.WithTimeout(time.Duration(c.Register.Timeout)*time.Millisecond),
		)
		if err == nil {
			// 1. Load remote config
			// if item, err := reg.GetConfig(context.Background(), c.Name, "DEFAULT_GROUP"); err == nil {
			// 	if err := yaml.Unmarshal([]byte(item.Value), &c); err != nil {
			// 		fmt.Printf("Failed to unmarshal remote config: %v\n", err)
			// 	} else {
			// 		fmt.Println("Loaded config from register center")
			// 	}
			// } else {
			// 	fmt.Printf("Failed to get config from register center: %v\n", err)
			// }

			// 2. Register service (Async)
			go func() {
				// Wait for server to start
				time.Sleep(3 * time.Second)

				// Parse port
				parts := strings.Split(c.ListenOn, ":")
				port := 8080
				if len(parts) > 1 {
					p, _ := strconv.Atoi(parts[1])
					port = p
				}

				// Get IP
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
					fmt.Printf("Failed to register service: %v\n", err)
				} else {
					fmt.Printf("Registered service %s to %s\n", c.Name, c.Register.Type)
				}
			}()
		} else {
			fmt.Printf("Failed to create register client: %v\n", err)
		}
	}

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterAuthServiceServer(grpcServer, authserviceServer.NewAuthServiceServer(ctx))
		pb.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(ctx))
		pb.RegisterRoleServiceServer(grpcServer, roleserviceServer.NewRoleServiceServer(ctx))
		pb.RegisterPermissionServiceServer(grpcServer, permissionserviceServer.NewPermissionServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(xcode.Interceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
