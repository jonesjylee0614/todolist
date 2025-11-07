@echo off
echo 开始编译 Go 项目到 CentOS 7.8...

REM 设置编译环境
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

REM 编译信息
for /f "tokens=*" %%i in ('git describe --tags --always 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" set VERSION=unknown
set BUILD_TIME=%date:~0,4%-%date:~5,2%-%date:~8,2%_%time:~0,2%:%time:~3,2%:%time:~6,2%
set BUILD_TIME=%BUILD_TIME: =0%

REM 编译
echo 编译目标: %GOOS%/%GOARCH%
echo 版本信息: %VERSION%
echo 构建时间: %BUILD_TIME%

go build -ldflags="-w -s -X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME%" -o todolist-server main.go

if %ERRORLEVEL% EQU 0 (
    echo 编译成功!
    echo 输出文件: todolist-server
    for %%F in (todolist-server) do echo 文件大小: %%~zF 字节
    
    REM 验证
    echo 部署步骤:
    echo 1. 传输到服务器: scp todolist-server user@server:/path/
    echo 2. 在服务器上运行: chmod +x todolist-server && ./todolist-server
) else (
    echo 编译失败!
    exit /b 1
)
