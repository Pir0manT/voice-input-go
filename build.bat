@echo off
setlocal

set APP_NAME=voice-input-go
set OUTPUT=%APP_NAME%.exe
set LDFLAGS=-s -w -extldflags '-static'
set CGO_LDFLAGS=-static -lportaudio -lole32 -lwinmm -lksuser -lsetupapi -luuid

echo ======================================
echo  Building %OUTPUT% (static linking)
echo ======================================
echo.
echo GOROOT:     %GOROOT%
echo CGO_LDFLAGS: %CGO_LDFLAGS%
echo LDFLAGS:     %LDFLAGS%
echo.
echo [1/2] Compiling... (this may take a minute)

go build -v -ldflags="%LDFLAGS%" -o %OUTPUT% ./cmd/%APP_NAME%/ 2>&1
if errorlevel 1 (
    echo.
    echo *** BUILD FAILED ***
    exit /b 1
)
echo [1/2] Build complete.
echo.

:: Find UPX: PATH, then WinGet links
set UPX_CMD=
where upx >nul 2>nul
if %errorlevel% equ 0 (
    set UPX_CMD=upx
) else if exist "%LOCALAPPDATA%\Microsoft\WinGet\Links\upx.exe" (
    set "UPX_CMD=%LOCALAPPDATA%\Microsoft\WinGet\Links\upx.exe"
)

if defined UPX_CMD (
    echo [2/2] Compressing with UPX...
    "%UPX_CMD%" --best %OUTPUT%
) else (
    echo [2/2] UPX not found, skipping compression
    echo     Install: winget install upx.upx
)

echo.
echo ======================================
for %%A in (%OUTPUT%) do echo  Done: %OUTPUT%  %%~zA bytes
echo ======================================
