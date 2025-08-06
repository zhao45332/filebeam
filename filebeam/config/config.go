package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	SharedDir      string
	UploadPassword string
	MaxFileSize    int64    // 最大文件大小（字节）
	AllowedTypes   []string // 允许的文件类型
}

func LoadConfig() *Config {
	config := &Config{
		Port:           getEnv("PORT", "8888"),
		SharedDir:      getEnv("SHARED_DIR", ""), // 默认为空，需要用户指定
		UploadPassword: getEnv("UPLOAD_PASSWORD", "123456"),
		MaxFileSize:    getEnvAsInt64("MAX_FILE_SIZE", 100*1024*1024), // 默认100MB
		AllowedTypes:   []string{},                                    // 空表示允许所有类型
	}

	// 如果设置了允许的文件类型
	if allowedTypes := os.Getenv("ALLOWED_TYPES"); allowedTypes != "" {
		config.AllowedTypes = []string{allowedTypes}
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
