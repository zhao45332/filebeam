package server

import (
	"log"
	"net"
	"net/http"
	"os"

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
	// 确保共享目录存在
	if err := os.MkdirAll(s.config.SharedDir, os.ModePerm); err != nil {
		return err
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

	// 获取本机IP地址
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				log.Printf("   局域网访问: http://%s:%s/", ipnet.IP.String(), s.config.Port)
			}
		}
	}
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🚀 服务已就绪，等待连接...")
}
