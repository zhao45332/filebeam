package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"filebeam/config"
)

type DownloadHandler struct {
	config *config.Config
}

func NewDownloadHandler(cfg *config.Config) *DownloadHandler {
	return &DownloadHandler{config: cfg}
}

func (h *DownloadHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	// 从URL中提取文件路径
	relPath := r.URL.Path[len("/download/"):]
	if relPath == "" {
		http.Error(w, "文件路径不能为空", http.StatusBadRequest)
		return
	}

	// 安全检查：防止路径遍历攻击
	if strings.Contains(relPath, "..") || strings.Contains(relPath, "/") {
		http.Error(w, "非法的文件路径", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.config.SharedDir, relPath)

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "文件不存在", http.StatusNotFound)
		} else {
			http.Error(w, "无法访问文件", http.StatusInternalServerError)
		}
		return
	}

	// 检查是否为目录
	if fileInfo.IsDir() {
		http.Error(w, "不能下载目录", http.StatusBadRequest)
		return
	}

	// 设置响应头
	filename := filepath.Base(filePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Last-Modified", fileInfo.ModTime().Format(time.RFC1123))

	// 提供文件下载
	http.ServeFile(w, r, filePath)
}

func (h *DownloadHandler) HandleFileInfo(w http.ResponseWriter, r *http.Request) {
	// 从URL中提取文件路径
	relPath := r.URL.Path[len("/info/"):]
	if relPath == "" {
		http.Error(w, "文件路径不能为空", http.StatusBadRequest)
		return
	}

	// 安全检查
	if strings.Contains(relPath, "..") || strings.Contains(relPath, "/") {
		http.Error(w, "非法的文件路径", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.config.SharedDir, relPath)

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "文件不存在", http.StatusNotFound)
		} else {
			http.Error(w, "无法访问文件", http.StatusInternalServerError)
		}
		return
	}

	// 返回文件信息
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"name": "%s",
		"size": %d,
		"modified": "%s",
		"size_human": "%s"
	}`,
		fileInfo.Name(),
		fileInfo.Size(),
		fileInfo.ModTime().Format(time.RFC3339),
		formatFileSize(fileInfo.Size()))
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
