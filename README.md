# Go Microstack

## 项目简介
Go Microstack 是一个基于 [Go Zero](https://go-zero.dev/) 微服务框架构建的后端系统。它旨在提供一套标准、高效、可扩展的微服务基础组件，包含用户中心、文件服务等核心模块。

## 模块列表

| 模块名称 | 目录 | 描述 |
| :--- | :--- | :--- |
| **User Center** | [usercenter](./usercenter) | 用户管理、认证鉴权 (JWT)、RBAC 权限控制 |
| **File Server** | [fileserver](./fileserver) | 文件上传下载、分片上传、对接 MinIO |

## 技术栈

*   **Golang**: 核心开发语言
*   **Go Zero**: 微服务框架 (RPC, API, Model)
*   **gRPC / Protobuf**: 服务间通信
*   **MySQL**: 关系型数据库
*   **Redis**: 缓存与会话管理
*   **Etcd**: 服务注册与发现
*   **MinIO**: 对象存储服务

## 快速开始

### 环境准备
1.  安装 Go 1.18+
2.  安装 `goctl` 工具: `go install github.com/zeromicro/go-zero/tools/goctl@latest`
3.  准备依赖服务: MySQL, Redis, Etcd, MinIO (推荐使用 Docker Compose 部署)

### 部署流程

1.  **初始化数据库**
    *   执行 `usercenter/model/*.sql`
    *   执行 `fileserver/model/*.sql`

2.  **启动 User Center**
    *   进入 `usercenter/rpc` 修改配置并启动 `go run usercenter.go`
    *   进入 `usercenter/api` 修改配置并启动 `go run usercenter.go`

3.  **启动 File Server**
    *   进入 `fileserver/rpc` 修改配置并启动 `go run fileserver.go`
    *   进入 `fileserver/api` 修改配置并启动 `go run fileserver.go`

## 可观测性接入指南 (Observability)

本项目已集成统一的审计日志与指标监控体系。

### 1. 审计日志 (Audit Log)
所有 API 请求会自动记录审计日志，包含 TraceID, UserID, Method, Path, Status, Duration, ClientIP 等信息。
日志默认输出到标准输出 (stdout)，格式为 JSON，可通过日志收集系统 (如 ELK/Loki) 采集 `[AUDIT]` 关键字。

### 2. 监控指标 (Metrics)
各服务已暴露 Prometheus 格式的监控指标，可直接对接 Prometheus/Grafana。

| 服务名称 | Metrics 端点 | 描述 |
| :--- | :--- | :--- |
| **User Center** | `http://<host>:41003/metrics` | 用户中心 API 指标 |
| **File Server** | `http://<host>:41004/metrics` | 文件服务 API 指标 |
| **Gateway** | `http://<host>:41000/metrics` | 网关层指标 (Gin) |

**指标说明**:
*   `microstack_requests_duration_ms_bucket`: 请求耗时分布 (Histogram)
*   `microstack_requests_code_total`: 请求状态码计数 (Counter)
*   `gateway_requests_duration_ms_bucket`: 网关请求耗时 (Histogram)
*   `gateway_requests_code_total`: 网关请求状态码 (Counter)

### 3. 链路追踪 (Tracing)
基于 Go Zero 的 OpenTelemetry 支持。需在 `etc/*.yaml` 中配置 `Telemetry` 节点对接 Jaeger/Zipkin。

## 开发规范

*   **API 定义**: 使用 `.api` 文件定义 HTTP 接口，通过 `goctl api go` 生成代码。
*   **RPC 定义**: 使用 `.proto` 文件定义 RPC 接口，通过 `goctl rpc protoc` 生成代码。
*   **错误处理**: 统一使用 `xerr` 包（需自定义）或标准错误码返回。
*   **文档**: 每个模块维护自己的 README，API 文档由 Swagger 自动生成（需配置）。

## 联系方式
*   Author: zhangdongping


**总体目标**

- 以“应用交付”和“平台治理”为核心，覆盖代码到生产的全链路，兼顾安全、可观测性、成本与治理。
- 提供低门槛的可视化体验：DAG 流程图、向导式操作、命令面板与快捷键，具备多语言与主题定制。
- 支持 K8s/Compose 双栈的一键部署、实时反馈与回滚，统一发布与环境策略。

**流水线与交付**

- 可视化流水线编排：支持串并行、条件分支、审批节点、手动 Gate、重试/回退策略。
- GitOps 集成：Argo CD/Flux 支持，自动化从 Git 推送到环境，带变更审计。
- 渐进式交付：蓝绿、金丝雀、分阶段发布与暂停/继续，带指标驱动（SLO/SLA）自动推进或回滚。
- 供应链安全：构建产物签名（Cosign）、SBOM 生成与验证（CycloneDX/SPDX）、镜像与依赖漏洞扫描（Trivy/Grype）。
- 多环境推广：Dev/Staging/Prod 推广流，与冻结窗口、变更门禁、风险分级审批联动。

**部署与环境管理**

- 应用模板中心：Helm、Kustomize、Compose 与 Terraform 模板统一管理；一键生成项目骨架。
- 环境编排与隔离：命名空间/项目/团队三级治理、配额与限流、优先级队列、资源池（GPU/ARM）。
- 策略与准入：OPA/Gatekeeper 策略（镜像白名单、网络策略、资源上限）、Admission Webhook。
- 配置与密钥：集中配置（Nacos/Consul/Etcd）、密钥管理（Vault/SOPS/Sealed Secrets），动态热更新与滚动注入。
- 灰度与流量治理：API 网关、路由与权重、熔断与重试、限速与配额（与现有 gateway 模块联动）。

**可观测性与运维**

- 指标/日志/链路追踪三件套：Prometheus + Grafana、Loki/ELK、OpenTelemetry/Jaeger/Tempo。
- SLO 管理与误差预算：服务级 SLO 定义、预算消耗看板、告警与发布门禁联动。
- 事件与变更审计：K8s 事件、发布记录、配置变更、审批轨迹统一归档与查询。
- 运行手册与自动化修复：Runbook 体系、故障自动化（重启/扩容/降级）、值班与升级路径。
- 真实用户监控与合成监控：RUM 探针与 Synthetic 检查，确保端到端体验。

**安全与合规**

- 身份与访问：SSO（OIDC/SAML）、MFA、细粒度 RBAC/ABAC、临时提权与双人审批、Just-in-Time 访问。
- 合规与审计：CIS/K8s Benchmark 检查、合规扫描（依赖/镜像/基础设施）、数据主权与审计日志留存策略。
- 供应链与镜像治理：私有镜像仓库策略、签名校验、漏洞与过期基线、强制重建/替换。
- 机密数据治理：密钥轮换、最小权限、机密访问审计；避免敏感信息进入仓库与日志（当前 .env 明文需整治）。

**平台治理与多租户**

- 多租户与项目空间：团队/项目分层、资源配额、隔离策略、越权检测。
- 服务目录与依赖映射：API Catalog、依赖图谱、版本与兼容矩阵、弃用告警。
- 成本与容量（FinOps）：资源成本归集与分摊、预算与告警、容量趋势与预测、自动扩缩容（HPA/VPA/KEDA）。
- 风险与发布窗口：变更风险评估、冻结窗口治理、强制门禁与回滚策略。

**开发者体验与 UI/交互**

- 主题与布局：自定义主题、深浅色、密度与布局切换，多语言（i18n）与区域化（日期/数值格式）。
- 快捷导航与命令面板：全局搜索、命令面板、键盘快捷键、上下文操作与工作流提示。
- 实时反馈：WebSocket/Server-Sent Events 推送，发布进度、日志流、指标就地显示。
- 指导式体验：向导与教程（Tour）、就地帮助（Tooltip/Popover）、错误自解释与修复建议。
- 性能优化：虚拟列表、增量加载、离线缓存、前后端协同分页；确保操作流畅。

**文档与帮助系统**

- 文档中心与知识库：快速入门、最佳实践、故障排查、蓝图模板与示例库。
- API 文档与规范：OpenAPI/AsyncAPI 中心，SDK/CLI 文档、版本与变更日志、弃用计划。
- 上下文帮助：界面元素级帮助、推荐下一步动作、模板和策略的预设注解。

**集成与扩展能力**

- 插件与扩展：流水线步骤插件、表单扩展、视图与图表小部件（Widget），事件总线触发（Webhooks/Kafka）。
- ChatOps：Slack/Teams 集成，工单与审批在 IM 里完成，机器人提醒与交互发布。
- 外部系统：工单/ITSM、CMDB、代码托管（GitLab/GitHub/Gitea）、制品库（Nexus/Artifactory）。

**灾备与韧性**

- 备份与恢复：应用/数据/配置备份，演练机制与恢复目标（RPO/RTO）。
- 多活与跨区域：读写分离/容灾切换策略、DNS 流量分配、线路健康检查。
- 混沌工程：注入故障与演练、韧性评估与改进计划。

**与当前代码库的落地建议**

- usercenter：扩展为平台的统一身份与授权中心（SSO/OIDC、MFA、细粒度 RBAC、审计日志）。
- gateway：接入限流、熔断、重试、灰度路由；可作为统一入口与发布门禁落地点。
- fileserver：接入日志归档与工件存储（构建产物/报告），支持生命周期与合规策略。
- common：新增审计日志、策略校验、告警与通知 SDK；整合 OpenTelemetry。
- deploy 目录：提供 Helm/Kustomize/Compose 与 Terraform 的模板中心（蓝图），支持一键部署与环境选择。
- 安全整改：将 [.env](file:///e:/project/demo/go-microstack/.env) 中所有敏感参数迁移至 Vault/SOPS/Nacos（加密），避免明文与仓库暴露；启用密钥轮换与访问审计。

**优先级路线图（建议）**

- 第1阶段（平台基础）：身份与权限中心、模板与一键部署、可视化流水线（含审批/回滚）、统一日志/指标接入。
- 第2阶段（治理与SLO）：GitOps、SLO/预算管理、发布门禁策略、供应链安全与镜像签名、审计与合规。
- 第3阶段（体验与扩展）：命令面板与快捷键、多语言与主题、ChatOps、插件市场、知识库与引导式体验。
- 第4阶段（韧性与FinOps）：灾备演练、混沌工程、多活容灾、成本治理与容量预测。

如果需要，我可以基于当前仓库先实现“统一审计日志与指标上报”最小可用版本，并为 usercenter/gateway/fileserver 增加健康指标、追踪与权限策略的骨架，随后补充 Helm 模板与一键部署流程。
