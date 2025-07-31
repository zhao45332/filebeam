package main

import (
	"log"

	"filebeam/config"
	"filebeam/server"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 创建服务器实例
	srv := server.NewServer(cfg)

	// 启动服务器
	log.Fatal(srv.Start())
}
