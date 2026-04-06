@echo off
setlocal enabledelayedexpansion

:: Check if the user passed an argument; if not, show help
if "%~1"=="" goto help

:: Jump to the label provided in the first argument
findstr /i /c:":%~1" "%~f0" >nul 2>&1
if %errorlevel% neq 0 (
    echo [Error] Unknown command: %~1
    goto help
)
goto %~1

:dc
    echo === docker-compose ===
    echo Spinning Up... 
    :: Add your build commands here
    docker compose -f docker/docker-compose.yaml up
    exit /b %errorlevel%

:sd
    echo === Server dev ===
    echo Watch Server in Air... 
    :: Add your build commands here
    air -c .\.air-server.toml
    @REM start air -c .\.air-server.toml
    exit /b %errorlevel%

:rd
    echo === Remote dev ===
    echo Watch Remote in Air... 
    :: Add your build commands here
    air -c .\.air-remote.toml
    exit /b %errorlevel%

:nk
    echo === Generate with NK ===
    echo Generating Key Pair... 
    :: Add your build commands here
    nk -gen user -pubout
    exit /b %errorlevel%

:nkeys
    echo === Generate Key Pair ===
    echo Generating Key Pair... 
    :: Add your build commands here
    go run .\cmd\nkeys-generator\
    exit /b %errorlevel%

:sf
    echo === Server Frontend  dev ===
    echo Running Server Frontend... 
    :: Add your build commands here
    pnpm --filter ./web/server dev
    exit /b %errorlevel%

:rf
    echo === Remote Frontend  dev ===
    echo Running Remote Frontend... 
    :: Add your build commands here
    pnpm --filter ./web/remote dev
    exit /b %errorlevel%

:help
    echo Usage: Makefile [command]
    echo.
    echo Commands:
    echo   dc - Spinning Up docker services
    echo   sd - launch Serve in Air
    echo   rd - launch Remote in Air
    echo   sf - launch server frontend dev server
    echo   rf - launch remote frontend dev server
    echo   nk - Generate Key Pair
    echo   nkeys - Generate Key Pair using nkeys
    exit /b 0