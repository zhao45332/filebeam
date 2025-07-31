@echo off
echo ========================================
echo FileBeam 文件共享服务启动脚本
echo ========================================
echo.

REM 检查是否存在可执行文件
if not exist "filebeam.exe" (
    echo 正在编译项目...
    go build -o filebeam.exe main.go
    if errorlevel 1 (
        echo 编译失败！请检查Go环境是否正确安装。
        pause
        exit /b 1
    )
)

echo 启动 FileBeam 服务...
echo.
echo 提示：
echo - 默认访问地址：http://localhost:8888/
echo - 默认上传密码：123456
echo - 按 Ctrl+C 停止服务
echo.
echo ========================================

filebeam.exe 