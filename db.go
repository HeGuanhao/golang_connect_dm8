package main

import (
	"database/sql"
	"fmt"
	"time"
)

// initDB 初始化数据库连接池
func initDB(cfg DBConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open db failed: %w", err)
	}

	// 连接池核心参数
	db.SetMaxOpenConns(cfg.MaxOpenConns)       // 最大打开连接数
	db.SetMaxIdleConns(cfg.MaxIdleConns)       // 最大空闲连接数
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 连接最大存活时间，防止数据库侧主动断开
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // 空闲连接超时回收，释放不必要的资源

	// 验证连接可用
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db failed: %w", err)
	}

	fmt.Printf("connected to \"%s\" succeed.\n", cfg.DataSourceName)
	printPoolStats(db)
	return db, nil
}

// closeDB 关闭数据库连接池
func closeDB(db *sql.DB) error {
	printPoolStats(db)
	if err := db.Close(); err != nil {
		return fmt.Errorf("close db failed: %w", err)
	}
	fmt.Println("db disconnected.")
	return nil
}

// printPoolStats 打印当前连接池状态
func printPoolStats(db *sql.DB) {
	s := db.Stats()
	fmt.Printf("[pool] open=%d idle=%d inUse=%d waitCount=%d waitDuration=%s maxIdleClosed=%d maxLifetimeClosed=%d\n",
		s.OpenConnections,
		s.Idle,
		s.InUse,
		s.WaitCount,
		s.WaitDuration.Round(time.Millisecond),
		s.MaxIdleClosed,
		s.MaxLifetimeClosed,
	)
}

// checkPool 连接池健康检查，验证连接仍然可用
func checkPool(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("pool health check failed: %w", err)
	}
	printPoolStats(db)
	return nil
}
