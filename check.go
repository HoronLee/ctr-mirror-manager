package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckConfig 检查当前镜像配置状态
func CheckConfig(certsDir string) error {
	// 检查目录是否存在
	info, err := os.Stat(certsDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("证书目录 %s 不存在", certsDir)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s 不是有效的目录", certsDir)
	}

	// 检查备份是否存在
	backupPath := certsDir + ".bak"
	_, backupErr := os.Stat(backupPath)
	if os.IsNotExist(backupErr) {
		fmt.Printf("备份目录 %s 不存在\n", backupPath)
	} else {
		fmt.Printf("已存在备份：%s\n", backupPath)
	}

	// 检查镜像配置
	var foundMirrors int
	err = filepath.Walk(certsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Base(path) == "hosts.toml" {
			relativePath, _ := filepath.Rel(certsDir, filepath.Dir(path))
			if relativePath == "." {
				relativePath = "<根目录>"
			}
			fmt.Printf("发现镜像配置：%s\n", relativePath)
			foundMirrors++
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("扫描目录失败: %v", err)
	}

	if foundMirrors == 0 {
		fmt.Println("未发现任何镜像配置")
	} else {
		fmt.Printf("共发现 %d 个镜像配置\n", foundMirrors)
	}

	return nil
}
