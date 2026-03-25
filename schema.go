package main

import (
	"database/sql"
	"fmt"
)

// listSchemas 查询所有 schema
func listSchemas(db *sql.DB) error {
	rows, err := db.Query("SELECT NAME FROM SYSOBJECTS WHERE TYPE$ = 'SCH' ORDER BY NAME")
	if err != nil {
		return fmt.Errorf("query schemas failed: %w", err)
	}
	defer rows.Close()

	fmt.Println("=== Schemas ===")
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan schema failed: %w", err)
		}
		fmt.Println(" -", name)
	}
	return rows.Err()
}

// listTables 查询指定 schema 下的所有表
func listTables(db *sql.DB, schema string) error {
	query := `SELECT o.NAME FROM SYSOBJECTS o
	          JOIN SYSOBJECTS s ON o.SCHID = s.ID
	          WHERE s.NAME = ? AND o.TYPE$ = 'SCHOBJ' AND o.SUBTYPE$ = 'UTAB'
	          ORDER BY o.NAME`
	rows, err := db.Query(query, schema)
	if err != nil {
		return fmt.Errorf("query tables failed: %w", err)
	}
	defer rows.Close()

	fmt.Printf("=== Tables in schema [%s] ===\n", schema)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan table failed: %w", err)
		}
		fmt.Println(" -", name)
	}
	return rows.Err()
}

// listColumns 查询指定表的列信息
func listColumns(db *sql.DB, schema, table string) error {
	query := `SELECT c.NAME, t.NAME, c.NULLABLE$
	          FROM SYSCOLUMNS c
	          JOIN SYSOBJECTS o ON c.ID = o.ID
	          JOIN SYSOBJECTS s ON o.SCHID = s.ID
	          JOIN SYSTYPES t ON c.TYPE$ = t.NAME
	          WHERE s.NAME = ? AND o.NAME = ?
	          ORDER BY c.COLID`
	rows, err := db.Query(query, schema, table)
	if err != nil {
		return fmt.Errorf("query columns failed: %w", err)
	}
	defer rows.Close()

	fmt.Printf("=== Columns of [%s.%s] ===\n", schema, table)
	for rows.Next() {
		var colName, dataType string
		var nullable int
		if err := rows.Scan(&colName, &dataType, &nullable); err != nil {
			return fmt.Errorf("scan column failed: %w", err)
		}
		nullableStr := "Y"
		if nullable == 0 {
			nullableStr = "N"
		}
		fmt.Printf("  %-30s %-20s nullable=%s\n", colName, dataType, nullableStr)
	}
	return rows.Err()
}
