# DM8 Go Demo

连接达梦数据库(DM8)并测试基本数据操作的 Go demo。

## 项目结构

```
src/
├── app/
│   ├── main.go         # 入口，串联所有演示步骤
│   ├── config.go       # 数据库连接配置
│   ├── db.go           # 连接池初始化/关闭
│   ├── schema.go       # Schema/Table/Column 元数据查询
│   ├── product.go      # 产品表 CRUD 操作
│   └── transaction.go  # 事务演示
├── dm/                 # DM8 Go 驱动源码
└── golang.org/         # 依赖
```

## 配置

修改 `src/app/config.go` 中的连接信息：

```go
var DefaultConfig = DBConfig{
    DriverName:     "dm",
    DataSourceName: "dm://SYSDBA:SYSDBA@<host>:<port>",
    MaxOpenConns:   10,
    MaxIdleConns:   5,
}
```

## 运行

```bash
# 设置 GOPATH
export GOPATH=$(pwd)

# 编译并运行
go run src/app/*.go
```

## 演示内容

1. 查询所有 Schema 和表信息
2. 创建 production.product 表
3. 单条插入产品
4. 查询单条产品（含 CLOB/BLOB 字段）
5. 更新产品
6. 事务批量插入（四大名著）
7. 分页查询所有产品
8. 删除产品
