package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Path     string `mapstructure:"path"`
	Interval int    `mapstructure:"interval"`
}

func main() {
	// 初始化 viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var apps []AppConfig
	if err := viper.UnmarshalKey("apps", &apps); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	for _, app := range apps {
		app := app
		go monitorApp(app)
	}

	select {} // 阻塞主 goroutine
}

// 监控和启动程序
func monitorApp(cfg AppConfig) {
	for {
		if !isProcessRunning(cfg.Path) {
			log.Printf("[%s] 未运行，正在启动...\n", cfg.Path)
			if err := startProcess(cfg.Path); err != nil {
				log.Printf("[%s] 启动失败: %v\n", cfg.Path, err)
			} else {
				log.Printf("[%s] 启动成功\n", cfg.Path)
			}
		} else {
			log.Printf("[%s] 已在运行\n", cfg.Path)
		}
		time.Sleep(time.Duration(cfg.Interval) * time.Second)
	}
}

// 判断进程是否已运行
func isProcessRunning(path string) bool {
	exeName := filepath.Base(path)
	// 在 Windows 下使用 tasklist 检查
	out, err := exec.Command("tasklist").Output()
	if err != nil {
		log.Printf("执行 tasklist 出错: %v", err)
		return false
	}
	//fmt.Println(string(out))
	return strings.Contains(string(out), exeName)
}

// 启动进程
func startProcess(path string) error {
	cmd := exec.Command(path)
	return cmd.Start() // 不等待完成
}
