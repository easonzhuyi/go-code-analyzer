package hook

import (
    "fmt"
    "os"
    "path/filepath"
)

const preCommitScript = `#!/bin/sh
# 获取暂存区中的 Go 文件
files=$(git diff --cached --name-only --diff-filter=d | grep "\.go$" || true)
if [ -n "$files" ]; then
    golangci-lint run $files
    if [ $? -ne 0 ]; then
        echo "代码检查失败，请修复问题后重试"
        exit 1
    fi
fi
`

type PreCommit struct {
    projectRoot string
}

func NewPreCommit(projectRoot string) *PreCommit {
    return &PreCommit{
        projectRoot: projectRoot,
    }
}

func (p *PreCommit) Install() error {
    hooksDir := filepath.Join(p.projectRoot, ".git", "hooks")
    hookPath := filepath.Join(hooksDir, "pre-commit")

    if err := os.MkdirAll(hooksDir, 0755); err != nil {
        return fmt.Errorf("failed to create hooks directory: %w", err)
    }

    if exists, _ := fileExists(hookPath); exists {
        backupPath := hookPath + ".backup"
        if err := os.Rename(hookPath, backupPath); err != nil {
            return fmt.Errorf("failed to backup existing hook: %w", err)
        }
        fmt.Printf("已备份现有 hook 到: %s\n", backupPath)
    }

    if err := os.WriteFile(hookPath, []byte(preCommitScript), 0755); err != nil {
        return fmt.Errorf("failed to write hook script: %w", err)
    }

    fmt.Printf("成功安装 pre-commit hook 到: %s\n", hookPath)
    return nil
}

func fileExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}
