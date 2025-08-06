package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"filebeam/config"
	"filebeam/server"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 如果未设置共享目录，则交互式询问用户
	if cfg.SharedDir == "" {
		cfg.SharedDir = promptForSharedDirectory()
	}

	// 创建服务器实例
	srv := server.NewServer(cfg)

	// 启动服务器
	log.Fatal(srv.Start())
}

func promptForSharedDirectory() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("FileBeam 文件共享服务")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("请输入要共享的文件夹路径：")
	fmt.Println("示例：")
	fmt.Println("  Windows: C:\\Users\\YourName\\Documents\\SharedFiles")
	fmt.Println("  Linux/Mac: /home/user/Documents/SharedFiles")
	fmt.Println("  Windows: D:\\MyFiles")
	fmt.Println()

	for {
		fmt.Print("共享文件夹路径: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取输入失败: %v\n", err)
			continue
		}

		// 清理输入
		path := strings.TrimSpace(input)
		if path == "" {
			fmt.Println("❌ 路径不能为空，请重新输入")
			continue
		}

		// 检查路径是否存在
		fileInfo, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("❌ 路径不存在: %s\n", path)
				fmt.Print("是否要创建这个文件夹？(y/n): ")
				create, _ := reader.ReadString('\n')
				if strings.ToLower(strings.TrimSpace(create)) == "y" {
					if err := os.MkdirAll(path, os.ModePerm); err != nil {
						fmt.Printf("❌ 创建文件夹失败: %v\n", err)
						continue
					}
					fmt.Printf("✅ 已创建文件夹: %s\n", path)
				} else {
					continue
				}
			} else {
				fmt.Printf("❌ 无法访问路径: %v\n", err)
				continue
			}
		} else {
			// 路径存在，检查是否是目录
			if !fileInfo.IsDir() {
				fmt.Printf("❌ 指定的路径不是文件夹: %s\n", path)
				continue
			}
		}

		// 检查目录权限
		if err := checkDirectoryPermissions(path); err != nil {
			fmt.Printf("❌ 权限检查失败: %v\n", err)
			continue
		}

		// 转换为绝对路径
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Printf("❌ 获取绝对路径失败: %v\n", err)
			continue
		}

		fmt.Printf("✅ 共享文件夹: %s\n", absPath)
		fmt.Println()
		return absPath
	}
}

func checkDirectoryPermissions(path string) error {
	// 尝试在目录中创建一个临时文件来测试写权限
	testFile := filepath.Join(path, ".filebeam_test")

	// 创建测试文件
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("无法在文件夹中创建文件，请检查写权限")
	}
	file.Close()

	// 删除测试文件
	os.Remove(testFile)

	return nil
}
