# Go Backend Final Lab v2.0 — Codex 脚手架开发任务书

> **接收方：Codex（可直接将本文全文作为开发指令）**  
> **交付目标：搭建供本人学习的可运行脚手架，不是实现完整任务管理应用。**  
> 仓库建议名：`go-backend-final-lab`  
> 对应学习文档：`02_学习执行手册.md`  
> 依据：`Go_Backend_Final_Lab_v2_Frozen_2026-10-14.md`（B 阶段范围已冻结）  
> 日期：2026-10-14

---

## 0. 最高优先级指令：不要替学习者完成实验

你是**脚手架工程师**，不是业务需求开发工程师。成功标准不是「所有 Task API 可用」，而是：

1. 一个新的、隔离的 Go 项目能编译、能启动、`GET /healthz` 返回 200。
2. MySQL 与 Redis 的**开发依赖容器**能通过基础 `compose.yaml` 启动。
3. 五个 Task API 路由**已经占位，但统一返回 501 Not Implemented**；真实业务未实现。
4. 留下能编译的极简 `Handler / Service / Repository` 代码骨架与 TODO，但**不完成**数据库和缓存操作。
5. 留下一个可打开的极简静态页面框架，但**不编写**调用 Task API 的 `fetch` 逻辑。
6. 有明确交接文档说明已交付内容、尚待学习者实现的内容及最小验证结果。

**严禁“好心顺便做完”**：真实 CRUD、GORM 数据操作、事务、EXPLAIN、Redis Client/Cache Aside、Dockerfile、完整 Go API Compose 服务、前端 fetch、业务测试、竞态测试、陌生需求变更，这些全部必须留给学习者。

如果本文指令与通常的软件工程“尽量交付完整产品”习惯冲突，**优先保留学习空缺**。不得借口“使测试通过”实现业务逻辑。

---

## 1. 安全边界与操作范围

- 仅在用户指定的**全新独立目录** `go-backend-final-lab` 中操作。若当前工作目录不是目标目录，或目录存在非脚手架文件，**先汇报，不覆盖、不删除、不重置**。
- **不得修改或向** `mrussss/repolens`、`mrussss/gateway-system` **提交任何内容**；它们是后续 C 阶段只读参考项目。
- 不扫描、读取、复制现有项目的 `.env`、API Key、数据库密码、Provider 凭证。只生成本地开发**示例**变量；`.env` 在 `.gitignore` 中。
- 不购买云服务器，不执行公网部署；MySQL/Redis 对宿主机的端口仅绑定 `127.0.0.1`。
- 不创建 JWT、RBAC、登录、Agent、MQ、gRPC、K8s、微服务、CI/CD、监控平台或 React/Vue 工程。
- 如需安装 Go 依赖，选择与 Go 1.22 兼容的稳定版本。优先沿用项目已存在的 Go 环境；不升级系统工具链。
- 如果 Docker、网络、下载依赖不可用，**报告阻碍与已做验证，不能伪称运行成功**；不要通过扩大实现范围来绕过。

## 2. 交付目录：最小但可编译

建议按照下面结构交付；必要时可以小范围微调文件名，但不能增加额外架构层。

```text
go-backend-final-lab/
├── cmd/api/main.go                # 读取 HTTP_PORT、启动基础 Gin 服务
├── internal/app/router.go         # /healthz、/、Task 路由注册
├── internal/app/router_test.go    # 只测试 /healthz 及服务基础可用性
├── internal/task/
│   ├── model.go                   # 领域类型占位 / 清晰 TODO，不完成 ORM 模型设计
│   ├── handler.go                 # 五个路由的 501 占位处理
│   ├── service.go                 # Service 空结构 / TODO，不实现业务
│   └── repository.go              # Repository 约定或空骨架，不实现 GORM
├── internal/cache/
│   └── doc.go                     # 只说明将在 B3 添加 go-redis，未实现 Client
├── web/index.html                 # 只有最小结构/表单/输出区域，无 fetch
├── compose.yaml                   # 只有 mysql、redis 两个依赖服务
├── .env.example                   # 本地示例参数，无真实秘密
├── .gitignore
├── README.md                      # 环境和阶段边界，极简运行说明
├── docs/SCAFFOLD_HANDOFF.md       # Codex 的交付报告
├── go.mod / go.sum                 # 正常生成
└── （Dockerfile 暂时不存在；由学习者在 B5 编写）
```

### 2.1 可由你实现的部分

**HTTP 启动：**

- `main.go` 读取 `HTTP_PORT`，默认 `18080`；从仓库根目录运行 `go run ./cmd/api`。
- 用 Gin 注册 `GET /healthz`，返回 `200` 与简短 JSON，例如 `{"status":"ok"}`。
- `/healthz` **仅代表 Go HTTP 进程存活**，不可声称 MySQL / Redis 已就绪。
- 用 Gin 提供同源 `GET /` 静态 `web/index.html`（运行假定从仓库根目录执行；README 写明）。
- Task API 占位路由：`POST /tasks`、`GET /tasks/:id`、`GET /tasks`、`PATCH /tasks/:id`、`DELETE /tasks/:id`；全部按 JSON 返回 `501` 与 `{"error":"not_implemented"}`。**不要**写成功的假数据、内存 CRUD 或 GORM。
- 基础路由注册和 Handler 骨架可完成，核心业务 TODO 放在对应文件。**不要**把一个已经写好的 Service/Store 藏在 TODO 背后。

**目录骨架：**

- `task/model.go` 只允许放注释、占位类型或最简单的字段名提示。**字段约束、GORM tags、事务、存取逻辑由学习者完成。**
- `task/service.go` 与 `task/repository.go` 可包含编译所需的空结构/接口签名，但不能连 DB、写 SQL、写缓存、实现业务规则。
- 允许在 README 中写业务契约说明，不能把契约转换成现成实现。
- 不要提前创建会在初始状态失败的完整业务测试，也不要使用“跳过测试”来伪装业务已完成。

**前端：**

- `index.html` 最多包含标题、创建任务的输入框/按钮、任务列表的静态容器、错误提示容器。
- 不生成任何真实 `fetch('/tasks')` 或 API 联动代码。无需 npm、Vite、CSS 库与构建步骤。
- 可以用简短文案提示“B4 阶段由学习者实现 JavaScript 联调”。

### 2.2 仅交付 MySQL / Redis 的依赖容器

使用与已有项目容易对照的开发版本：**MySQL 8.0、Redis 7-alpine**。

`compose.yaml` **仅有** `mysql`、`redis`；此时不存在 `api` 服务：

- MySQL 的容器端口 `3306`，宿主开发端口默认 `127.0.0.1:13306`。
- Redis 的容器端口 `6379`，宿主开发端口默认 `127.0.0.1:16379`。
- MySQL 创建本地数据库 `final_lab`、用户 `final_lab`；密码通过 `${LAB_MYSQL_PASSWORD:-lab_only_change_me}` 形式的**本地开发默认值**或 `.env` 注入。Root 密码也只允许本地开发默认值。不得使用历史项目密钥。README 明确这些密码不能用于生产。
- MySQL **使用独立 named volume**，数据不会因为普通 `docker compose down` 自动丢失。Redis 仅作练习缓存，无需持久化设计。
- 配置最基础健康检查；不额外创建 Adminer、phpMyAdmin、Prometheus、Nginx 等服务。
- 在 README 说明本阶段 Go API 由 WSL 本机运行，通过 `127.0.0.1:13306` / `127.0.0.1:16379` 访问依赖；**直到 B5 才**让学习者自行增加 API 容器，容器之间使用 `mysql:3306` / `redis:6379`。
- 不生成 Dockerfile，不提前增加 `api` service、构建上下文和最终阶段运行命令。
- 明确 Compose project name 需要显式使用 `-p go-backend-final-lab`，避免与真实项目容器混淆。

**注意**：本阶段 Codex 写基础 Compose 是在搭“实验台”，不代表用户已经通过 Docker 学习验收；使用、检查、修改与排错仍由用户亲自完成。

### 2.3 初始 README 的内容边界

README 只包含：

1. 项目用途与强制保留的未实现内容。
2. Go/WSL/Docker Desktop 基础环境要求（Windows/WSL2 用户）。
3. 从仓库根目录运行 `go run ./cmd/api`、访问 `http://127.0.0.1:18080/healthz` 与 `/`。
4. 启动/停止依赖的示例命令；说明数据库/缓存均为开发环境。
5. 标识 `B0 → B6` 应当逐步完成的 TODO；**不能写成已通过。**
6. 简单安全提醒：`.env` 与凭据不提交；不要误用 `docker compose down -v` 处理真实项目。
7. 指向学习文档 `02_学习执行手册.md`，如果该文件尚未放入仓库，标注由用户提供，**不要自行创作其替代版本**。

---

## 3. 不得越过的“学习者保留区”

| 阶段 | 需要由学习者亲手完成 | Codex 在脚手架阶段必须停止的位置 |
|---|---|---|
| B0 | Linux 命令、`docker` / `compose` 的运行与检查、redis-cli/MySQL CLI | 只生成最小依赖 Compose 和操作说明，**不得代替操作学习** |
| B1 | Task/TaskEvent 模型、Gin 绑定/校验、GORM 连接和 CRUD、错误处理、分页 | 仅 501 占位路由；不生成 GORM CRUD 代码 |
| B2 | 关联双写事务、故障注入、回滚验证、一次 EXPLAIN | 不写 `db.Transaction`，不写 SQL 优化实验 |
| B3 | `redis.NewClient`、Cache Aside、TTL、失效、故障读回源 | `internal/cache` 仅文档占位；没有 Client 实现 |
| B4 | 原生 JS `fetch` GET/POST、Network 观察和错误显示 | 页面只有静态元素，无 JS 网络调用 |
| B5 | Dockerfile、Compose 中新增 API service、容器网络与环境变量、运行诊断 | **不创建 Dockerfile，不新增 API Compose 服务** |
| B6 | `go test` 业务测试、`-race` 微实验、Debug 记录、陌生需求变更 | 只有健康检查基础测试；不生成业务测试答案 |

请严格按此表做范围控制；不要为了“让业务运行”添加任何列为学习者保留的内容。

---

## 4. 脚手架自身的验收——只证明骨架可用

如果环境允许，请执行并在交接文档中逐项记录：

```bash
go fmt ./...
go test ./...
go build ./...
docker compose -p go-backend-final-lab config -q
# 仅在 Docker Desktop/daemon 可用时：
docker compose -p go-backend-final-lab up -d mysql redis
docker compose -p go-backend-final-lab ps
# 启动 Go API 后在另一个终端：
curl -i http://127.0.0.1:18080/healthz
curl -i http://127.0.0.1:18080/tasks
```

**预期状态：**

| 验证 | 预期结果 |
|---|---|
| `go test`、`go build` | 仅脚手架通过，不代表业务已实现 |
| `GET /healthz` | 200，返回基础健康 JSON |
| `GET /tasks` | **501**，返回 `not_implemented` |
| `GET /` | 200，可打开静态 HTML 页面 |
| Compose | 可解析、如环境允许则可启动 MySQL / Redis |
| `git status` | 仅有新 Lab 文件，无历史真实项目修改 |

必须明确 **“五个 Task API 尚未实现”**；这是交付成功条件，不是 Bug。不得将项目描述为完成 CRUD 或具备数据库/Redis 功能。

**如无法下载依赖或 Docker 不可用：**可先完成不依赖它的部分，交接中列出失败的原命令、报错、未验证项及用户下一步动作。不得宣称已通过运行验收。

---

## 5. 交付方式

请在 `docs/SCAFFOLD_HANDOFF.md` 中提供：

- 当前实际文件树、Go module 名、Go/Gin 版本；如有版本选择原因用一句话说明。
- 哪些文件由 Codex 写了**可执行的脚手架**，哪些仅是 TODO。
- 逐项验证命令：**已执行 / 成功 / 失败 / 未执行**；失败命令附核心错误原因。
- `GET /healthz` 和一个 Task 501 返回的真实状态；如果未运行，标注未验证。
- 当前 Compose 服务名、端口、卷、必要环境变量与如何停止。
- 明确下一步从 **B0** 开始；用户应先熟悉 `docker ps / compose ps / logs / exec`，再进入 B1。
- 罗列全部**刻意未实现**事项，不写“剩余功能很简单，我已经顺便补全”。

如果当前目录尚未初始化 Git，可以 `git init`；只有在 Git 身份配置已有且目录确认无其他用户文件时，才可创建一个 `chore: scaffold only` 基线提交。不改全局 Git 配置，不强推、不自动上传 GitHub。

**最终答复格式：**简述已创建的内容、执行了哪些验证、有哪些环境阻碍、从哪个 TODO 开始；不要在答复里输出整套未来业务实现代码。

---

## 6. 给 Codex 的自检提问

在结束前，请自行检查并在交接文档回答：

1. 是否能在**不连数据库**时启动 `/healthz`？如果不能，说明脚手架依赖边界错误，应调整为基础健康检查独立。
2. 五个 Task API 是否都保持 `501`？是否暗中加入内存 CRUD、GORM、go-redis 或事务？如有，删掉越界实现。
3. Compose 是否只有 MySQL + Redis，没有 API？是否仅绑定本机 loopback？
4. 项目是否没有 Dockerfile、没有真实前端 `fetch`、没有业务测试实现？
5. 是否绝对没有改动 RepoLens、Gateway、系统级凭据或全局 Git 配置？
6. 失败或未运行的检查是否如实标明？

六项都符合，才可以交付脚手架。**后续任何补全 Task API 的开发都必须由用户按学习文档逐阶段显式授权，不能由本次任务自动延续。**
