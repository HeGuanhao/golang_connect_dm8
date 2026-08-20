package main

import (
	_ "dm"
	"fmt"
	"log"
	"time"
)

func main() {
	// 初始化数据库连接
	db, err := initDB(DefaultConfig)
	if err != nil {
		log.Fatalf("init db failed: %v", err)
	}
	defer func() {
		if err := closeDB(db); err != nil {
			log.Printf("close db error: %v", err)
		}
	}()

	// ── 1. 查询 Schema 信息 ──────────────────────────────────────
	fmt.Println("\n[1] Schema Info")
	if err = listSchemas(db); err != nil {
		log.Printf("listSchemas error: %v", err)
	}
	if err = listTables(db, "SYSDBA"); err != nil {
		log.Printf("listTables error: %v", err)
	}

	// ── 2. 建表 ──────────────────────────────────────────────────
	fmt.Println("\n[2] Create Table")
	if err = createProductTable(db); err != nil {
		log.Printf("createProductTable error: %v", err)
	}
	if err = listColumns(db, "SYSDBA", "PRODUCT"); err != nil {
		log.Printf("listColumns error: %v", err)
	}

	// ── 3. 单条插入 ───────────────────────────────────────────────
	fmt.Println("\n[3] Insert Product")
	p := Product{
		Name:        "三国演义",
		Author:      "罗贯中",
		Publisher:   "中华书局",
		PublishTime: time.Date(2005, 4, 1, 0, 0, 0, 0, time.Local),
		ProductNo:   "9787101046121",
		OrigPrice:   19.0,
		NowPrice:    15.2,
		Discount:    8.0,
		Description: "《三国演义》是中国第一部长篇章回体小说。",
	}
	insertedID, err := insertProduct(db, p)
	if err != nil {
		log.Printf("insertProduct error: %v", err)
	}

	// ── 4. 查询单条 ───────────────────────────────────────────────
	fmt.Println("\n[4] Query Product")
	if insertedID > 0 {
		if err = queryProduct(db, int(insertedID)); err != nil {
			log.Printf("queryProduct error: %v", err)
		}
	}

	// ── 5. 更新 ───────────────────────────────────────────────────
	fmt.Println("\n[5] Update Product")
	if insertedID > 0 {
		if err = updateProduct(db, int(insertedID), "三国演义（上）"); err != nil {
			log.Printf("updateProduct error: %v", err)
		}
	}

	// ── 6. 事务批量插入 ───────────────────────────────────────────
	fmt.Println("\n[6] Transaction Insert")
	if err = txInsertProducts(db); err != nil {
		log.Printf("txInsertProducts error: %v", err)
	}

	// ── 7. 分页查询 ───────────────────────────────────────────────
	fmt.Println("\n[7] Query All Products (page 1)")
	if err = queryAllProducts(db, 1, 10); err != nil {
		log.Printf("queryAllProducts error: %v", err)
	}

	// ── 8. 删除 ───────────────────────────────────────────────────
	fmt.Println("\n[8] Delete Product")
	if insertedID > 0 {
		if err = deleteProduct(db, int(insertedID)); err != nil {
			log.Printf("deleteProduct error: %v", err)
		}
	}

	// ── 9. 连接池健康检查 ─────────────────────────────────────────
	fmt.Println("\n[9] Pool Health Check")
	if err = checkPool(db); err != nil {
		log.Printf("checkPool error: %v", err)
	}

	fmt.Println("\nAll done.")
}
