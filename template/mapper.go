package template

import (
	"bytes"
	"fmt"
	"text/template"
)

type mapper struct {
	PackageInfo string
	ClassName   string
}

// 生成 XXXMapper.java

var mapperJava = `
package {{.PackageInfo}};

public interface {{.ClassName}}Mapper extends BaseMapper<{{.ClassName}}> {


	//--PLACE_HOLDER--

}
`

// 生成 XXXMapper.xml

var mapperXml = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="{{.PackageInfo}}.{{.ClassName}}Mapper">

	<!--PLACE_HOLDER-->

</mapper>
`

func GenerateMapperJava() {
	mapper := &mapper{
		PackageInfo: "com.relengxing.demo",
		ClassName:   "User",
	}
	// 创建一个模板对象
	t := template.Must(template.New("mapperJava").Parse(mapperJava))
	// 创建一个缓冲区
	var buf bytes.Buffer
	// 渲染模板到标准输出
	err := t.Execute(&buf, mapper)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 从缓冲区中获取最终的字符串
	output := buf.String()
	// 打印输出结果
	fmt.Println(output)
}

func GenerateMapperXml() {
	mapper := &mapper{
		PackageInfo: "com.relengxing.demo",
		ClassName:   "User",
	}
	// 创建一个模板对象
	t := template.Must(template.New("mapperXml").Parse(mapperXml))
	// 创建一个缓冲区
	var buf bytes.Buffer
	// 渲染模板到标准输出
	err := t.Execute(&buf, mapper)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 从缓冲区中获取最终的字符串
	output := buf.String()
	// 打印输出结果
	fmt.Println(output)
}
