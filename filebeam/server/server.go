package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"filebeam/config"
	"filebeam/handlers"
)

type Server struct {
	config          *config.Config
	uploadHandler   *handlers.UploadHandler
	downloadHandler *handlers.DownloadHandler
	filesHandler    *handlers.FilesHandler
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config:          cfg,
		uploadHandler:   handlers.NewUploadHandler(cfg),
		downloadHandler: handlers.NewDownloadHandler(cfg),
		filesHandler:    handlers.NewFilesHandler(cfg),
	}
}

func (s *Server) Start() error {
	// 验证共享目录配置
	if s.config.SharedDir == "" {
		return fmt.Errorf("未指定共享目录，请设置 SHARED_DIR 环境变量或通过命令行参数指定")
	}

	// 检查共享目录是否存在
	fileInfo, err := os.Stat(s.config.SharedDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("共享目录不存在: %s", s.config.SharedDir)
		}
		return fmt.Errorf("无法访问共享目录: %v", err)
	}

	// 检查是否是目录
	if !fileInfo.IsDir() {
		return fmt.Errorf("指定的路径不是目录: %s", s.config.SharedDir)
	}

	// 检查目录权限
	if err := s.checkDirectoryPermissions(); err != nil {
		return fmt.Errorf("目录权限检查失败: %v", err)
	}

	// 设置路由
	s.setupRoutes()

	// 打印启动信息
	s.printStartupInfo()

	// 启动服务器
	return http.ListenAndServe("0.0.0.0:"+s.config.Port, nil)
}

func (s *Server) setupRoutes() {
	// 静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API路由
	http.HandleFunc("/api/files", s.filesHandler.HandleFileList)
	http.HandleFunc("/upload", s.uploadHandler.HandleUpload)
	http.HandleFunc("/download/", s.downloadHandler.HandleDownload)
	http.HandleFunc("/info/", s.downloadHandler.HandleFileInfo)

	// 首页路由
	http.HandleFunc("/", s.handleIndex)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// 只处理根路径
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// 提供静态HTML页面
	http.ServeFile(w, r, "static/index.html")
}

func (s *Server) printStartupInfo() {
	log.Println("✅ FileBeam 文件共享服务启动成功！")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📁 共享目录: %s", s.config.SharedDir)
	log.Printf("🔒 上传密码: %s", s.config.UploadPassword)
	log.Printf("📏 最大文件: %d MB", s.config.MaxFileSize/(1024*1024))
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🌐 访问地址:")
	log.Printf("   本地访问: http://localhost:%s/", s.config.Port)

	// 获取可访问的局域网地址
	s.printAccessibleAddresses()

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🚀 服务已就绪，等待连接...")
}

func (s *Server) printAccessibleAddresses() {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range interfaces {
		// 跳过回环接口和未启用的接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				ip := ipnet.IP
				// 只显示IPv4地址，排除回环地址
				if ip.To4() != nil && !ip.IsLoopback() {
					// 检查是否是私有网络地址（可访问的局域网地址）
					if s.isPrivateNetwork(ip) {
						log.Printf("   局域网访问: http://%s:%s/", ip.String(), s.config.Port)
					}
				}
			}
		}
	}
}

func (s *Server) isPrivateNetwork(ip net.IP) bool {
	// 检查是否是私有网络地址
	// 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16

	// 排除链路本地地址 (169.254.0.0/16)
	if ip[0] == 169 && ip[1] == 254 {
		return false
	}

	// 排除回环地址 (127.0.0.0/8)
	if ip[0] == 127 {
		return false
	}

	// 检查私有网络地址
	if ip[0] == 10 {
		return true
	}
	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return true
	}
	if ip[0] == 192 && ip[1] == 168 {
		return true
	}

	return false
}

func (s *Server) checkDirectoryPermissions() error {
	// 尝试在目录中创建一个临时文件来测试写权限
	testFile := filepath.Join(s.config.SharedDir, ".filebeam_test")

	// 创建测试文件
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("无法在共享目录中创建文件，请检查写权限")
	}
	file.Close()

	// 删除测试文件
	os.Remove(testFile)

	return nil
}
