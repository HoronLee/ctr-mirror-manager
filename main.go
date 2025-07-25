package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "v1.0.1"
)

// 应用命令相关的配置
type applyOptions struct {
	configPath string
	backupPath string
}

// 恢复命令相关的配置
type restoreOptions struct {
	backupPath string
	configPath string // 添加配置文件路径
}

// 检查命令相关的配置
type checkOptions struct {
	configPath string // 添加配置文件路径
}

// 创建根命令
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "cmm",
		Short:        "Containerd 镜像源批量配置工具",
		SilenceUsage: true,
	}

	// 添加子命令
	rootCmd.AddCommand(newApplyCmd())
	rootCmd.AddCommand(newRestoreCmd())
	rootCmd.AddCommand(newCheckCmd())
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}

// 创建应用配置命令
func newApplyCmd() *cobra.Command {
	opts := &applyOptions{}

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "✍️应用镜像源配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApply(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.configPath, "config", "c", "", "📄镜像源配置 TOML 文件路径")
	cmd.MarkFlagRequired("config")

	return cmd
}

// 创建恢复备份命令
func newRestoreCmd() *cobra.Command {
	opts := &restoreOptions{
		backupPath: "",
	}

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "🔄从备份恢复镜像源配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestore(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.configPath, "config", "c", "", "📄镜像源配置 TOML 文件路径 (可选，用于读取 certs_dir)")

	return cmd
}

// 创建检查命令
func newCheckCmd() *cobra.Command {
	opts := &checkOptions{
		configPath: "",
	}

	cmd := &cobra.Command{
		Use:   "check",
		Short: "检查当前镜像配置状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.configPath, "config", "c", "", "📄镜像源配置 TOML 文件路径 (可选，用于读取 certs_dir)")

	return cmd
}

// 执行检查命令
func runCheck(opts *checkOptions) error {
	var certsDir string

	// 如果提供了配置文件，从配置读取证书目录
	if opts.configPath != "" {
		cfg, err := ParseConfig(opts.configPath)
		if err != nil {
			return fmt.Errorf("❌解析配置失败: %v", err)
		}
		certsDir = cfg.CertsDir
	} else {
		// 使用默认证书目录
		certsDir = "/etc/containerd/certs.d"
	}

	fmt.Printf("🔍正在检查 %s 目录的配置...\n", certsDir)
	return CheckConfig(certsDir)
}

// 执行恢复备份命令
func runRestore(opts *restoreOptions) error {
	var certsDir string

	// 如果提供了配置文件，从配置读取证书目录
	if opts.configPath != "" {
		cfg, err := ParseConfig(opts.configPath)
		if err != nil {
			return fmt.Errorf("❌解析配置失败: %v", err)
		}
		certsDir = cfg.CertsDir
	} else {
		// 使用默认证书目录
		certsDir = "/etc/containerd/certs.d"
	}

	opts.backupPath = certsDir + ".bak"

	// 检查备份是否存在
	if _, err := os.Stat(opts.backupPath); os.IsNotExist(err) {
		return fmt.Errorf("❌备份目录 %s 不存在，无法恢复", opts.backupPath)
	}

	fmt.Printf("🔁正在从 %s 恢复备份...\n", opts.backupPath)
	if err := RestoreBackup(opts.backupPath, certsDir); err != nil {
		return fmt.Errorf("❌恢复失败: %v", err)
	}

	fmt.Println("✅备份恢复成功")
	return nil
}

// 执行应用配置命令
func runApply(opts *applyOptions) error {
	// 解析配置文件
	cfg, err := ParseConfig(opts.configPath)
	if err != nil {
		return fmt.Errorf("❌解析配置失败: %v", err)
	}

	// 使用配置文件中指定的证书目录
	certsDir := cfg.CertsDir
	// 检查 certsDir 是否存在，如果不存在则创建
	if _, err := os.Stat(certsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(certsDir, 0755); err != nil {
			return fmt.Errorf("❌创建目录失败: %v", err)
		}
		fmt.Printf("⚠️目录 %s 不存在，已自动创建\n", certsDir)
	}

	opts.backupPath = certsDir + ".bak"
	// 1. 备份
	if err := BackupCertsDir(certsDir, opts.backupPath); err != nil {
		return fmt.Errorf("❌备份失败: %v", err)
	}

	// 重新创建证书目录
	if err := os.MkdirAll(certsDir, 0755); err != nil {
		return fmt.Errorf("❌重新创建目录失败: %v", err)
	}

	// 3. 生成 hosts.toml
	if err := ApplyConfig(cfg, certsDir); err != nil {
		return fmt.Errorf("❌写入配置失败: %v", err)
	}

	fmt.Println("✅操作成功，镜像源已更新。")
	return nil
}

// 创建版本命令
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cmm version:", Version)
		},
	}
}

func main() {
	rootCmd := newRootCmd()
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
