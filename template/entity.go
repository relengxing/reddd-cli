package template

import (
	"bytes"
	"fmt"
	"reddd-cli/column"
	"reddd-cli/utils"
	"text/template"
)

var tmpl = `
package {{.PackageInfo}};
import java.util.*;

@Data
@EqualsAndHashCode(callSuper = false)
@Accessors(chain = true)
@TableName("{{.TableName}}")
@ApiModel(value = "{{.Comment}}", description = "{{.Comment}}")
public class {{.ClassName}}Entity implements Serializable {

   private static final long serialVersionUID = 1L;

   private Long id;

	{{range .Field}}
   @ApiModelProperty(value = "{{.FieldComment}}")
   private {{.FieldType}} {{.FieldName}};

	{{end}}
}
`

// 生成 XXXEntity.java
type entity struct {
	PackageInfo string
	TableName   string
	ClassName   string
	Comment     string
	Field       []field
}

type field struct {
	FieldComment string
	FieldType    string
	FieldName    string
}

func GenerateEntity() {
	queryParam := &column.QueryParam{
		Type:     "mysql",
		Username: "root",
		Password: "zip94303",
		Ip:       "localhost",
		Port:     "3306",
		Db:       "td",
		Table:    "td_user",
	}
	tableInfo, err := column.QueryTableColumn(*queryParam)
	if err != nil {
		fmt.Println("数据库连接失败：", err)
	}
	// 创建Employee结构体数组
	fieldArr := make([]field, len(*tableInfo.ColumnInfo))

	// 遍历结构体数组，进行转换
	for i, column := range *tableInfo.ColumnInfo {
		fieldArr[i] = field{
			FieldComment: column.Comment,
			FieldType:    column.DataTypeJava(),
			FieldName:    column.Name,
		}
	}
	entity := entity{
		PackageInfo: "com.relengxing",
		TableName:   *tableInfo.TableName,
		Comment:     *tableInfo.Comment,
		ClassName:   utils.ToCamelCase(*tableInfo.TableName),
		Field:       fieldArr,
	}

	// 创建一个模板对象
	t := template.Must(template.New("entity").Parse(tmpl))

	// 创建一个缓冲区
	var buf bytes.Buffer
	// 渲染模板到标准输出
	err = t.Execute(&buf, entity)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// 从缓冲区中获取最终的字符串
	output := buf.String()

	// 打印输出结果
	fmt.Println(output)
}

// 生成 XXXDTO.java

// 生成 XXXReq.java

// 生成 XXXResp.java
