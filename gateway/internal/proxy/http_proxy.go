package proxy

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cy77cc/go-microstack/common/pkg/register/types"
	"github.com/cy77cc/go-microstack/gateway/internal/config"
	"github.com/cy77cc/go-microstack/gateway/pkg/loadbalance"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Handler 代理处理器
type Handler struct {
	discovery    types.Register
	loadBalancer loadbalance.LoadBalancer
	cfg          *config.MergedConfig
}

// NewProxyHandler 创建新的代理处理器
func NewProxyHandler(d types.Register, lb loadbalance.LoadBalancer) *Handler {
	return &Handler{
		discovery:    d,
		loadBalancer: lb,
	}
}

// HandleGeneric 处理通用的 /api/:service/*path 请求
func (h *Handler) HandleGeneric(c *gin.Context) {
	service := c.Param("service")
	path := c.Param("path")
	h.proxy(c, service, path, "")
}

// HandleRoute 返回特定路由的处理器函数
func (h *Handler) HandleRoute(serviceName, stripPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.proxy(c, serviceName, c.Request.URL.Path, stripPrefix)
	}
}

// proxy 核心反向代理逻辑
func (h *Handler) proxy(c *gin.Context, serviceName, path, stripPrefix string) {
	// 检查是否是 WebSocket 升级请求
	if h.isWebSocketRequest(c.Request) {
		h.proxyWebSocket(c, serviceName)
		return
	}

	// 1. 服务发现
	instances, err := h.discovery.GetService(c, serviceName, "")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("service discovery error: %v", err)})
		c.Abort()
		return
	}

	if len(instances) == 0 {
		c.JSON(http.StatusBadGateway, gin.H{"error": "no available instances"})
		c.Abort()
		return
	}

	// 2. 负载均衡选择实例
	instance, err := h.loadBalancer.Select(instances)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "no available instances"})
		c.Abort()
		return
	}

	// 3. 构建目标 URL
	targetURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", instance.Host, instance.Port))

	// 4. 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// 5. 路径重写
		targetPath := path
		if stripPrefix != "" && strings.HasPrefix(targetPath, stripPrefix) {
			targetPath = strings.TrimPrefix(targetPath, stripPrefix)
			if !strings.HasPrefix(targetPath, "/") {
				targetPath = "/" + targetPath
			}
		}

		req.URL.Path = targetPath
		req.Host = targetURL.Host

		// 6. 注入安全请求头 (身份信息 & 签名)
		h.injectSecurityHeaders(req)
	}

	// 7. 错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if !c.Writer.Written() {
			c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("proxy error: %v", err)})
		}
	}

	// 8. 执行代理
	proxy.ServeHTTP(c.Writer, c.Request)
}

// 判断是否为 WebSocket 升级请求
func (h *Handler) isWebSocketRequest(req *http.Request) bool {
	connectionHeader := strings.ToLower(req.Header.Get("Connection"))
	upgradeHeader := strings.ToLower(req.Header.Get("Upgrade"))
	return strings.Contains(connectionHeader, "upgrade") && upgradeHeader == "websocket"
}

// WebSocket 代理实现
func (h *Handler) proxyWebSocket(c *gin.Context, serviceName string) {
	instances, err := h.discovery.GetService(c, serviceName, "")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("service discovery error: %v", err)})
		c.Abort()
		return
	}

	instance, err := h.loadBalancer.Select(instances)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "no available instances"})
		c.Abort()
		return
	}

	// 构建目标 URL
	targetURL := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", instance.Host, instance.Port),
		Path:   c.Request.URL.Path,
	}

	// 创建 WebSocket 连接
	dialer := websocket.DefaultDialer
	targetConn, resp, err := dialer.Dial(targetURL.String(), c.Request.Header)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to connect to backend: %v", err)})
		if resp != nil {
			resp.Body.Close()
		}
		c.Abort()
		return
	}
	defer targetConn.Close()

	// 升级客户端连接
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许跨域
		},
	}
	clientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to upgrade client connection: %v", err)})
		c.Abort()
		return
	}
	defer clientConn.Close()

	// 在两个连接之间转发消息
	errChan := make(chan error, 2)

	// 从客户端转发到后端服务
	go func() {
		defer close(errChan)
		for {
			messageType, message, err := clientConn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}

			if err := targetConn.WriteMessage(messageType, message); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// 从后端服务转发到客户端
	go func() {
		for {
			messageType, message, err := targetConn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}

			if err := clientConn.WriteMessage(messageType, message); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// 等待任一方向出现错误
	<-errChan
}

func (h *Handler) OnConfigChange(config *config.MergedConfig) {
	h.cfg = config
}

func (h *Handler) injectSecurityHeaders(req *http.Request) {
	// Request ID
	if req.Header.Get("X-Request-Id") == "" {
		req.Header.Set("X-Request-Id", uuid.NewString())
	}
	// Timestamp
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	req.Header.Set("X-Timestamp", ts)
	// User ID from JWT or header
	uidStr := h.extractUID(req)
	if uidStr == "" {
		uidStr = "0"
	}
	req.Header.Set("X-User-Id", uidStr)
	// Signature
	secret := ""
	if h.cfg != nil {
		secret = h.cfg.Gateway.Sign.Secret
	}
	if secret == "" {
		return
	}
	canonical := req.Method + "\n" + req.URL.Path + "\n" + uidStr + "\n" + ts
	sig := hmacSha256Hex([]byte(secret), []byte(canonical))
	req.Header.Set("X-Signature", sig)
}

func (h *Handler) extractUID(req *http.Request) string {
	// Prefer JWT Authorization: Bearer <token>
	if h.cfg != nil && h.cfg.Gateway.Auth.AccessSecret != "" {
		auth := req.Header.Get("Authorization")
		prefix := "Bearer "
		if strings.HasPrefix(auth, prefix) {
			tokenStr := strings.TrimPrefix(auth, prefix)
			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(h.cfg.Gateway.Auth.AccessSecret), nil
			})
			if token != nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if v, ok := claims["uid"]; ok {
						switch vv := v.(type) {
						case float64:
							return strconv.FormatInt(int64(vv), 10)
						case int64:
							return strconv.FormatInt(vv, 10)
						case string:
							return vv
						}
					}
				}
			}
		}
	}
	// Fallback to header
	if hv := req.Header.Get("X-User-Id"); hv != "" {
		return hv
	}
	return ""
}

func hmacSha256Hex(key, data []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
