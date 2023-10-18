package main

import "reddd-cli/template"

func main() {
	//queryParam := column.QueryParam{
	//	Type:     "mysql",
	//	Username: "root",
	//	Password: "zip94303",
	//	Ip:       "localhost",
	//	Port:     "3306",
	//	Db:       "td",
	//	Table:    "td_user",
	//}
	//tableColumn, _ := column.QueryTableColumn(queryParam)
	//marshal, _ := json.Marshal(tableColumn)
	//fmt.Println(string(marshal))
	template.GenerateEntity()
	//template.GenerateController()
}
