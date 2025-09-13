# GoRepo

## Go 知识驿站

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


## Go 工具集合

- https://github.com/incu6us/goimports-reviser  
Tool for Golang to sort goimports

- https://github.com/kaptinlin/jsonrepair  
Easily repair invalid JSON documents with the Golang JSONRepair Library.

- https://github.com/tbphp/gpt-load  
A high-performance, enterprise-grade AI API transparent proxy service designed specifically for enterprises and developers who need to integrate multiple AI services. Built with Go, featuring intelligent key management, load balancing, and comprehensive monitoring capabilities, designed for high-concurrency production environments.

