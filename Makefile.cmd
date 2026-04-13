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
    @REM pnpm --filter ./web/apps/server dev
    pnpm --filter server dev
    exit /b %errorlevel%

:rf
    echo === Remote Frontend  dev ===
    echo Running Remote Frontend... 
    :: Add your build commands here
    @REM pnpm --filter ./web/apps/remote dev
    pnpm --filter remote dev
    exit /b %errorlevel%

:spec-gen:
    echo === Generating API Spec and Code ===
    echo Generating API Spec and Code... 
    :: Add your build commands here
        @REM this generates typespec/spec.yaml
        cmd /c scripts\spec-gen-server.bat
        cmd /c scripts\spec-gen-remote.bat
        @REM this generates the GO API code in 
        cmd /c scripts\oapi-codegen-server.bat
        cmd /c scripts\oapi-codegen-remote.bat
        @REM Generate the client API code in
        cmd /c scripts\update-client-api-server.bat
        cmd /c scripts\update-client-api-remote.bat
    exit /b %errorlevel%

:help
    echo Usage: Makefile [command]
    echo.
    echo Commands:
    echo   spec-gen - Generate API Spec and Code
    echo   dc       - Spinning Up docker services
    echo   ---      - - - - - - - - - - - - - - - -
    echo   nk       - Generate Key Pair (nk Util)
    echo   nkeys    - Generate Key Pair using nkeys (GO app)
    echo   ---      - - - - - - - - - - - - - - - -
    echo   ---      - depricated commands (replaced by docker compose and air)
    echo   sd       - launch Serve in Air
    echo   rd       - launch Remote in Air
    echo   sf       - launch server frontend dev server
    echo   rf       - launch remote frontend dev server
    exit /b 0