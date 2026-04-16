$envFile = "cmd\server\env\.env"
$targetFile = "typespec\out\spec-server.yaml"

# load env
Get-Content $envFile | Where-Object {$_ -match '^[^#].+=.+'} | ForEach-Object {
    $k,$v = $_ -split '=',2
    $v = $v.Trim().Trim('"')
    Set-Item -Path "Env:$k" -Value $v
}

$search = "- url: http://localhost:8080"

$fullUrl = $env:APP_URL + ":" + $env:APP_PORT
$replace = "- url: " + $fullUrl

(Get-Content $targetFile) -replace [regex]::Escape($search), $replace |
Set-Content $targetFile