package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reddd-cli/cmd/base"
	"reddd-cli/utils"
	"strings"
)

type PackageInfo struct {
	GroupId     string   `json:"groupId"`
	ArtifactId  string   `json:"artifactId"`
	PackageName string   `json:"packageName"`
	Title       string   `json:"title"`
	Ignore      []string `json:"ignore"`
	WhiteList   []string `json:"whiteList"`
}

func createFolder(project string) {
	// 创建文件夹
	if _, err := os.Stat(project); os.IsNotExist(err) {
		if err := os.MkdirAll(project, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	// 创建 package-info.json
	content := `{
	"groupId":"com.relengxing.demo",
	"artifactId":"demo",
	"packageName":"com.relengxing.demo",
	"title":"MyDemo"
}
	`
	packageTemplate := filepath.Join(project, "package-template.json")
	os.WriteFile(packageTemplate, []byte(content), os.ModePerm)
}

func clone(ctx context.Context) (PackageInfo, *base.Repo) {
	repo := base.NewRepo(repoURL, branch)
	repo.Clone(ctx)
	packageInfo := filepath.Join(repo.Path(), "package-info.json")
	byteValue, _ := os.ReadFile(packageInfo)
	var packageInfoRepo PackageInfo
	_ = json.Unmarshal(byteValue, &packageInfoRepo)
	return packageInfoRepo, repo
}

// 根据文件夹内的 Package-info 信息生成标准工程
func generate(ctx context.Context, project string) error {
	if _, err := os.Stat(project); os.IsNotExist(err) {
		return errors.New("文件夹不存在，请先创建文件夹")
	}
	packageTemplate := filepath.Join(project, "package-template.json")
	byteValue, _ := os.ReadFile(packageTemplate)
	var packageInfoNew PackageInfo
	_ = json.Unmarshal(byteValue, &packageInfoNew)
	infoRepo, repo := clone(ctx)

	if err := filepath.Walk(repo.Path(), walk(repo.Path(), project, infoRepo, packageInfoNew)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func walk(home string, to string, infoRepo, infoNew PackageInfo) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && utils.ContainsString(infoRepo.Ignore, info.Name()) {
			return filepath.SkipDir
		}

		if !info.IsDir() { // 如果当前路径不是一个目录，则将其添加到文件路径数组中
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if !utils.ContainsString(infoRepo.WhiteList, filepath.Ext(path)) {
				// rewrite file
				content = reWriteFileContent(content, infoRepo, infoNew)
			}
			//return file.getPath().replace(projectBaseDir, projectBaseDirNew) // 新目录
			//.replace(PACKAGE_NAME.replaceAll("\\.", Matcher.quoteReplacement(separator)),
			//packageNameNew.replaceAll("\\.", Matcher.quoteReplacement(separator)))
			//.replace(ARTIFACT_ID, artifactIdNew) //
			//.replaceAll(StrUtil.upperFirst(ARTIFACT_ID), StrUtil.upperFirst(artifactIdNew));
			newPath := strings.ReplaceAll(path, home, to)
			newPath = strings.ReplaceAll(newPath, strings.ReplaceAll(infoRepo.PackageName, "\\.", string(filepath.Separator)), strings.ReplaceAll(infoNew.PackageName, "\\.", string(filepath.Separator)))
			newPath = strings.ReplaceAll(newPath, infoRepo.ArtifactId, infoNew.ArtifactId)
			ArtifactIdRepoUpper := strings.ToUpper(string(infoRepo.ArtifactId[0])) + infoRepo.ArtifactId[1:]
			ArtifactIdNewUpper := strings.ToUpper(string(infoNew.ArtifactId[0])) + infoNew.ArtifactId[1:]
			newPath = strings.ReplaceAll(newPath, ArtifactIdRepoUpper, ArtifactIdNewUpper)

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

func reWriteFileContent(old []byte, infoRepo, infoNew PackageInfo) []byte {
	oldStr := string(old)
	oldStr = strings.ReplaceAll(oldStr, infoRepo.GroupId, infoNew.GroupId)
	oldStr = strings.ReplaceAll(oldStr, infoRepo.PackageName, infoNew.PackageName)
	oldStr = strings.ReplaceAll(oldStr, infoRepo.ArtifactId, infoNew.ArtifactId)
	ArtifactIdRepoUpper := strings.ToUpper(string(infoRepo.ArtifactId[0])) + infoRepo.ArtifactId[1:]
	ArtifactIdNewUpper := strings.ToUpper(string(infoNew.ArtifactId[0])) + infoNew.ArtifactId[1:]
	oldStr = strings.ReplaceAll(oldStr, ArtifactIdRepoUpper, ArtifactIdNewUpper)
	oldStr = strings.ReplaceAll(oldStr, infoRepo.Title, infoNew.Title)
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
	return os.WriteFile(path, content, srcinfo.Mode())
}