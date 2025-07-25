package main

import (
	"os"
)

func BackupCertsDir(srcDir, backupPath string) error {
	// 如果已存在备份，先删除
	if _, err := os.Stat(backupPath); err == nil {
		if err := os.RemoveAll(backupPath); err != nil {
			return err
		}
	}

	// 简单地将源目录重命名为备份路径
	return os.Rename(srcDir, backupPath)
}
