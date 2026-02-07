@echo off
echo ====================================
echo 智慧养老系统 - 依赖安装脚本
echo ====================================
echo.

echo [1/2] 安装 Go 依赖...
cd /d "%~dp0"
go mod download
if %errorlevel% neq 0 (
    echo 错误: Go 依赖安装失败
    pause
    exit /b 1
)
echo Go 依赖安装完成
echo.

echo [2/2] 安装前端依赖...
cd frontend
call npm install
if %errorlevel% neq 0 (
    echo 错误: 前端依赖安装失败
    pause
    exit /b 1
)
cd ..
echo.

echo ====================================
echo 所有依赖安装完成！
echo 运行 'build.bat' 启动开发服务器
echo ====================================
pause
