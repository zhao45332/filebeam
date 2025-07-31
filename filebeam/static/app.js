// FileBeam 前端应用
class FileBeamApp {
    constructor() {
        this.currentXHR = null;
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadFileList();
    }

    bindEvents() {
        // 文件上传表单
        const uploadForm = document.getElementById('uploadForm');
        uploadForm.addEventListener('submit', (e) => this.handleUpload(e));

        // 刷新按钮
        const refreshBtn = document.getElementById('refreshBtn');
        refreshBtn.addEventListener('click', () => this.loadFileList());

        // 取消上传按钮
        const cancelBtn = document.getElementById('cancelBtn');
        cancelBtn.addEventListener('click', () => this.cancelUpload());

        // 文件选择显示
        const fileInput = document.getElementById('fileInput');
        fileInput.addEventListener('change', (e) => this.updateFileLabel(e));

        // 密码输入框回车提交
        const passwordInput = document.getElementById('passwordInput');
        passwordInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                uploadForm.dispatchEvent(new Event('submit'));
            }
        });
    }

    // 加载文件列表
    async loadFileList() {
        const fileList = document.getElementById('fileList');
        fileList.innerHTML = '<div class="loading">加载中...</div>';

        try {
            const response = await fetch('/api/files');
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }

            const data = await response.json();
            this.renderFileList(data);
        } catch (error) {
            console.error('加载文件列表失败:', error);
            fileList.innerHTML = '<div class="empty-state">❌ 加载失败，请刷新重试</div>';
        }
    }

    // 渲染文件列表
    renderFileList(data) {
        const fileList = document.getElementById('fileList');
        
        if (!data.files || data.files.length === 0) {
            fileList.innerHTML = '<div class="empty-state">📁 暂无共享文件</div>';
            return;
        }

        const fileItems = data.files.map(file => `
            <div class="file-item">
                <div class="file-info">
                    <div class="file-name">${this.escapeHtml(file.name)}</div>
                    <div class="file-meta">
                        📏 ${file.size_human} | 📅 ${file.modified_str}
                    </div>
                </div>
                <div class="file-actions">
                    <a href="${file.download_url}" class="download-link" download>
                        ⬇ 下载
                    </a>
                    <a href="${file.info_url}" class="info-link" target="_blank">
                        ℹ 详情
                    </a>
                </div>
            </div>
        `).join('');

        fileList.innerHTML = fileItems;
    }

    // 处理文件上传
    handleUpload(event) {
        event.preventDefault();

        const form = event.target;
        const formData = new FormData(form);
        const fileInput = document.getElementById('fileInput');
        const passwordInput = document.getElementById('passwordInput');

        // 验证文件是否选择
        if (!fileInput.files[0]) {
            this.showMessage('请选择要上传的文件', 'error');
            return;
        }

        // 验证密码是否输入
        if (!passwordInput.value.trim()) {
            this.showMessage('请输入上传密码', 'error');
            return;
        }

        // 禁用上传按钮，启用取消按钮
        const uploadBtn = document.getElementById('uploadBtn');
        const cancelBtn = document.getElementById('cancelBtn');
        uploadBtn.disabled = true;
        cancelBtn.disabled = false;

        // 显示进度条
        const progressContainer = document.getElementById('progressContainer');
        const progressBar = document.getElementById('progressBar');
        const progressText = document.getElementById('progressText');
        
        progressContainer.style.display = 'block';
        progressBar.style.width = '0%';
        progressText.textContent = '准备上传...';

        // 创建XHR请求
        this.currentXHR = new XMLHttpRequest();

        // 监听上传进度
        this.currentXHR.upload.addEventListener('progress', (e) => {
            if (e.lengthComputable) {
                const percent = (e.loaded / e.total) * 100;
                progressBar.style.width = percent.toFixed(2) + '%';
                progressText.textContent = `上传进度：${percent.toFixed(2)}%`;
            }
        });

        // 监听完成事件
        this.currentXHR.addEventListener('load', () => {
            this.handleUploadComplete();
        });

        // 监听错误事件
        this.currentXHR.addEventListener('error', () => {
            this.handleUploadError('网络错误');
        });

        // 监听中止事件
        this.currentXHR.addEventListener('abort', () => {
            this.handleUploadError('上传已取消');
        });

        // 发送请求
        this.currentXHR.open('POST', '/upload');
        this.currentXHR.send(formData);
    }

    // 处理上传完成
    handleUploadComplete() {
        const uploadBtn = document.getElementById('uploadBtn');
        const cancelBtn = document.getElementById('cancelBtn');
        const progressText = document.getElementById('progressText');

        uploadBtn.disabled = false;
        cancelBtn.disabled = true;

        if (this.currentXHR.status === 200) {
            try {
                const response = JSON.parse(this.currentXHR.responseText);
                if (response.success) {
                    progressText.textContent = '✅ 上传成功！正在刷新文件列表...';
                    this.showMessage('文件上传成功！', 'success');
                    
                    // 重置表单
                    document.getElementById('uploadForm').reset();
                    document.querySelector('.file-input-text').textContent = '选择文件';
                    
                    // 延迟刷新文件列表
                    setTimeout(() => {
                        this.loadFileList();
                        document.getElementById('progressContainer').style.display = 'none';
                    }, 1000);
                } else {
                    this.handleUploadError('上传失败：' + (response.message || '未知错误'));
                }
            } catch (e) {
                this.handleUploadError('响应格式错误');
            }
        } else {
            let errorMessage = '上传失败';
            switch (this.currentXHR.status) {
                case 409:
                    errorMessage = '文件已存在，禁止重复上传';
                    break;
                case 413:
                    errorMessage = '文件过大，超出限制';
                    break;
                case 403:
                    errorMessage = '上传密码错误';
                    break;
                case 400:
                    errorMessage = '请求参数错误';
                    break;
                default:
                    errorMessage = `服务器错误 (${this.currentXHR.status})`;
            }
            this.handleUploadError(errorMessage);
        }
    }

    // 处理上传错误
    handleUploadError(message) {
        const uploadBtn = document.getElementById('uploadBtn');
        const cancelBtn = document.getElementById('cancelBtn');
        const progressText = document.getElementById('progressText');

        uploadBtn.disabled = false;
        cancelBtn.disabled = true;
        progressText.textContent = `❌ ${message}`;
        
        this.showMessage(message, 'error');
        
        // 3秒后隐藏进度条
        setTimeout(() => {
            document.getElementById('progressContainer').style.display = 'none';
        }, 3000);
    }

    // 取消上传
    cancelUpload() {
        if (this.currentXHR) {
            this.currentXHR.abort();
            this.currentXHR = null;
        }
    }

    // 更新文件选择标签
    updateFileLabel(event) {
        const file = event.target.files[0];
        const label = document.querySelector('.file-input-text');
        
        if (file) {
            label.textContent = file.name;
        } else {
            label.textContent = '选择文件';
        }
    }

    // 显示消息提示
    showMessage(message, type = 'info') {
        // 创建消息元素
        const messageEl = document.createElement('div');
        messageEl.className = `message message-${type}`;
        messageEl.textContent = message;
        
        // 添加样式
        messageEl.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 15px 20px;
            border-radius: 8px;
            color: white;
            font-weight: 600;
            z-index: 1000;
            animation: slideIn 0.3s ease;
            max-width: 300px;
        `;
        
        // 根据类型设置背景色
        switch (type) {
            case 'success':
                messageEl.style.background = '#28a745';
                break;
            case 'error':
                messageEl.style.background = '#dc3545';
                break;
            case 'warning':
                messageEl.style.background = '#ffc107';
                messageEl.style.color = '#333';
                break;
            default:
                messageEl.style.background = '#17a2b8';
        }
        
        // 添加到页面
        document.body.appendChild(messageEl);
        
        // 3秒后自动移除
        setTimeout(() => {
            messageEl.style.animation = 'slideOut 0.3s ease';
            setTimeout(() => {
                if (messageEl.parentNode) {
                    messageEl.parentNode.removeChild(messageEl);
                }
            }, 300);
        }, 3000);
    }

    // HTML转义
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// 添加动画样式
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(100%);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);

// 初始化应用
document.addEventListener('DOMContentLoaded', () => {
    new FileBeamApp();
}); 