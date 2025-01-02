package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/easonzhuyi/go-code-analyzer/debug"
	"github.com/easonzhuyi/go-code-analyzer/internal/analyzer"
	"github.com/easonzhuyi/go-code-analyzer/internal/config"
	"github.com/easonzhuyi/go-code-analyzer/internal/hook"
)

func main() {
	if isRunningAsPlugin() {
		analyzer.RunAsPlugin()
		return
	}

	if err := autoInstall(); err != nil {
		fmt.Fprintf(os.Stderr, "Installation failed: %v\n", err)
		os.Exit(1)
	}
}

func isRunningAsPlugin() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "golangci-lint") {
			return true
		}
	}
	return false
}

func autoInstall() error {
	defer debug.NewTimer("autoInstall").Stop()

	debug.DumpEnv()

	debug.Log("Checking golangci-lint installation...")
	// 1. 确保 golangci-lint 已安装
	if err := ensureGolangciLint(); err != nil {
		return err
	}

	// 2. 获取项目根目录
	debug.Log("Getting project root...")
	projectRoot, err := getProjectRoot()
	if err != nil {
		return err
	}
	debug.Log("Project root: %s", projectRoot)

	// 3. 安装 pre-commit hook
	debug.Log("Installing pre-commit hook...")
	precommit := hook.NewPreCommit(projectRoot)
	if err := precommit.Install(); err != nil {
		return fmt.Errorf("failed to install pre-commit hook: %w", err)
	}

	// 4. 创建或更新 golangci-lint 配置
	debug.Log("Setting up golangci-lint config...")
	if err := config.Setup(projectRoot); err != nil {
		return fmt.Errorf("failed to setup golangci-lint config: %w", err)
	}

	return nil
}

func ensureGolangciLint() error {
	if err := exec.Command("golangci-lint", "--version").Run(); err != nil {
		fmt.Println("Installing golangci-lint...")

		cmd := exec.Command("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install golangci-lint: %w", err)
		}

		time.Sleep(time.Second)
		if err := exec.Command("golangci-lint", "--version").Run(); err != nil {
			return fmt.Errorf("golangci-lint installation verification failed: %w", err)
		}

		fmt.Println("✅ golangci-lint installed successfully")
	}
	return nil
}

func getProjectRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get project root: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func init() {
	if err := ensureGopathInPath(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to setup PATH: %v\n", err)
	}
}

func ensureGopathInPath() error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		gopath = filepath.Join(homeDir, "go")
	}

	gopathBin := filepath.Join(gopath, "bin")
	path := os.Getenv("PATH")

	if !strings.Contains(path, gopathBin) {
		newPath := fmt.Sprintf("%s%s%s", gopathBin, string(os.PathListSeparator), path)
		os.Setenv("PATH", newPath)
	}

	return nil
}
