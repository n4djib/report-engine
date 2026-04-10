cmd /c tsp compile typespec/main-remote.tsp --config typespec/tspconfig-remote.yaml

@echo off
setlocal

powershell -NoProfile -ExecutionPolicy Bypass -File scripts\update-spec-url-remote.ps1

endlocal