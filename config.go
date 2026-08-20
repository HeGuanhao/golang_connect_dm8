package main

import (
	"os"
	"strings"
	"time"
)

// DBConfig 数据库连接池配置
type DBConfig struct {
	DriverName      string
	DataSourceName  string
	MaxOpenConns    int           // 最大打开连接数（0 = 不限制）
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大存活时间
	ConnMaxIdleTime time.Duration // 连接最大空闲时间
}

// MaskedDSN 返回脱敏后的 DSN（隐藏密码），用于日志输出
func (c DBConfig) MaskedDSN() string {
	dsn := c.DataSourceName
	// dm://user:password@host:port -> dm://user:****@host:port
	if i := strings.Index(dsn, "://"); i >= 0 {
		if j := strings.LastIndex(dsn, "@"); j > i {
			userInfo := dsn[i+3 : j]
			if k := strings.Index(userInfo, ":"); k >= 0 {
				return dsn[:i+3] + userInfo[:k] + ":****" + dsn[j:]
			}
		}
	}
	return dsn
}

// DefaultConfig 默认配置，可用环境变量 DM_DSN 覆盖连接串（避免密码硬编码进代码）
var DefaultConfig = DBConfig{
	DriverName:      "dm",
	DataSourceName:  envOr("DM_DSN", "dm://SYSDBA:Dameng@7777@10.47.80.80:31775"),
	MaxOpenConns:    10,
	MaxIdleConns:    5,
	ConnMaxLifetime: 30 * time.Minute, // 连接最多复用 30 分钟
	ConnMaxIdleTime: 10 * time.Minute, // 空闲超过 10 分钟自动回收
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
