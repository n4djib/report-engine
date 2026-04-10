cmd /c tsp compile typespec/main-server.tsp --config typespec/tspconfig-server.yaml

@echo off
setlocal

powershell -NoProfile -ExecutionPolicy Bypass -File scripts\update-spec-url-server.ps1

endlocal