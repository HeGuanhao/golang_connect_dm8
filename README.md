# DM8 Go Demo

连接达梦数据库(DM8)并测试基本数据操作的 Go demo。

## 项目结构

```
├── main.go         # 入口，串联所有演示步骤
├── config.go       # 数据库连接配置（支持环境变量覆盖）
├── db.go           # 连接池初始化/关闭/健康检查
├── schema.go       # Schema/Table/Column 元数据查询
├── product.go      # 产品表 CRUD 操作（含 CLOB/BLOB 读取、分页）
├── transaction.go  # 事务演示（批量插入、模拟转账）
├── go.mod          # 通过 replace 引用本地驱动
└── src/
    ├── dm/         # DM8 Go 驱动源码
    └── golang.org/ # 驱动依赖
```

## 配置

默认连接信息在 `config.go` 的 `DefaultConfig` 中，可通过环境变量覆盖（避免密码硬编码）：

```bash
export DM_DSN="dm://SYSDBA:SYSDBA@<host>:<port>"
```

连接池参数：

```go
var DefaultConfig = DBConfig{
    DriverName:      "dm",
    DataSourceName:  envOr("DM_DSN", "dm://SYSDBA:SYSDBA@<host>:<port>"),
    MaxOpenConns:    10,
    MaxIdleConns:    5,
    ConnMaxLifetime: 30 * time.Minute,
    ConnMaxIdleTime: 10 * time.Minute,
}
```

## 运行

Go Module 模式，直接运行：

```bash
go run .
```

## 演示内容

1. 查询所有 Schema、表和列信息
2. 创建 `SYSDBA.product` 表（已存在则跳过）
3. 单条插入产品
4. 查询单条产品（读取 CLOB 内容、BLOB 长度）
5. 更新产品
6. 事务批量插入（四大名著，失败自动回滚）
7. 分页查询所有产品（内层排序 + 外层 ROWNUM 过滤）
8. 删除产品
9. 连接池健康检查
