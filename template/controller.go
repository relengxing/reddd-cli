package template

import (
	"bytes"
	"fmt"
	"text/template"
)

// 生成  XXXController.java

type controller struct {
	PackageInfo string
	ClassName   string
}

var controllerTpl = `
package {{.PackageInfo}}.controller;

@Api(tags = "")
@RestController
@RequestMapping("/{{.ClassName}}")
public class {{.ClassName}}Controller {

}
`

func GenerateController() {
	ctrl := &controller{
		PackageInfo: "com.relengxing.demo",
		ClassName:   "User",
	}
	// 创建一个模板对象
	t := template.Must(template.New("controllerTpl").Parse(controllerTpl))
	// 创建一个缓冲区
	var buf bytes.Buffer
	// 渲染模板到标准输出
	err := t.Execute(&buf, ctrl)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 从缓冲区中获取最终的字符串
	output := buf.String()
	// 打印输出结果
	fmt.Println(output)
}
