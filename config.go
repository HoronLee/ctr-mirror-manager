package main

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Host struct {
	URL          string   `toml:"url"`
	Capabilities []string `toml:"capabilities"`
	Username     string   `toml:"username,omitempty"`
	Password     string   `toml:"password,omitempty"`
	SkipVerify   bool     `toml:"skip_verify,omitempty"`
}

type Mirror struct {
	Name   string `toml:"name"`
	Server string `toml:"server"`
	Hosts  []Host `toml:"host"`
}

type Config struct {
	CertsDir string   `toml:"certs_dir,omitempty"`
	Mirrors  []Mirror `toml:"mirror"`
}

func ParseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = toml.Unmarshal(data, &cfg)

	// 如果未指定证书目录，使用默认值
	if cfg.CertsDir == "" {
		cfg.CertsDir = "/etc/containerd/certs.d"
	} else {
		// 将相对路径转换为绝对路径
		if !filepath.IsAbs(cfg.CertsDir) {
			configDir := filepath.Dir(path)
			cfg.CertsDir = filepath.Join(configDir, cfg.CertsDir)
		}
	}

	return &cfg, err
}
