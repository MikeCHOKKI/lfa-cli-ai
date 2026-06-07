@echo off
setlocal enabledelayedexpansion

set REPO=MikeCHOKKI/lfa-cli-ai
set BINARY=lfa
set VERSION=%1
if "%VERSION%"=="" set VERSION=latest

REM Detect architecture
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" set ARCH=amd64
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" set ARCH=arm64
if "%ARCH%"=="" set ARCH=amd64

REM Get latest version
if "%VERSION%"=="latest" (
  for /f "tokens=*" %%a in ('curl -fsSL "https://api.github.com/repos/%REPO%/releases/latest" ^| findstr "tag_name"') do (
    set LINE=%%a
    set VERSION=!LINE:*"v"=!
    set VERSION=!VERSION:"=!
  )
)

set BINARY_NAME=lfa-windows-%ARCH%.exe
set DOWNLOAD_URL=https://github.com/%REPO%/releases/download/v%VERSION%/%BINARY_NAME%
set INSTALL_DIR=%USERPROFILE%\.local\bin

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

echo Downloading LFA CLI v%VERSION% for windows/%ARCH%...
curl -fsSL "%DOWNLOAD_URL%" -o "%INSTALL_DIR%\%BINARY%.exe"

REM Add to PATH if not already (use setx with safe length)
echo %PATH% | findstr /C:"%INSTALL_DIR%" >nul 2>&1
if errorlevel 1 (
  echo Adding %INSTALL_DIR% to PATH...
  setx PATH "%PATH%;%INSTALL_DIR%" >nul 2>&1
  if errorlevel 1 (
    REM Fallback: set user-level PATH only to avoid 1024-char truncation
    for /f "tokens=2*" %%a in ('reg query HKCU\Environment /v PATH 2^>nul') do if not "%%b"=="" (
      setx PATH "%%b;%INSTALL_DIR%" >nul 2>&1
    )
  )
)

echo.
echo LFA CLI v%VERSION% installed successfully!
echo Run 'lfa setup' to configure your environment.
endlocal
