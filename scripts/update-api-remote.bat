cmd /c pnpx openapi-zod-client typespec\out\spec-remote.yaml --output web/remote/src/api --group-strategy=tag-file --strict-objects


@echo off
set "API_PATH=web/remote/src/api"

echo Stripping .passthrough() from files in %API_PATH%...

powershell -Command "Get-ChildItem '%API_PATH%\*.ts' | ForEach-Object { (Get-Content $_.FullName) -replace '\.passthrough\(\)', '' | Set-Content $_.FullName }"

echo Success!
