@echo off
setlocal enabledelayedexpansion

set "PROJECT_ROOT=%~dp0"
set "BACKEND_TITLE=TodoList Backend"
set "FRONTEND_TITLE=TodoList Frontend"

call :validate_dirs || exit /b 1

set "ACTION=%~1"
if "%ACTION%"=="" set "ACTION=a"

if /I "%ACTION%"=="a" goto restart_all
if /I "%ACTION%"=="b" goto restart_backend
if /I "%ACTION%"=="f" goto restart_frontend

echo 用法:
echo   start-dev.bat        ^(或 start-dev.bat a^) 重启前后端
echo   start-dev.bat b      只重启后端
echo   start-dev.bat f      只重启前端
exit /b 1

:restart_backend
echo Restarting backend server...
call :stop_backend
call :start_backend
goto end_script

:restart_frontend
echo Restarting frontend server...
call :stop_frontend
call :start_frontend
goto end_script

:restart_all
echo Restarting backend and frontend servers...
call :stop_backend
call :stop_frontend
call :start_backend
call :start_frontend
goto end_script

:start_backend
start "%BACKEND_TITLE%" cmd /k "cd /d ""%PROJECT_ROOT%backend"" && go run ."
echo Backend window launched.
exit /b 0

:start_frontend
start "%FRONTEND_TITLE%" cmd /k "cd /d ""%PROJECT_ROOT%frontend"" && npm run dev"
echo Frontend window launched.
exit /b 0

:stop_backend
call :kill_window "%BACKEND_TITLE%" && (
  echo Backend window stopped.
) || (
  echo No backend window found to stop.
)
exit /b 0

:stop_frontend
call :kill_window "%FRONTEND_TITLE%" && (
  echo Frontend window stopped.
) || (
  echo No frontend window found to stop.
)
exit /b 0

:kill_window
set "WINDOW_TITLE=%~1"
set "TERMINATED=0"
for /f "skip=1 tokens=2 delims=," %%P in ('tasklist /v /fo csv /fi "WINDOWTITLE eq %WINDOW_TITLE%" ^| findstr /I /C:"%WINDOW_TITLE%" 2^>nul') do (
  set "PID=%%~P"
  if not "!PID!"=="" (
    taskkill /pid !PID! /t /f >nul 2>&1
    set "TERMINATED=1"
  )
)
if "!TERMINATED!"=="0" (
  exit /b 1
) else (
  exit /b 0
)

:validate_dirs
if not exist "%PROJECT_ROOT%backend" (
  echo [ERROR] backend directory not found next to this script.
  exit /b 1
)
if not exist "%PROJECT_ROOT%frontend" (
  echo [ERROR] frontend directory not found next to this script.
  exit /b 1
)
exit /b 0

:end_script
echo Done.
endlocal
exit /b 0
