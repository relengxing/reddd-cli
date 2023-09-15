package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type QueryParam struct {
	Type     string
	Username string
	Password string
	Ip       string
	Port     string
	Db       string
	Table    string
}

func (q QueryParam) datasourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", q.Username, q.Password, q.Ip, q.Port, q.Db)
}

func (q QueryParam) showColumnSql() string {
	return fmt.Sprintf("SHOW COLUMNS FROM %s", q.Table)

}

type ColumnInfo struct {
	Name        string
	DataType    string
	Length      int64
	NullableStr string
	Key         string
	Default     sql.NullString
	Extra       string
}

func (c ColumnInfo) Nullable() bool {
	if c.NullableStr == "YES" {
		return true
	} else {
		return false
	}
}

func QueryTableColumn(queryParam QueryParam) (*[]ColumnInfo, error) {
	// 连接到 MySQL 数据库
	db, err := sql.Open(queryParam.Type, queryParam.datasourceName())
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		return nil, err
	}
	defer db.Close()

	// 执行查询语句
	rows, err := db.Query(queryParam.showColumnSql())
	if err != nil {
		fmt.Println("执行查询语句失败：", err)
		return nil, err
	}
	defer rows.Close()

	// 解析查询结果
	var tableStructure []ColumnInfo
	for rows.Next() {
		column := ColumnInfo{}
		err := rows.Scan(&column.Name, &column.DataType, &column.NullableStr, &column.Key, &column.Default, &column.Extra)
		if err != nil {
			fmt.Println("解析查询结果失败：", err)
			return nil, err
		}
		tableStructure = append(tableStructure, column)
	}
	return &tableStructure, nil
}

func main() {
	queryParam := QueryParam{
		Type:     "mysql",
		Username: "root",
		Password: "zip94303",
		Ip:       "localhost",
		Port:     "3306",
		Db:       "td",
		Table:    "td_user",
	}
	// 连接到 MySQL 数据库
	db, err := sql.Open(queryParam.Type, queryParam.datasourceName())
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		return
	}
	defer db.Close()

	// 执行查询语句
	rows, err := db.Query(queryParam.showColumnSql())
	if err != nil {
		fmt.Println("执行查询语句失败：", err)
		return
	}
	defer rows.Close()

	// 解析查询结果
	var tableStructure []ColumnInfo
	for rows.Next() {
		column := ColumnInfo{}
		err := rows.Scan(&column.Name, &column.DataType, &column.NullableStr, &column.Key, &column.Default, &column.Extra)
		if err != nil {
			fmt.Println("解析查询结果失败：", err)
			return
		}
		tableStructure = append(tableStructure, column)
	}

	// 输出表结构
	for _, column := range tableStructure {
		fmt.Printf("Name: %s\n", column.Name)
		fmt.Printf("DataType: %s\n", column.DataType)
		fmt.Printf("Length: %d\n", column.Length)
		fmt.Printf("Nullable: %v\n", column.NullableStr)
		fmt.Printf("Key: %s\n", column.Key)
		if column.Default.Valid {
			fmt.Printf("Default: %s\n", column.Default.String)
		} else {
			fmt.Println("Default: NULL")
		}
		fmt.Printf("Extra: %s\n", column.Extra)
		fmt.Println("----------------------")
	}
}
