package template

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 生成文件夹

func GenerateDomainFolder(domain string) error {
	// 判断文件夹是否存在，存在则抛出异常

	// 创建对应目录
	/*
		- xxxx
		-- model (实体类)
		-- dao (存储接口)
		-- api	(对外调用)
		-- service	(业务逻辑)
	*/
	target := ""
	target = target
	wd, _ := os.Getwd()
	err := filepath.Walk(wd, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && strings.HasSuffix(path, "-domain") {
			target = filepath.Join(target, path)
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	target = filepath.Join(target, "src", "main", "java")
	target = filepath.Join(target, "package", domain)

	// 创建文件夹
	if _, err := os.Stat(target); os.IsNotExist(err) {
		if err := os.MkdirAll(target, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	// 创建子文件夹
	model := filepath.Join(target, "model")
	if _, err := os.Stat(model); os.IsNotExist(err) {
		if err := os.MkdirAll(model, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
