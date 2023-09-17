package column

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reddd-cli/utils"
	"strings"
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

func (q QueryParam) queryTableInfoSql() string {
	return fmt.Sprintf(`SELECT TABLE_NAME,TABLE_COMMENT FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = '%s' AND table_name = '%s'`, q.Db, q.Table)
}

func (q QueryParam) queryColumnInfoSql() string {
	return fmt.Sprintf(`SELECT COLUMN_NAME name, COLUMN_COMMENT comment FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = '%s' AND table_name = '%s'`, q.Db, q.Table)
}

type TableInfo struct {
	TableName  *string
	Comment    *string
	ColumnInfo *[]ColumnInfo
}
type ColumnInfo struct {
	Name     string
	DataType string
	Comment  string
}

func (c *ColumnInfo) DataTypeJava() string {
	if strings.Contains(c.DataType, "char") {
		return "String"
	} else if strings.Contains(c.DataType, "bigint") {
		return "Long"
	}
	switch c.DataType {
	case "int":
		return "Integer"
	case "bigint":
		return "Long"
	case "varchar":
		return "String"
	case "datetime":
		return "Date"
	case "date":
		return "Date"
	default:
		return "Object"
	}
}

type MysqlColumn struct {
	Name        string
	DataType    string
	Length      int64
	NullableStr string
	Key         string
	Default     sql.NullString
	Extra       string
	Comment     string
}

func (c *MysqlColumn) Nullable() bool {
	if c.NullableStr == "YES" {
		return true
	} else {
		return false
	}
}

func (c *MysqlColumn) SetComment(comment string) {
	c.Comment = comment
}

func (c *MysqlColumn) camelName() string {
	camelName := utils.ToCamelCase(c.Name)
	return strings.ToLower(string(camelName[0])) + camelName[1:]
}

func QueryTableColumn(queryParam QueryParam) (*TableInfo, error) {
	// 连接到 MySQL 数据库
	db, err := sql.Open(queryParam.Type, queryParam.datasourceName())
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		return nil, err
	}
	defer db.Close()

	// 表信息
	tableInfo := &TableInfo{
		TableName:  nil,
		Comment:    nil,
		ColumnInfo: nil,
	}
	rows, err := db.Query(queryParam.queryTableInfoSql())
	if err != nil {
		fmt.Println("执行查询语句失败：", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tableInfo.TableName, &tableInfo.Comment)
		if err != nil {
			fmt.Println("查询表信息失败：", err)
			return nil, err
		}
	}
	// 列信息
	rows, err = db.Query(queryParam.showColumnSql())
	if err != nil {
		fmt.Println("执行查询语句失败：", err)
		return nil, err
	}
	defer rows.Close()

	var mysqlColumn []MysqlColumn
	for rows.Next() {
		column := MysqlColumn{}
		err := rows.Scan(&column.Name, &column.DataType, &column.NullableStr, &column.Key, &column.Default, &column.Extra)
		if err != nil {
			fmt.Println("解析查询结果失败：", err)
			return nil, err
		}
		mysqlColumn = append(mysqlColumn, column)
	}
	// 查询列备注
	rows, err = db.Query(queryParam.queryColumnInfoSql())
	if err != nil {
		fmt.Println("执行注释查询语句失败：", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		column := MysqlColumn{}
		err := rows.Scan(&column.Name, &column.Comment)
		if err != nil {
			fmt.Println("解析查询结果失败：", err)
			return nil, err
		}
		for i, _ := range mysqlColumn {
			if column.Name == mysqlColumn[i].Name {
				mysqlColumn[i].SetComment(column.Comment)
				continue
			}
		}
	}
	var columnInfoList []ColumnInfo
	for _, column := range mysqlColumn {
		columnInfo := ColumnInfo{
			Name:     column.Name,
			DataType: column.DataType,
			Comment:  column.Comment,
		}
		columnInfoList = append(columnInfoList, columnInfo)

	}
	tableInfo.ColumnInfo = &columnInfoList

	return tableInfo, nil
}
