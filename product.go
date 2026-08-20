package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"dm"
)

// Product 产品信息
type Product struct {
	ProductID   int
	Name        string
	Author      string
	Publisher   string
	PublishTime time.Time
	ProductNo   string
	OrigPrice   float64
	NowPrice    float64
	Discount    float64
	Description string
}

// createProductTable 创建产品表（如果不存在则创建）
func createProductTable(db *sql.DB) error {
	// DM8 不支持 IF NOT EXISTS，先尝试建表，已存在则忽略错误
	ddl := `CREATE TABLE SYSDBA.product (
		productid            INT IDENTITY(1,1) PRIMARY KEY,
		name                 VARCHAR(50)  NOT NULL,
		author               VARCHAR(100),
		publisher            VARCHAR(100),
		publishtime          DATE,
		product_subcategoryid INT,
		productno            VARCHAR(50),
		satetystocklevel     INT,
		originalprice        DECIMAL(10,4),
		nowprice             DECIMAL(10,4),
		discount             DECIMAL(4,2),
		description          CLOB,
		photo                BLOB,
		type                 CHAR(2),
		papertotal           INT,
		wordtotal            INT,
		sellstarttime        DATE,
		sellendtime          DATE
	)`
	if _, err := db.Exec(ddl); err != nil {
		// -2124 = 对象已存在，忽略
		if !strings.Contains(err.Error(), "-2124") {
			return fmt.Errorf("create table failed: %w", err)
		}
		fmt.Println("table SYSDBA.product already exists, skip.")
		return nil
	}
	fmt.Println("table SYSDBA.product created.")
	return nil
}

// insertProduct 插入一条产品记录
func insertProduct(db *sql.DB, p Product) (int64, error) {
	sqlStr := `INSERT INTO SYSDBA.product(name,author,publisher,publishtime,
		product_subcategoryid,productno,satetystocklevel,originalprice,nowprice,discount,
		description,type,papertotal,wordtotal,sellstarttime,sellendtime)
		VALUES(:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16)`

	publishTime := p.PublishTime
	if publishTime.IsZero() {
		publishTime = time.Now()
	}
	sellStart := time.Now()
	sellEnd := time.Date(2099, 12, 31, 0, 0, 0, 0, time.Local)

	result, err := db.Exec(sqlStr,
		p.Name, p.Author, p.Publisher, publishTime,
		4, p.ProductNo, 10, p.OrigPrice, p.NowPrice, p.Discount,
		p.Description, "25", 943, 93000, sellStart, sellEnd,
	)
	if err != nil {
		return 0, fmt.Errorf("insert product failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %w", err)
	}
	fmt.Printf("insertProduct succeed, id=%d\n", id)
	return id, nil
}

// updateProduct 更新产品名称
func updateProduct(db *sql.DB, productID int, newName string) error {
	sqlStr := "UPDATE SYSDBA.product SET name = :1 WHERE productid = :2"
	result, err := db.Exec(sqlStr, newName, productID)
	if err != nil {
		return fmt.Errorf("update product failed: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected failed: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no product found with id=%d", productID)
	}
	fmt.Printf("updateProduct succeed, rows affected=%d\n", rows)
	return nil
}

// queryProduct 查询单条产品（含 CLOB 内容）
func queryProduct(db *sql.DB, productID int) error {
	sqlStr := "SELECT productid,name,author,description,photo FROM SYSDBA.product WHERE productid=:1"
	rows, err := db.Query(sqlStr, productID)
	if err != nil {
		return fmt.Errorf("query product failed: %w", err)
	}
	defer rows.Close()

	fmt.Println("=== queryProduct results ===")
	found := false
	for rows.Next() {
		found = true
		var id int
		var name, author string
		var description dm.DmClob
		var photo dm.DmBlob
		if err = rows.Scan(&id, &name, &author, &description, &photo); err != nil {
			return fmt.Errorf("scan product failed: %w", err)
		}
		desc, err := readClob(&description)
		if err != nil {
			return fmt.Errorf("read description clob failed: %w", err)
		}
		blobLen, _ := photo.GetLength()
		fmt.Printf("  id=%d name=%s author=%s\n  description=%s\n  photoLen=%d\n",
			id, name, author, desc, blobLen)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	if !found {
		fmt.Printf("  no product found with id=%d\n", productID)
	}
	return nil
}

// readClob 读取 DmClob 全部内容（位置从 1 开始）
func readClob(clob *dm.DmClob) (string, error) {
	length, err := clob.GetLength()
	if err != nil {
		return "", err
	}
	if length == 0 {
		return "", nil
	}
	s, err := clob.ReadString(1, int(length))
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return s, nil
}

// queryAllProducts 查询所有产品（分页）
// 注意：ROWNUM 在 ORDER BY 之前分配，必须先在内层排序、外层再用 ROWNUM 过滤，
// 否则分页结果顺序不确定。
func queryAllProducts(db *sql.DB, page, pageSize int) error {
	if page < 1 || pageSize < 1 {
		return fmt.Errorf("invalid pagination: page=%d pageSize=%d", page, pageSize)
	}
	offset := (page-1)*pageSize + 1
	end := page * pageSize
	sqlStr := `SELECT productid, name, author, originalprice, nowprice FROM (
		           SELECT ROWNUM rn, t.* FROM (
		               SELECT productid, name, author, originalprice, nowprice
		               FROM SYSDBA.product ORDER BY productid
		           ) t
		       ) WHERE rn >= :1 AND rn <= :2`
	rows, err := db.Query(sqlStr, offset, end)
	if err != nil {
		return fmt.Errorf("query all products failed: %w", err)
	}
	defer rows.Close()

	fmt.Printf("=== Products (page=%d, size=%d) ===\n", page, pageSize)
	count := 0
	for rows.Next() {
		var id int
		var name, author string
		var origPrice, nowPrice float64
		if err = rows.Scan(&id, &name, &author, &origPrice, &nowPrice); err != nil {
			return fmt.Errorf("scan product row failed: %w", err)
		}
		count++
		fmt.Printf("  id=%-5d name=%-20s author=%-15s price=%.2f->%.2f\n",
			id, name, author, origPrice, nowPrice)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	if count == 0 {
		fmt.Println("  (empty page)")
	}
	return nil
}

// deleteProduct 删除产品
func deleteProduct(db *sql.DB, productID int) error {
	sqlStr := "DELETE FROM SYSDBA.product WHERE productid = :1"
	result, err := db.Exec(sqlStr, productID)
	if err != nil {
		return fmt.Errorf("delete product failed: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected failed: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no product found with id=%d", productID)
	}
	fmt.Printf("deleteProduct succeed, rows affected=%d\n", rows)
	return nil
}
