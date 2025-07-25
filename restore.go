package main

import (
	"fmt"
	"os"
)

// RestoreBackup 从备份恢复证书目录
func RestoreBackup(backupPath, targetDir string) error {
	// 检查备份是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份目录 %s 不存在", backupPath)
	}

	// 如果目标目录存在，先删除
	if _, err := os.Stat(targetDir); err == nil {
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("清理目标目录失败: %v", err)
		}
	}

	// 从备份恢复
	return os.Rename(backupPath, targetDir)
}
