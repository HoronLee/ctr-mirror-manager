package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ApplyConfig(cfg *Config, baseDir string) error {
	for _, mirror := range cfg.Mirrors {
		targetDir := filepath.Join(baseDir, mirror.Name)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}
		targetFile := filepath.Join(targetDir, "hosts.toml")
		content, err := RenderHostsToml(mirror)
		if err != nil {
			return err
		}
		if err := os.WriteFile(targetFile, []byte(content), 0644); err != nil {
			return err
		}
		fmt.Printf("✅已更新: %s\n", targetFile)
	}
	return nil
}

func RenderHostsToml(m Mirror) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("server = %q\n\n", m.Server))
	for _, h := range m.Hosts {
		sb.WriteString(fmt.Sprintf("[host.%q]\n", h.URL))
		sb.WriteString("  capabilities = [")
		for i, cap := range h.Capabilities {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", cap))
		}
		sb.WriteString("]\n")
		if h.Username != "" {
			sb.WriteString(fmt.Sprintf("  username = %q\n", h.Username))
		}
		if h.Password != "" {
			sb.WriteString(fmt.Sprintf("  password = %q\n", h.Password))
		}
		// 添加 skip_verify 选项处理
		if h.SkipVerify {
			sb.WriteString("  skip_verify = true\n")
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
