# Go Backend Final Lab

这是学习实验的初始脚手架。HTTP 服务设计为独立运行，`/healthz` 只证明进程存活，不检查 MySQL/Redis。五个 Task API 全部占位为 `501 {"error":"not_implemented"}`；CRUD、事务、缓存、前端网络调用、Dockerfile 和业务测试由学习者完成。

Go 依赖已下载并校验，基础测试、构建和 HTTP 实测已通过；五个 Task API 均保持 501。Docker 镜像拉取仍超时，MySQL/Redis 容器尚未启动，详见 [交接报告](docs/SCAFFOLD_HANDOFF.md)。

## 环境与启动

使用 Go 1.22 兼容环境；当前验证环境和依赖版本见 [交接报告](docs/SCAFFOLD_HANDOFF.md)。Windows 用户在 WSL2 中运行 Go，使用启用 WSL Integration 的 Docker Desktop。Docker 只用于开发依赖，不影响初始 HTTP 服务启动。

必须从项目根目录（当前工作区的 `go-backend-final-lab/`）运行，以便找到 `web/index.html`：

```bash
cd go-backend-final-lab
go mod tidy
go run ./cmd/api
```

访问 <http://127.0.0.1:18080/healthz> 与 <http://127.0.0.1:18080/>。自定义端口可用 `HTTP_PORT=18081 go run ./cmd/api`；Go 不自动加载 `.env`。

## 开发依赖

所有 Compose 命令显式使用独立项目名，避免混淆真实项目：

```bash
docker compose -p go-backend-final-lab config -q
docker compose -p go-backend-final-lab up -d mysql redis
docker compose -p go-backend-final-lab ps
docker compose -p go-backend-final-lab logs --tail=50 mysql redis
docker compose -p go-backend-final-lab down
```

仅有 MySQL 8.0 和 Redis 7-alpine。当前 Go API 由 WSL 本机运行，后续使用 `127.0.0.1:13306` / `127.0.0.1:16379` 连接依赖。MySQL 数据库和用户均为 `final_lab`，配置示例见 [.env.example](.env.example)。可自行复制为 `.env` 配置 Compose，示例密码仅用于本地开发，不能用于生产。

MySQL 使用项目独立命名卷，普通 `down` 会保留数据；Redis 仅作缓存练习，不保存数据。B5 才由学习者增加 API 容器，届时容器之间使用 `mysql:3306` / `redis:6379`。

## 架构与文件约定

| 位置 | 职责与边界 |
| --- | --- |
| `cmd/api/main.go` | 启动入口、读取 HTTP_PORT |
| `internal/app/` | 组装路由、进程健康检查、基础 HTTP 测试 |
| `internal/task/handler.go` | HTTP 边界；目前只有 501 响应 |
| `internal/task/service.go` | 业务与事务职责占位；不依赖 Gin Context |
| `internal/task/repository.go` | 持久化职责占位；尚未连接数据库 |
| `internal/task/model.go` | Task/TaskEvent 空类型；模型由学习者设计 |
| `internal/cache/` | B3 文档占位；尚无 Redis Client |
| `web/` | 无构建步骤的静态页面；B4 接入请求 |
| `docs/` | 脚手架交接与验证记录 |

沿用现有领域目录，不增加通用框架或额外架构层。未来调用方向是 Handler → Service → Repository，缓存由 Service 协调；当前尚未接通这些业务层。Go 文件使用小写英文名称，按职责放置；测试与被测包同目录；文档使用 Markdown。根目录仅保留运行配置、依赖清单和项目说明，凭据不进入源码。

## 学习 TODO（均待学习者验收）

- [ ] B0：亲手操作 Linux、Compose、MySQL CLI 和 redis-cli，检查日志与端口。
- [ ] B1：定义模型、请求校验，连接 GORM/MySQL，实现五个 Task API 与分页。
- [ ] B2：PATCH 事务关联写入、故障回滚验证与一次 EXPLAIN。
- [ ] B3：go-redis、Cache Aside、TTL、缓存失效与故障回源。
- [ ] B4：原生 JavaScript 联调、浏览器 Network 观察与错误显示。
- [ ] B5：自己编写 Dockerfile，向 Compose 增加 API 服务并验证容器网络。
- [ ] B6：业务测试、race 微实验、Debug 记录与教练给出的未知微变更。

五个占位路由：`POST /tasks`、`GET /tasks/:id`、`GET /tasks`、`PATCH /tasks/:id`、`DELETE /tasks/:id`。501 是当前交付的预期行为，不能算作 CRUD 已完成。

原始文档保留在上一级工作区：[Codex脚手架开发任务书](../01_Codex脚手架开发任务书.md)、[学习执行手册](../02_学习执行手册.md)。单独搬移项目时请同时保留这些用户提供的文档。

## 安全习惯

`.env` 与凭据不提交，`.env.example` 只提供示例。仅操作 `go-backend-final-lab` Compose 项目，不操作 RepoLens/Gateway。不要误用 `docker compose down -v` 处理真实项目，它会删除数据卷。初始脚手架不提供公网部署或生产配置。
