package main

import "time"

// DBConfig 数据库连接池配置
type DBConfig struct {
	DriverName      string
	DataSourceName  string
	MaxOpenConns    int           // 最大打开连接数（0 = 不限制）
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大存活时间
	ConnMaxIdleTime time.Duration // 连接最大空闲时间
}

// DefaultConfig 默认配置，修改为你的实际 DM8 地址
var DefaultConfig = DBConfig{
	DriverName:      "dm",
	DataSourceName:  "dm://SYSDBA:Dameng@7777@10.47.80.80:31775",
	MaxOpenConns:    10,
	MaxIdleConns:    5,
	ConnMaxLifetime: 30 * time.Minute, // 连接最多复用 30 分钟
	ConnMaxIdleTime: 10 * time.Minute, // 空闲超过 10 分钟自动回收
}
