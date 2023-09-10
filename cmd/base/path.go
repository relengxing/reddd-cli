package base

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reddd-cli/config"
	"reddd-cli/utils"
	"strings"

	"github.com/fatih/color"
)

var ignoreArr = []string{".git", ".idea", "target"}
var whiteArr = []string{".jpg", ".png", ".jpeg"}

func redddHome() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	home := filepath.Join(dir, ".reddd")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func homeWithDir(dir string) string {
	home := filepath.Join(redddHome(), dir)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func copyFile(src, dst string, replaces []string) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	buf, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	var old string
	for i, next := range replaces {
		if i%2 == 0 {
			old = next
			continue
		}
		buf = bytes.ReplaceAll(buf, []byte(old), []byte(next))
	}
	return os.WriteFile(dst, buf, srcinfo.Mode())
}

func copyDir(src, dst string, replaces, ignores []string) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcinfo.Mode())
	if err != nil {
		return err
	}

	fds, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, fd := range fds {
		if hasSets(fd.Name(), ignores) {
			continue
		}
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())
		var e error
		if fd.IsDir() {
			e = copyDir(srcfp, dstfp, replaces, ignores)
		} else {
			e = copyFile(srcfp, dstfp, replaces)
		}
		if e != nil {
			return e
		}
	}
	return nil
}

func hasSets(name string, sets []string) bool {
	for _, ig := range sets {
		if ig == name {
			return true
		}
	}
	return false
}

func Tree(path string, dir string) {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			fmt.Printf("%s %s (%v bytes)\n", color.GreenString("CREATED"), strings.Replace(path, dir+"/", "", -1), info.Size())
		}
		return nil
	})
}

func walk(home string, to string) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && utils.ContainsString(ignoreArr, info.Name()) {
			return filepath.SkipDir
		}

		if !info.IsDir() { // 如果当前路径不是一个目录，则将其添加到文件路径数组中
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if !utils.ContainsString(whiteArr, filepath.Ext(path)) {
				// rewrite file
				content = reWriteFileContent(content)
			}
			//return file.getPath().replace(projectBaseDir, projectBaseDirNew) // 新目录
			//.replace(PACKAGE_NAME.replaceAll("\\.", Matcher.quoteReplacement(separator)),
			//	packageNameNew.replaceAll("\\.", Matcher.quoteReplacement(separator)))
			//.replace(ARTIFACT_ID, artifactIdNew) //
			//.replaceAll(StrUtil.upperFirst(ARTIFACT_ID), StrUtil.upperFirst(artifactIdNew));
			newPath := strings.ReplaceAll(path, home, to)
			newPath = strings.ReplaceAll(newPath, strings.ReplaceAll(config.PackageNameOld, "\\.", string(filepath.Separator)), strings.ReplaceAll(config.PackageNameOld, "\\.", string(filepath.Separator)))
			newPath = strings.ReplaceAll(newPath, config.ArtifactIdOld, config.ArtifactIdNew)
			ArtifactIdOldUpper := strings.ToUpper(string(config.ArtifactIdOld[0])) + config.ArtifactIdOld[1:]
			ArtifactIdNewUpper := strings.ToUpper(string(config.ArtifactIdNew[0])) + config.ArtifactIdNew[1:]
			newPath = strings.ReplaceAll(newPath, ArtifactIdOldUpper, ArtifactIdNewUpper)

			//fmt.Println("newPath %s , path %s", path, newPath)
			err = writeFileTo(content, newPath, info)
			if err != nil {
				fmt.Println(err)
				return err
			}

		}
		return nil
	}
}

func reWriteFileContent(old []byte) []byte {
	oldStr := string(old)
	oldStr = strings.ReplaceAll(oldStr, config.GroupIdOld, config.GroupIdNew)
	oldStr = strings.ReplaceAll(oldStr, config.PackageNameOld, config.PackageNameNew)
	oldStr = strings.ReplaceAll(oldStr, config.ArtifactIdOld, config.ArtifactIdNew)
	ArtifactIdOldUpper := strings.ToUpper(string(config.ArtifactIdOld[0])) + config.ArtifactIdOld[1:]
	ArtifactIdNewUpper := strings.ToUpper(string(config.ArtifactIdNew[0])) + config.ArtifactIdNew[1:]
	oldStr = strings.ReplaceAll(oldStr, ArtifactIdOldUpper, ArtifactIdNewUpper)
	oldStr = strings.ReplaceAll(oldStr, config.TitleOld, config.TitleNew)
	// 遍历所有文件，白名单直接复制
	// 非白名单，重写文件内容，再生成文件
	//        return content.replaceAll(GROUP_ID, groupIdNew)
	//                .replaceAll(PACKAGE_NAME, packageNameNew)
	//                .replaceAll(ARTIFACT_ID, artifactIdNew) // 必须放在最后替换，因为 ARTIFACT_ID 太短！
	//                .replaceAll(StrUtil.upperFirst(ARTIFACT_ID), StrUtil.upperFirst(artifactIdNew))
	//                .replaceAll(TITLE, titleNew);
	return []byte(oldStr)
}

func writeFileTo(content []byte, path string, srcinfo os.FileInfo) error {
	err := os.MkdirAll(strings.Replace(path, srcinfo.Name(), "", -1), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("path: %s  mode: %s \r\n", path, srcinfo.Mode())
	return os.WriteFile(path, content, srcinfo.Mode())
}
