package main

import (
	"database/sql"
	"fmt"
	"time"
)

// txInsertProducts 使用事务批量插入产品，演示事务回滚
func txInsertProducts(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}

	// 提交后 Rollback 是安全的 no-op；未提交时自动回滚
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
			fmt.Println("transaction rolled back.")
		}
	}()

	sqlStr := `INSERT INTO SYSDBA.product(name,author,publisher,publishtime,
		product_subcategoryid,productno,satetystocklevel,originalprice,nowprice,discount,
		description,type,papertotal,wordtotal,sellstarttime,sellendtime)
		VALUES(:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16)`

	publishTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	sellStart := time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)
	sellEnd := time.Date(2099, 12, 31, 0, 0, 0, 0, time.Local)

	books := []struct{ name, author, no string }{
		{"水浒传", "施耐庵", "9787020015016"},
		{"西游记", "吴承恩", "9787020008742"},
		{"红楼梦", "曹雪芹", "9787020002207"},
	}

	// 先删除同名旧数据，保证重复执行结果一致
	for _, b := range books {
		if _, err = tx.Exec("DELETE FROM SYSDBA.product WHERE name = :1", b.name); err != nil {
			return fmt.Errorf("tx delete [%s] failed: %w", b.name, err)
		}
	}

	for _, b := range books {
		if _, err = tx.Exec(sqlStr,
			b.name, b.author, "人民文学出版社", publishTime,
			4, b.no, 10, 39.0, 29.0, 7.5,
			b.name+"是中国古典四大名著之一。", "25", 800, 80000, sellStart, sellEnd,
		); err != nil {
			return fmt.Errorf("tx insert [%s] failed: %w", b.name, err)
		}
		fmt.Printf("  tx inserted: %s\n", b.name)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit failed: %w", err)
	}
	committed = true
	fmt.Println("transaction committed.")
	return nil
}

// txTransferDemo 演示事务中的转账逻辑（模拟）
func txTransferDemo(db *sql.DB, fromID, toID int, amount float64) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
			fmt.Println("transfer rolled back.")
		}
	}()

	// 扣减 from 的价格（模拟转账）
	if _, err = tx.Exec(
		"UPDATE SYSDBA.product SET nowprice = nowprice - :1 WHERE productid = :2",
		amount, fromID,
	); err != nil {
		return fmt.Errorf("deduct failed: %w", err)
	}

	// 增加 to 的价格（模拟转账）
	if _, err = tx.Exec(
		"UPDATE SYSDBA.product SET nowprice = nowprice + :1 WHERE productid = :2",
		amount, toID,
	); err != nil {
		return fmt.Errorf("credit failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit failed: %w", err)
	}
	committed = true
	fmt.Printf("transfer %.2f from product#%d to product#%d succeed.\n", amount, fromID, toID)
	return nil
}
