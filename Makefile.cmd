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

:dcd
    echo === docker-compose dev ===
    echo Spinning Up... 
    docker compose -f docker/docker-compose.dev-build.yaml build
    docker compose -f docker/docker-compose.dev.yaml up
    exit /b %errorlevel%
:dcp
    echo === docker-compose prod ===
    echo Spinning Up...
    docker compose -f docker/docker-compose.prod-build.yaml build
    docker compose -f docker/docker-compose.prod.yaml up
    exit /b %errorlevel%

:cd
    echo === Central dev ===
    @REM echo Watch Central in Air... 
    :: Add your build commands here
    @REM air -c .\.air-central.toml
    @REM start air -c .\.air-central.toml
    go run .\cmd\central\
    exit /b %errorlevel%

:rd
    echo === Remote dev ===
    @REM echo Watch Remote in Air... 
    :: Add your build commands here
    @REM air -c .\.air-remote.toml
    go run .\cmd\remote\
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

:cf
    echo === Central Frontend  dev ===
    echo Running Central Frontend... 
    :: Add your build commands here
    pnpm -C web --filter central dev
    exit /b %errorlevel%

:rf
    echo === Remote Frontend  dev ===
    echo Running Remote Frontend... 
    :: Add your build commands here
    pnpm -C web --filter remote dev
    exit /b %errorlevel%
    
:flint
    echo === Lint Frontends ===
    echo linting Frontends... 
    :: Add your build commands here
    pnpm -C web lint
    exit /b %errorlevel%

:spec-gen:
    echo === Generating API Spec and Code ===
    echo Generating API Spec and Code... 
    :: Add your build commands here
        @REM this generates typespec/spec.yaml
        cmd /c scripts\spec-gen-central.bat
        cmd /c scripts\spec-gen-remote.bat
        @REM this generates the GO API code in 
        cmd /c scripts\oapi-codegen-central.bat
        cmd /c scripts\oapi-codegen-remote.bat
        @REM Generate the client API code in
        cmd /c scripts\update-client-api-central.bat
        cmd /c scripts\update-client-api-remote.bat
    exit /b %errorlevel%

:help
    echo Usage: Makefile [command]
    echo.
    echo Commands:
    echo   spec-gen - Generate API Spec and Code
    echo   dcd      - Spinning Up docker services (development)
    echo   dcp      - Spinning Up docker services (production)
    echo   flint    - Lint frontend code 
    echo   ---      - - - - - - - - - - - - - - - -
    echo   nk       - Generate Key Pair (nk Util)
    echo   nkeys    - Generate Key Pair using nkeys (GO app)
    echo   ---      - - - - - - - - - - - - - - - -
    echo   ---      - depricated commands (replaced by docker compose and air) (# TODO to be removed)
    echo   cd       - launch Central backend
    echo   rd       - launch Remote backend
    echo   cf       - launch Central frontend dev central
    echo   rf       - launch Remote frontend dev central
    exit /b 0