package config

import (
	"os"
	"path/filepath"
	"strings"
)

const defaultConfig = `linters:
  enable:
    - gocomment

linters-settings:
  gocomment:
    # 自定义配置项
`

func Setup(projectRoot string) error {
	configPath := filepath.Join(projectRoot, ".golangci.yml")

	// 如果配置文件不存在，直接创建
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return os.WriteFile(configPath, []byte(defaultConfig), 0644)
	}

	// 如果配置文件存在，合并配置
	return mergeConfig(configPath)
}

func mergeConfig(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	config := string(content)

	if !strings.Contains(config, "gocomment") {
		if strings.Contains(config, "linters:") {
			if strings.Contains(config, "enable:") {
				config = strings.Replace(config,
					"enable:",
					"enable:\n    - gocomment",
					1)
			} else {
				config = strings.Replace(config,
					"linters:",
					"linters:\n  enable:\n    - gocomment",
					1)
			}
		} else {
			config += "\n" + defaultConfig
		}

		return os.WriteFile(configPath, []byte(config), 0644)
	}

	return nil
}
