@REM @REM this generates typespec/spec.yaml
cmd /c tsp compile typespec


@REM this generates the GO API code in 
oapi-codegen -generate types       -o internal/api/server/oapi-gen/types.gen.go  -package oapi typespec/spec.yaml
oapi-codegen -generate echo-server -o internal/api/server/oapi-gen/server.gen.go -package oapi typespec/spec.yaml
oapi-codegen -generate spec        -o internal/api/server/oapi-gen/spec.gen.go   -package oapi typespec/spec.yaml

@REM change the generated code to echo v5
powershell -Command "(Get-Content internal/api/server/oapi-gen/server.gen.go -Raw) -replace 'echo/v4', 'echo/v5' | Set-Content internal/api/server/oapi-gen/server.gen.go"
powershell -Command "(Get-Content internal/api/server/oapi-gen/server.gen.go -Raw) -replace 'ctx echo.Context', 'ctx *echo.Context' | Set-Content internal/api/server/oapi-gen/server.gen.go"


@REM cmd /c npm -C web run prebuild
