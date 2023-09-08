/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package project

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// projectCmd represents the project command
var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "项目相关",
	Long: `项目相关命令:

	1. 创建项目`,

	Run: run,
}

var (
	repoURL string
	branch  string
	timeout string
)

func init() {
	if repoURL = os.Getenv("REDDD_LAYOUT_REPO"); repoURL == "" {
		repoURL = "git@github.com:relengxing/reddd.git"
	}
	timeout = "60s"
	ProjectCmd.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	ProjectCmd.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
	ProjectCmd.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	ctx = ctx
	defer cancel()
	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is project name ?",
			Help:    "Created project name.",
		}
		err = survey.AskOne(prompt, &name)
		if err != nil || name == "" {
			return
		}
	} else {
		name = args[0]
	}
	projectName, workingDir := processProjectParams(name, wd)
	p := &Project{Name: projectName}
	fmt.Println(projectName, workingDir)

	p.New(ctx, workingDir, repoURL, branch)

	//done := make(chan error, 1)
	//go func() {
	//	projectRoot := workingDir
	//	packagePath, e := filepath.Rel(projectRoot, filepath.Join(workingDir, projectName))
	//	if e != nil {
	//		done <- fmt.Errorf("🚫 failed to get relative path: %v", err)
	//		return
	//	}
	//	packagePath = strings.ReplaceAll(packagePath, "\\", "/")
	//
	//	mod, e := base.ModulePath(filepath.Join(projectRoot, "go.mod"))
	//	if e != nil {
	//		done <- fmt.Errorf("🚫 failed to parse `go.mod`: %v", e)
	//		return
	//	}
	//	// Get the relative path for adding a project based on Go modules
	//	p.Path = filepath.Join(strings.TrimPrefix(workingDir, projectRoot+"/"), p.Name)
	//	done <- p.Add(ctx, workingDir, repoURL, branch, mod, packagePath)
	//}()
	//select {
	//case <-ctx.Done():
	//	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
	//		fmt.Fprint(os.Stderr, "\033[31mERROR: project creation timed out\033[m\n")
	//		return
	//	}
	//	fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to create project(%s)\033[m\n", ctx.Err().Error())
	//	//case err = <-done:
	//	//	if err != nil {
	//	//		fmt.Fprintf(os.Stderr, "\033[31mERROR: Failed to create project(%s)\033[m\n", err.Error())
	//	//	}
	//}
}

func processProjectParams(projectName string, workingDir string) (projectNameResult, workingDirResult string) {
	_projectDir := projectName
	_workingDir := workingDir
	// Process ProjectName with system variable
	if strings.HasPrefix(projectName, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// cannot get user home return fallback place dir
			return _projectDir, _workingDir
		}
		_projectDir = filepath.Join(homeDir, projectName[2:])
	}

	// check path is relative
	if !filepath.IsAbs(projectName) {
		absPath, err := filepath.Abs(projectName)
		if err != nil {
			return _projectDir, _workingDir
		}
		_projectDir = absPath
	}

	return filepath.Base(_projectDir), filepath.Dir(_projectDir)
}
