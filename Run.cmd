@echo off
chcp 65001&cls

cd /d "%~dp0"

set DEBUG=0
set PLATFORM=pc
set MODE=run
set CLEAN=0

:parse_args
if "%~1"=="" goto done_args
if "%~1"=="-d" (
    set DEBUG=1
    set MODE=build
)
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
if exist "src-tauri\target" rd /s /q "src-tauri\target"
if exist "src-vue\dist" rd /s /q "src-vue\dist"
echo [INFO] Clean complete

:after_clean
where cargo >NUL 2>&1
if errorlevel 1 (
    echo [FATAL] Rust not installed
    pause
    exit /b 1
)

if "%PLATFORM%"=="mobile" goto mobile
goto pc

:pc
if not exist "src-vue\node_modules" (
    echo [INFO] Installing frontend dependencies...
    cd src-vue
    call npm install
    cd ..
)

if "%MODE%"=="build" goto pc_build
goto pc_run

:pc_run
if "%DEBUG%"=="1" (
    echo [INFO] Starting in debug mode...
    cargo tauri dev -- -- -d
) else (
    echo [INFO] Starting application...
    start "" /b cargo tauri dev
)
goto end

:pc_build
if not exist "runtime\build" mkdir "runtime\build"
if "%DEBUG%"=="1" (
    echo [INFO] Building debug version...
    cargo tauri build --debug
    if errorlevel 1 (
        echo [FATAL] Debug build failed
        pause
        goto end
    )
    echo [INFO] Copying debug build to runtime\build...
    copy /y "src-tauri\target\debug\tavern.exe" "runtime\build\Tavern-dev.exe"
    echo [INFO] Starting debug version...
    start "" "runtime\build\Tavern-dev.exe" -d
    goto end
)
echo [INFO] Building release version...
cargo tauri build
if errorlevel 1 (
    echo [FATAL] Build failed
    pause
    goto end
)
echo [INFO] Copying release build to runtime\build...
copy /y "src-tauri\target\release\tavern.exe" "runtime\build\Tavern.exe"
echo [INFO] Build complete
pause
goto end

:mobile
if not exist "src-vue\node_modules" (
    echo [INFO] Installing frontend dependencies...
    cd src-vue
    call npm install
    cd ..
)

if not exist "src-tauri\gen\android" (
    echo [INFO] Initializing Android project...
    cargo tauri android init
)

if "%MODE%"=="build" goto mobile_build
goto mobile_run

:mobile_run
if "%DEBUG%"=="1" (
    echo [INFO] Starting Android in debug mode...
    cargo tauri android dev -- -- -d
) else (
    echo [INFO] Starting Android...
    start "" /b cargo tauri android dev
)
goto end

:mobile_build
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
