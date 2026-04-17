cmd /c tsp compile typespec/main-central.tsp --config typespec/tspconfig-central.yaml

@echo off
setlocal

powershell -NoProfile -ExecutionPolicy Bypass -File scripts\update-spec-url-central.ps1

endlocal