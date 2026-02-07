@echo off
echo ====================================
echo 智慧养老系统 - 开发环境启动脚本
echo ====================================
echo.

echo [1/3] 检查 Go 环境...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到 Go，请先安装 Go 1.21+
    pause
    exit /b 1
)
echo Go 环境检查通过
echo.

echo [2/3] 安装前端依赖...
cd frontend
if not exist "node_modules" (
    echo 首次运行，正在安装前端依赖...
    call npm install
    if %errorlevel% neq 0 (
        echo 错误: 前端依赖安装失败
        pause
        exit /b 1
    )
) else (
    echo 前端依赖已存在
)
cd ..
echo.

echo [3/3] 启动 Wails 开发服务器...
echo 应用将在浏览器中自动打开
echo.
wails dev
pause
