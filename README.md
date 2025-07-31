
# 📦 FileBeam - 简洁易用的局域网文件共享工具

`FileBeam` 是一个基于 Go 开发的轻量级文件共享服务，支持：

- ✅ 浏览器文件上传（支持进度条）
- ✅ 密码保护上传接口
- ✅ 文件下载列表页面
- ✅ 多人局域网访问共享内容
- ✅ 防止重复上传
- ✅ 文件上传取消与进度展示
- ✅ 文件夹上传（支持结构保留，可拓展）

---

## 🚀 快速开始

### 1. 构建可执行文件（Windows 示例）：

```bash
go build -o filebeam.exe
```

### 2. 运行服务：

```bash
./filebeam.exe
```

输出示例：

```
✅ FileBeam 文件共享服务启动：
  本地访问:     http://localhost:8888/
  局域网访问:   http://192.168.1.5:8888/
```

---

## 📂 默认共享目录

程序会在当前目录自动创建一个名为 `shared` 的目录：

```
./shared
```

你上传的文件都会保存到这里。

---

## 📤 上传文件说明

- 访问主页面 [http://localhost:8888](http://localhost:8888)
- 填入正确的上传密码（默认为 `123456`）
- 支持上传取消、上传进度展示
- 如果文件重名，将拒绝上传（防止覆盖）

---

## 🛡️ 安全特性

| 功能                  | 描述                     |
|-----------------------|--------------------------|
| 上传密码保护          | 防止任意人上传文件       |
| 拒绝重复上传          | 不允许文件名冲突上传     |
| 取消上传              | 在上传过程中可随时取消   |
| 文件类型限制（可拓展）| 目前未启用，可自定义扩展 |

---

## 📁 支持文件夹上传（拓展）

> 要实现文件夹上传，请将前端 `<input type="file">` 替换为：

```html
<input type="file" name="file" webkitdirectory directory multiple>
```

并搭配 `file.webkitRelativePath` 在后端处理完整路径上传。详情可查看项目扩展说明。

---

## 🛠️ 开发者指南

- 使用 Go 原生 `net/http`
- HTML 页面为内嵌模板，支持前端上传进度展示
- 支持局域网自动打印访问地址
- 可拓展支持文件夹压缩打包下载、Token认证等高级功能

---

## 🔒 修改上传密码

在源码中修改以下变量即可：

```go
const uploadPassword = "123456"
```

建议配合 nginx 或 frp 进行公网访问代理，并加上 HTTPS 保护。

---

## 📌 依赖环境

- Go 1.18+
- 支持 Windows / Linux / macOS

---

## 📷 示例截图

> 主页面展示上传进度与文件列表

![upload-example](https://dummyimage.com/600x150/ddd/000.png&text=FileBeam+上传界面+示意图)

---

## 📄 License

MIT License © 2025 FileBeam Contributors
