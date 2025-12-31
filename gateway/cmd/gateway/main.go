package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cy77cc/go-microstack/common/logx"
	"github.com/cy77cc/go-microstack/common/register"
	"github.com/cy77cc/go-microstack/common/register/types"
	"github.com/cy77cc/go-microstack/common/utils"
	"github.com/cy77cc/go-microstack/gateway/internal/config"
	"github.com/cy77cc/go-microstack/gateway/internal/proxy"
	"github.com/cy77cc/go-microstack/gateway/internal/router"
	"github.com/cy77cc/go-microstack/gateway/pkg/loadbalance"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. 初始化上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. 解析命令行参数
	var configPath string
	flag.StringVar(&configPath, "config", "etc/config.yaml", "config path")

	var gatewayConfigPath string
	flag.StringVar(&gatewayConfigPath, "gateway-router", "etc/gateway-router.json", "gateway config path")

	flag.Parse()

	// 3. 加载 .env 文件
	_ = godotenv.Load(".env")

	// 4. 初始化配置管理器
	configManager := config.NewConfigManager()

	// 5. 加载本地配置 (作为基础配置)
	localConfig, err := config.LoadLocalConfig(configPath)
	if err != nil {
		logx.Warnf("Failed to load local config: %v", err)
		// 如果本地配置加载失败，且没有配置注册中心，可能无法继续
	} else {
		configManager.SetLocalConfig(localConfig)
	}

	// 6. 加载本地路由配置
	routes, err := config.LoadRoutesFromJSON(gatewayConfigPath)
	if err != nil {
		logx.Warnf("Failed to load local routes: %v", err)
	} else {
		remoteConfig := &config.RemoteConfig{
			Routes: routes,
		}
		configManager.SetRemoteConfig(remoteConfig)
	}

	// 7. 连接注册中心并加载配置
	// 设置日志级别
	logx.SetLevel(logx.DEBUG)

	var regClient types.Register

	if localConfig != nil && localConfig.Register.Type != "" {

		localConfig.Register.Endpoints = append(localConfig.Register.Endpoints, os.Getenv("NACOS_ADDR"))
		localConfig.Register.Namespace = os.Getenv("NACOS_NAMESPACE")
		localConfig.Register.Username = os.Getenv("NACOS_USERNAME")
		localConfig.Register.Password = os.Getenv("NACOS_PASSWORD")

		regClient, err = register.NewRegister(ctx,
			localConfig.Register.Type,
			register.WithEndpoints(localConfig.Register.Endpoints...),
			register.WithAuth(localConfig.Register.Username, localConfig.Register.Password),
			register.WithNamespace(localConfig.Register.Namespace),
			register.WithTimeout(time.Duration(localConfig.Register.Timeout)*time.Millisecond),
		)

		if err == nil {
			logx.Infof("Connected to Register Center: %s", localConfig.Register.Type)

			// 加载全局配置
			item, err := regClient.GetConfig(ctx, "gateway", "DEFAULT_GROUP")
			if err != nil {
				logx.Errorf("Failed to load global config from register center: %v", err)
			} else {
				if err := configManager.ParseRemoteConfig(item); err != nil {
					logx.Errorf("Failed to parse global config from register center: %v", err)
				}
				logx.Info("Loaded global config from register center")
			}

			// 加载路由配置
			item, err = regClient.GetConfig(ctx, "gateway-router", "DEFAULT_GROUP")
			if err != nil {
				logx.Errorf("Failed to load router config from register center: %v", err)
			} else {
				if err := configManager.ParseRemoteConfig(item); err != nil {
					logx.Errorf("Failed to parse router config from register center: %v", err)
				}
				logx.Info("Loaded router config from register center")
			}

			// 注册 Gateway 服务自身
			go func() {
				// Wait for server to start
				time.Sleep(3 * time.Second)

				// Get IP
				ip := utils.GetMachineIP()
				port := localConfig.Server.Port

				inst := &types.ServiceInstance{
					ID:          fmt.Sprintf("%s-%s-%d", localConfig.Server.Name, ip, port),
					ServiceName: localConfig.Server.Name,
					Host:        ip,
					Port:        port,
					GroupName:   "DEFAULT_GROUP",
					ClusterName: "DEFAULT",
					Weight:      10,
					Metadata:    map[string]string{"http_port": strconv.Itoa(port)},
					HealthCheck: &types.HealthCheck{
						Type:     "tcp",
						URL:      fmt.Sprintf("%s:%d", ip, port),
						Interval: 10 * time.Second,
						Timeout:  2 * time.Second,
					},
				}

				if err := regClient.Register(context.Background(), inst); err != nil {
					logx.Errorf("Failed to register service: %v", err)
				} else {
					logx.Infof("Registered service %s to %s", localConfig.Server.Name, localConfig.Register.Type)
				}
			}()

		} else {
			logx.Errorf("Failed to connect to register center: %v", err)
		}
	}

	// 8. 初始化路由和代理
	// 获取初始配置
	currentConfig := configManager.GetConfig()

	// 创建负载均衡器
	lb := loadbalance.NewRoundRobinLoadBalancer()

	// 创建代理处理器
	proxyHandler := proxy.NewProxyHandler(regClient, lb)

	// 创建路由器
	r := router.NewRouter()

	// 注册为配置观察者
	configManager.RegisterWatcher(proxyHandler) // 如果代理也需要更新配置

	// 9. 启动 HTTP 服务器
	engine := gin.New()
	engine.Use(gin.Recovery())

	// 添加访问日志中间件
	engine.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		cfg := configManager.GetConfig()
		if cfg.Logging.AccessLog {
			// 简单的访问日志
			end := time.Now()
			latency := end.Sub(start)
			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()

			if raw != "" {
				path = path + "?" + raw
			}

			logx.Infof("[GIN] %v | %3d | %13v | %15s | %-7s %s",
				end.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		}
	})

	r.RegisterRoutes(engine, currentConfig.Routes, proxyHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", currentConfig.Server.Port),
		Handler: engine,
	}

	go func() {
		logx.Infof("Starting gateway server at %s:%d...", currentConfig.Server.Host, currentConfig.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logx.Fatalf("listen: %s\n", err)
		}
	}()

	// 10. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logx.Info("Shutting down server...")

	// 如果注册了服务，需要注销
	if regClient != nil && localConfig != nil {
		ip := utils.GetMachineIP()
		port := localConfig.Server.Port
		inst := &types.ServiceInstance{
			ID:          fmt.Sprintf("%s-%s-%d", localConfig.Server.Name, ip, port),
			ServiceName: localConfig.Server.Name,
			Host:        ip,
			Port:        port,
		}
		if err := regClient.Deregister(context.Background(), inst.ID); err != nil {
			logx.Errorf("Failed to deregister service: %v", err)
		} else {
			logx.Infof("Deregistered service %s", localConfig.Server.Name)
		}
	}

	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctx); err != nil {
		logx.Fatal("Server forced to shutdown: ", err)
	}

	logx.Info("Server exiting")
}
