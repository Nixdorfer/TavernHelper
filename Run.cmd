@echo off
chcp 65001 >NUL 2>&1

cd /d "%~dp0"

set VER=go
set DEBUG=0
set PLATFORM=pc
set MODE=run
set CLEAN=0

:parse_args
if "%~1"=="" goto done_args
if "%~1"=="-d" set DEBUG=1
if "%~1"=="-r" set VER=rust
if "%~1"=="-m" set PLATFORM=mobile
if "%~1"=="-b" set MODE=build
if "%~1"=="-c" set CLEAN=1
shift
goto parse_args
:done_args

if "%CLEAN%"=="1" goto do_clean
goto after_clean

:do_clean
echo [INFO] Cleaning build artifacts...
if exist "src-wails\dist" rd /s /q "src-wails\dist"
if exist "src-tauri\target" rd /s /q "src-tauri\target"
if exist "src-wails\build\bin" (
    for /d %%D in ("src-wails\build\bin\*") do (
        if /i not "%%~nxD"=="creations" if /i not "%%~nxD"=="images" if /i not "%%~nxD"=="text" rd /s /q "%%D"
    )
    for %%F in ("src-wails\build\bin\*.*") do del /q "%%F"
)
echo [INFO] Clean complete

:after_clean
if "%VER%"=="rust" goto rust_mode
goto go_mode

:go_mode
if not exist "src-vue\node_modules" (
    echo [INFO] Installing frontend dependencies...
    cd src-vue
    call npm install
    cd ..
)

where wails >NUL 2>&1
if errorlevel 1 goto install_wails
goto wails_ready

:install_wails
echo [FATAL] Wails not installed, installing...
go install github.com/wailsapp/wails/v2/cmd/wails@latest
if errorlevel 1 (
    echo [FATAL] Wails installation failed
    pause
    exit /b 1
)
echo [INFO] Wails installed

:wails_ready
cd src-wails
if "%MODE%"=="build" goto go_build
goto go_run

:go_run
if "%DEBUG%"=="1" (
    echo [INFO] Starting in debug mode...
    set OPEN_DEVTOOLS=1
    wails dev
) else (
    echo [INFO] Starting application...
    start "" /b wails dev -skipbindings
)
goto end

:go_build
echo [INFO] Building application...
wails build -clean -platform windows/amd64
if errorlevel 1 (
    echo [FATAL] Build failed
) else (
    echo [INFO] Build complete
)
pause
goto end

:rust_mode
where cargo >NUL 2>&1
if errorlevel 1 (
    echo [FATAL] Rust not installed
    pause
    exit /b 1
)

if "%PLATFORM%"=="mobile" goto rust_mobile
goto rust_pc

:rust_pc
if not exist "src-vue\node_modules" (
    echo [INFO] Installing PC frontend dependencies...
    cd src-vue
    call npm install
    cd ..
)

if "%MODE%"=="build" goto rust_pc_build
goto rust_pc_run

:rust_pc_run
if "%DEBUG%"=="1" (
    echo [INFO] Starting Rust PC in debug mode...
    cargo tauri dev
) else (
    echo [INFO] Starting Rust PC...
    start "" /b cargo tauri dev
)
goto end

:rust_pc_build
echo [INFO] Building Rust PC...
cargo tauri build
if errorlevel 1 (
    echo [FATAL] Build failed
) else (
    echo [INFO] Build complete
)
pause
goto end

:rust_mobile
if not exist "src-vue\node_modules" (
    echo [INFO] Installing mobile frontend dependencies...
    cd src-vue
    call npm install
    cd ..
)

if not exist "src-tauri\gen\android" (
    echo [INFO] Initializing Android project...
    cargo tauri android init
)

if "%MODE%"=="build" goto rust_mobile_build
goto rust_mobile_run

:rust_mobile_run
if "%DEBUG%"=="1" (
    echo [INFO] Starting Android in debug mode...
    cargo tauri android dev
) else (
    echo [INFO] Starting Android...
    start "" /b cargo tauri android dev
)
goto end

:rust_mobile_build
echo [INFO] Building Android APK...
cargo tauri android build
if errorlevel 1 (
    echo [FATAL] Build failed
) else (
    echo [INFO] Build complete
)
pause
goto end

:end
