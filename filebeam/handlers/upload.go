package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"filebeam/config"
)

type UploadHandler struct {
	config *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{config: cfg}
}

func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST 上传", http.StatusMethodNotAllowed)
		return
	}

	// 验证密码
	password := r.FormValue("password")
	if password != h.config.UploadPassword {
		http.Error(w, "上传密码错误", http.StatusForbidden)
		return
	}

	// 获取上传的文件
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "读取上传文件失败", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 验证文件大小
	if header.Size > h.config.MaxFileSize {
		http.Error(w, fmt.Sprintf("文件过大，最大允许 %d MB", h.config.MaxFileSize/(1024*1024)), http.StatusRequestEntityTooLarge)
		return
	}

	// 验证文件类型
	if !h.isAllowedFileType(header.Filename) {
		http.Error(w, "不支持的文件类型", http.StatusBadRequest)
		return
	}

	// 生成安全的文件名
	safeFilename := h.generateSafeFilename(header.Filename)
	dstPath := filepath.Join(h.config.SharedDir, safeFilename)

	// 检查文件是否已存在
	if _, err := os.Stat(dstPath); err == nil {
		http.Error(w, "文件已存在，禁止重复上传", http.StatusConflict)
		return
	}

	// 创建目标文件
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "保存文件失败", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 计算MD5并写入文件
	hash := md5.New()
	teeReader := io.TeeReader(file, hash)

	_, err = io.Copy(dst, teeReader)
	if err != nil {
		// 删除部分写入的文件
		os.Remove(dstPath)
		http.Error(w, "写入文件失败", http.StatusInternalServerError)
		return
	}

	// 获取MD5值
	md5Hash := hex.EncodeToString(hash.Sum(nil))

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true, "filename": "%s", "md5": "%s", "size": %d}`, safeFilename, md5Hash, header.Size)
}

func (h *UploadHandler) isAllowedFileType(filename string) bool {
	if len(h.config.AllowedTypes) == 0 {
		return true // 允许所有类型
	}

	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range h.config.AllowedTypes {
		if strings.ToLower(allowedType) == ext {
			return true
		}
	}
	return false
}

func (h *UploadHandler) generateSafeFilename(originalName string) string {
	// 获取文件扩展名
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)

	// 清理文件名，只保留字母、数字、下划线和连字符
	var safeName strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-' {
			safeName.WriteRune(r)
		}
	}

	// 如果清理后为空，使用时间戳
	if safeName.Len() == 0 {
		safeName.WriteString(fmt.Sprintf("file_%d", time.Now().Unix()))
	}

	return safeName.String() + ext
}
