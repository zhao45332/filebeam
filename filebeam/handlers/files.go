package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"filebeam/config"
)

type FileInfo struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	SizeHuman   string    `json:"size_human"`
	Modified    time.Time `json:"modified"`
	ModifiedStr string    `json:"modified_str"`
	DownloadURL string    `json:"download_url"`
	InfoURL     string    `json:"info_url"`
}

type FilesHandler struct {
	config *config.Config
}

func NewFilesHandler(cfg *config.Config) *FilesHandler {
	return &FilesHandler{config: cfg}
}

func (h *FilesHandler) HandleFileList(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(h.config.SharedDir)
	if err != nil {
		http.Error(w, "无法读取目录", http.StatusInternalServerError)
		return
	}

	var fileList []FileInfo
	for _, f := range files {
		if !f.IsDir() {
			fileInfo, err := f.Info()
			if err != nil {
				continue // 跳过无法获取信息的文件
			}

			file := FileInfo{
				Name:        f.Name(),
				Size:        fileInfo.Size(),
				SizeHuman:   formatFileSizeHuman(fileInfo.Size()),
				Modified:    fileInfo.ModTime(),
				ModifiedStr: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
				DownloadURL: "/download/" + f.Name(),
				InfoURL:     "/info/" + f.Name(),
			}
			fileList = append(fileList, file)
		}
	}

	// 按修改时间排序，最新的在前
	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].Modified.After(fileList[j].Modified)
	})

	// 返回JSON格式的文件列表
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"files": fileList,
		"count": len(fileList),
		"total_size": func() int64 {
			var total int64
			for _, f := range fileList {
				total += f.Size
			}
			return total
		}(),
		"total_size_human": func() string {
			var total int64
			for _, f := range fileList {
				total += f.Size
			}
			return formatFileSizeHuman(total)
		}(),
	})
}

func formatFileSizeHuman(size int64) string {
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
