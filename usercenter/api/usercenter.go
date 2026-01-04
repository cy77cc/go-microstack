package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/cy77cc/go-microstack/common/pkg/middleware"
	"github.com/cy77cc/go-microstack/common/pkg/register"
	"github.com/cy77cc/go-microstack/common/pkg/register/types"
	"github.com/cy77cc/go-microstack/common/pkg/utils"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/config"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/handler"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/usercenter.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	if c.Register.Type != "" {
		reg, err := register.NewRegister(context.Background(), c.Register.Type,
			register.WithEndpoints(c.Register.Endpoints...),
			register.WithAuth(c.Register.Username, c.Register.Password),
			register.WithNamespace(c.Register.Namespace),
			register.WithTimeout(time.Duration(c.Register.Timeout)*time.Millisecond),
		)
		if err == nil {
			// 1. Load config
			// if item, err := reg.GetConfig(context.Background(), c.Name, "DEFAULT_GROUP"); err == nil {
			// 	if err := yaml.Unmarshal([]byte(item.Value), &c); err != nil {
			// 		fmt.Printf("Failed to unmarshal remote config: %v\n", err)
			// 	} else {
			// 		fmt.Println("Loaded config from register center")
			// 	}
			// }

			// 2. Discover RPC services (Manual discovery and injection)
			// This allows us to use any registry supported by our package
			// provided we update the endpoints before creating ServiceContext
			if len(c.UserCenterRpc.Endpoints) == 0 {
				insts, err := reg.GetService(context.Background(), "usercenter.rpc", "DEFAULT_GROUP")
				if err == nil && len(insts) > 0 {
					var endpoints []string
					for _, inst := range insts {
						endpoints = append(endpoints, fmt.Sprintf("%s:%d", inst.Host, inst.Port))
					}
					c.UserCenterRpc.Endpoints = endpoints
					fmt.Printf("Discovered usercenter.rpc endpoints: %v\n", endpoints)
				}
			}

			// 3. Register API service (Async)
			go func() {
				time.Sleep(3 * time.Second)
				port := c.Port
				ip := utils.GetMachineIP()

				inst := &types.ServiceInstance{
					ID:          fmt.Sprintf("%s-%s-%d", c.Name, ip, port),
					ServiceName: c.Name,
					Host:        ip,
					Port:        port,
					GroupName:   "DEFAULT_GROUP",
					ClusterName: "DEFAULT",
					Weight:      10,
					Metadata:    map[string]string{"http_port": strconv.Itoa(port)},
				}
				if err := reg.Register(context.Background(), inst); err != nil {
					fmt.Printf("Failed to register service: %v\n", err)
				} else {
					fmt.Printf("Registered service %s to %s\n", c.Name, c.Register.Type)
				}
			}()
		}
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// Add Global Middleware
	server.Use(middleware.NewAuditMiddleware().Handle)
	server.Use(middleware.NewMetricMiddleware().Handle)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
