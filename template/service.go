package template

import (
	"bytes"
	"fmt"
	"text/template"
)

// 生成 XXXService.java

type service struct {
	PackageInfo string
	ClassName   string
}

var serviceTpl = `
package {{.PackageInfo}};

public interface I{{.ClassName}}Service extends IService<{{.ClassName}}> {

}
`

func GenerateService() {
	service := &service{
		PackageInfo: "com.relengxing.demo",
		ClassName:   "User",
	}
	// 创建一个模板对象
	t := template.Must(template.New("serviceTpl").Parse(serviceTpl))
	// 创建一个缓冲区
	var buf bytes.Buffer
	// 渲染模板到标准输出
	err := t.Execute(&buf, service)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 从缓冲区中获取最终的字符串
	output := buf.String()
	// 打印输出结果
	fmt.Println(output)
}

var serviceImplTpl = `
package {{.PackageInfo}}.impl;

@Service
@Slf4j
public class {{.ClassName}}ServiceImpl extends ServiceImpl<{{.ClassName}}Mapper, {{.ClassName}}> implements I{{.ClassName}}Service {

}
`

func GenerateServiceImpl() {
	service := &service{
		PackageInfo: "com.relengxing.demo",
		ClassName:   "User",
	}
	// 创建一个模板对象
	t := template.Must(template.New("serviceImplTpl").Parse(serviceImplTpl))
	// 创建一个缓冲区
	var buf bytes.Buffer
	// 渲染模板到标准输出
	err := t.Execute(&buf, service)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 从缓冲区中获取最终的字符串
	output := buf.String()
	// 打印输出结果
	fmt.Println(output)
}
