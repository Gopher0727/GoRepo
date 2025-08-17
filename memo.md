1. Go 的 JSON/YAML/TOML 等库在反序列化时，都需要结构体字段是导出的（大写开头）

2. 格式化 Go 文件的导入顺序
```shell
goimports-reviser -rm-unused -format -recursive .
```

3. 为什么 MySQL 要用 database/sql 的统一接口，而 Redis 测试直接用驱动就可以？

- MySQL/SQL 数据库的情况：
> SQL 数据库很多种：MySQL、PostgreSQL、SQLite、SQL Server……，每种数据库的底层 API 都不一样，但 SQL 语法有共通点，Go 提供 database/sql 这个统一接口的目的就是 `一次写代码，多种数据库通用`，如果你直接用 MySQL 驱动：你的代码就只能跑在 MySQL 上，换成 PostgreSQL、SQLite 等数据库就得重写，连接池、事务管理、Prepare、QueryRow 等功能都得自己实现，所以 _ "github.com/go-sql-driver/mysql" + database/sql 是为了 可移植、方便、统一管理资源。

- Redis 的情况：
> Redis 本身就是一个 key-value 存储，没有 SQL。Go 没有一个“统一接口”管理所有 NoSQL，Redis 驱动直接提供了高层 API，用 github.com/redis/go-redis/v9 就能直接操作 key-value，没有必要像 database/sql 那样抽象。

4. GORM

- 最流行的 Go ORM。  
- 全面支持 CRUD、关系、事务、钩子、预加载等。  
- 提供自动迁移（AutoMigrate）功能。  
- 链式 API，非常直观。

> GORM 是 关系型数据库 ORM，专门处理 MySQL、PostgreSQL、SQLite、SQL Server 等 SQL 数据库。它的功能都是围绕表、行、列、SQL 查询展开的，比如 CRUD、事务、关联、预加载等。
>
> Redis 是 内存键值数据库，通常用于缓存、队列、计数器等。其数据结构多样（String、Hash、List、Set、SortedSet 等），并不是关系型表格结构；操作是 key-value 风格，不走 SQL，所以 ORM 的映射意义不大。


# Go 常见框架 & 组件库

## Web 框架/路由

- Gin
    轻量、高性能、API 风格清晰，适合构建高 QPS 的 REST 服务与中小型后台。常用作“路由 + 中间件”主干
- Echo
    极简且功能完备（middleware、group、websocket 等），API 友好，适合快速开发 REST 服务与中间件链路
- Fiber
    从 Express 借鉴设计，基于 fasthttp，强调极致性能与低延迟，适合对吞吐要求非常高的场景
- Beego / Revel / Buffalo
    全栈型或脚手架型框架（有生成器、前端 pipeline、MVC 支持），适合想要“开箱即用”全套工具的团队

## 微服务/分布式架构

- Go kit 
    面向微服务的工具包，强调可观测性、可插拔、契约化（transport、endpoint、service 分层），适合企业级、可维护的微服务架构
- Kratos / go-micro / kitex
    提供 RPC、服务发现、治理等能力，适合需要完整微服务平台能力的项目

## 数据库/ORM/SQL

- GORM
    功能齐全、社区大、上手快的 ORM（关联、事务、钩子等），适合多数 CRUD 密集型应用
- ent
    “Schema as code”、代码生成、类型安全，适合大型复杂数据模型和需要静态类型保证的系统

## 缓存、KV、消息与队列

- go-redis / redigo
    Redis 客户端（go-redis 目前主流且活跃）
- NATS / Sarama (Kafka) / NSQ
    不同消息系统的客户端生态，按可靠性/吞吐/运维成本选型

## 日志/可观测性/追踪

- zap（Uber）
    快速、结构化日志（JSON），生产环境常用
- logrus / zerolog

## 依赖注入/应用生命周期

- uber-go/fx / dig
    提供 DI、启动/停止钩子与组件注入，适合大型服务或需要明确依赖管理的项目

## 验证、迁移、测试、调度等工具 

- validator.v10（go-playground）
    字段校验常用库
- migrate / goose
    数据库 schema 迁移。
- testify / ginkgo + gomega / httpexpect
    单元与集成测试框架
- robfig/cron
    定时任务
- gorilla/websocket / gorilla/mux
    websocket 与经典 router（历史悠久、生态丰富）

## 常用实用组件（配置、CLI、存储）

- viper / envconfig
    配置管理（文件 + 环境变量）。
- cobra / urfave/cli
    CLI 工具开发。
- aws-sdk-go-v2 / minio-go
    对象存储与云服务。