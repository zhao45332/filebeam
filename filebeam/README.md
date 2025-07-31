# FileBeam - 局域网文件共享工具

一个安全、稳定、现代化的局域网文件共享服务，支持文件上传、下载和管理。

## ✨ 功能特性

- 🔒 **安全上传** - 密码保护的文件上传功能
- 📁 **文件管理** - 直观的文件列表和详细信息
- 🚀 **现代化界面** - 响应式设计，支持移动端
- 📊 **实时进度** - 上传进度条和状态提示
- 🛡️ **安全检查** - 文件类型验证、大小限制、路径遍历防护
- 📱 **移动友好** - 完全响应式设计
- ⚡ **高性能** - 优化的文件处理和服务

## 🏗️ 项目结构

```
filebeam/
├── config/          # 配置管理
│   └── config.go
├── handlers/        # 请求处理器
│   ├── upload.go    # 上传处理
│   ├── download.go  # 下载处理
│   └── files.go     # 文件列表处理
├── server/          # 服务器核心
│   └── server.go
├── static/          # 前端资源
│   ├── index.html   # 主页面
│   ├── style.css    # 样式文件
│   └── app.js       # 前端逻辑
├── shared/          # 共享文件目录
├── main.go          # 程序入口
├── go.mod           # Go模块文件
└── README.md        # 项目文档
```

## 🚀 快速开始

### 环境要求

- Go 1.24 或更高版本
- 支持的操作系统：Windows、macOS、Linux

### 安装和运行

1. **克隆或下载项目**
   ```bash
   git clone <repository-url>
   cd filebeam
   ```

2. **运行服务**

   **方法一：直接运行（推荐）**
   ```bash
   # Windows
   start.bat
   
   # Linux/macOS
   chmod +x start.sh
   ./start.sh
   ```

   **方法二：手动编译运行**
   ```bash
   go build -o filebeam main.go
   ./filebeam
   ```

   **方法三：开发模式**
   ```bash
   go run main.go
   ```

3. **访问服务**
   - 本地访问：http://localhost:8888/
   - 局域网访问：http://[你的IP]:8888/

## ⚙️ 配置选项

通过环境变量可以自定义配置：

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `PORT` | `8888` | 服务端口 |
| `SHARED_DIR` | `./shared` | 共享文件目录 |
| `UPLOAD_PASSWORD` | `123456` | 上传密码 |
| `MAX_FILE_SIZE` | `104857600` | 最大文件大小（字节，默认100MB） |
| `ALLOWED_TYPES` | 空（允许所有） | 允许的文件类型（如：.jpg,.png,.pdf） |

### 配置示例

```bash
# Windows
set PORT=9000
set UPLOAD_PASSWORD=mypassword
set MAX_FILE_SIZE=524288000
go run main.go

# Linux/macOS
export PORT=9000
export UPLOAD_PASSWORD=mypassword
export MAX_FILE_SIZE=524288000
go run main.go
```

## 📖 使用说明

### 文件上传

1. 在网页界面选择要上传的文件
2. 输入上传密码（默认：123456）
3. 点击"上传文件"按钮
4. 等待上传完成

### 文件下载

1. 在文件列表中点击"下载"按钮
2. 文件将自动下载到本地

### 文件信息

1. 点击文件列表中的"详情"按钮
2. 查看文件的详细信息（大小、修改时间等）

## 🔧 API 接口

### 获取文件列表
```
GET /api/files
```

响应示例：
```json
{
  "files": [
    {
      "name": "example.pdf",
      "size": 1024000,
      "size_human": "1.0 MB",
      "modified": "2024-01-01T12:00:00Z",
      "modified_str": "2024-01-01 12:00:00",
      "download_url": "/download/example.pdf",
      "info_url": "/info/example.pdf"
    }
  ],
  "count": 1,
  "total_size": 1024000,
  "total_size_human": "1.0 MB"
}
```

### 文件上传
```
POST /upload
Content-Type: multipart/form-data

参数：
- file: 要上传的文件
- password: 上传密码
```

### 文件下载
```
GET /download/{filename}
```

### 文件信息
```
GET /info/{filename}
```

## 🛡️ 安全特性

- **密码保护** - 上传需要密码验证
- **文件类型验证** - 可配置允许的文件类型
- **大小限制** - 防止过大文件上传
- **路径遍历防护** - 防止恶意文件路径
- **文件名清理** - 自动清理不安全的文件名
- **MD5校验** - 文件完整性验证

## 🎨 界面特性

- **响应式设计** - 完美适配桌面和移动设备
- **现代化UI** - 美观的渐变背景和卡片式布局
- **实时反馈** - 上传进度条和状态提示
- **动画效果** - 流畅的交互动画
- **错误处理** - 友好的错误提示

## 🔄 更新日志

### v2.0.0
- ✨ 完全重构，采用模块化架构
- 🎨 全新的现代化界面设计
- 🔒 增强的安全功能
- 📱 完全响应式设计
- ⚡ 性能优化
- 📊 详细的文件信息显示

### v1.0.0
- 🚀 基础文件共享功能
- 📁 简单的文件列表
- ⬆️ 基础文件上传功能

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🆘 常见问题

**Q: 如何修改上传密码？**
A: 设置环境变量 `UPLOAD_PASSWORD` 或在 `config/config.go` 中修改默认值。

**Q: 如何限制文件大小？**
A: 设置环境变量 `MAX_FILE_SIZE`（以字节为单位）。

**Q: 如何只允许特定文件类型？**
A: 设置环境变量 `ALLOWED_TYPES`，例如：`.jpg,.png,.pdf`。

**Q: 服务无法启动怎么办？**
A: 检查端口是否被占用，可以修改 `PORT` 环境变量。

**Q: 文件上传失败？**
A: 检查文件大小是否超出限制，密码是否正确，以及磁盘空间是否充足。 