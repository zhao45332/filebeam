package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

const (
	sharedDir      = "./shared"
	uploadPassword = "123456" // ✅ 上传密码
)

func main() {
	_ = os.MkdirAll(sharedDir, os.ModePerm)

	// 上传接口
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "只支持 POST 上传", http.StatusMethodNotAllowed)
			return
		}

		password := r.FormValue("password")
		if password != uploadPassword {
			http.Error(w, "上传密码错误", http.StatusForbidden)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "读取上传文件失败", http.StatusBadRequest)
			return
		}
		defer file.Close()

		dstPath := filepath.Join(sharedDir, header.Filename)

		// ✅ 防止重复上传
		if _, err := os.Stat(dstPath); err == nil {
			http.Error(w, "文件已存在，禁止重复上传", http.StatusConflict)
			return
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "保存文件失败", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = dst.ReadFrom(file)
		if err != nil {
			http.Error(w, "写入文件失败", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// 下载接口
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		relPath := r.URL.Path[len("/download/"):]
		filePath := filepath.Join(sharedDir, relPath)

		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, filePath)
	})

	// 首页页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir(sharedDir)
		if err != nil {
			http.Error(w, "无法读取目录", http.StatusInternalServerError)
			return
		}

		var fileList []string
		for _, f := range files {
			if !f.IsDir() {
				fileList = append(fileList, f.Name())
			}
		}

		tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>FileBeam - 局域网共享工具</title>
</head>
<body>
	<h2>📂 共享文件列表</h2>
	<ul>
		{{range .}}
			<li><a href="/download/{{.}}">{{.}}</a></li>
		{{else}}
			<li>没有可用文件</li>
		{{end}}
	</ul>
	<hr>
	<h3>⬆ 上传文件（需密码）</h3>
	<form id="uploadForm" enctype="multipart/form-data">
		<input type="file" name="file" required>
		<br><br>
		密码：<input type="password" name="password" required>
		<br><br>
		<input type="submit" value="上传" id="submitBtn">
		<button type="button" id="cancelBtn" disabled>取消上传</button>
	</form>

	<!-- 进度条 -->
	<div id="progressContainer" style="width: 300px; background: #eee; margin-top: 20px; display: none;">
		<div id="progressBar" style="width: 0%; height: 20px; background: green;"></div>
	</div>
	<div id="progressText" style="margin-top: 5px;"></div>

	<script>
		let xhr = null;
		const form = document.getElementById('uploadForm');
		const progressBar = document.getElementById('progressBar');
		const progressContainer = document.getElementById('progressContainer');
		const progressText = document.getElementById('progressText');
		const cancelBtn = document.getElementById('cancelBtn');
		const submitBtn = document.getElementById('submitBtn');

		form.addEventListener('submit', function(event) {
			event.preventDefault();

			submitBtn.disabled = true;
			cancelBtn.disabled = false;

			const formData = new FormData(form);
			xhr = new XMLHttpRequest();

			xhr.upload.addEventListener("progress", function(e) {
				if (e.lengthComputable) {
					const percent = (e.loaded / e.total) * 100;
					progressBar.style.width = percent.toFixed(2) + "%";
					progressText.innerText = "上传进度：" + percent.toFixed(2) + "%";
				}
			});

			xhr.onload = function() {
				submitBtn.disabled = false;
				cancelBtn.disabled = true;

				if (xhr.status === 200 || xhr.status === 303) {
					progressText.innerText = "✅ 上传成功，正在刷新...";
					setTimeout(() => location.reload(), 1000);
				} else if (xhr.status === 409) {
					progressText.innerText = "⚠️ 文件已存在，禁止重复上传";
				} else {
					progressText.innerText = "❌ 上传失败：" + xhr.responseText;
				}
			};

			xhr.onerror = function() {
				progressText.innerText = "❌ 网络错误";
				submitBtn.disabled = false;
				cancelBtn.disabled = true;
			};

			xhr.open("POST", "/upload");
			xhr.send(formData);

			progressContainer.style.display = "block";
			progressBar.style.width = "0%";
			progressText.innerText = "开始上传...";
		});

		cancelBtn.addEventListener("click", function() {
			if (xhr) {
				xhr.abort();
				progressText.innerText = "❌ 上传已取消";
				progressBar.style.width = "0%";
				submitBtn.disabled = false;
				cancelBtn.disabled = true;
			}
		});
	</script>
</body>
</html>`
		t, _ := template.New("page").Parse(tmpl)
		t.Execute(w, fileList)
	})

	// 打印本机可访问地址
	log.Println("✅ FileBeam 文件共享服务启动：")
	log.Println("  本地访问: http://localhost:8888/")
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				log.Printf("  局域网访问: http://%s:8888/", ipnet.IP.String())
			}
		}
	}

	err := http.ListenAndServe("0.0.0.0:8888", nil)
	if err != nil {
		log.Fatal("❌ 启动失败:", err)
	}
}
