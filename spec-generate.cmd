@REM this generates typespec/spec.yaml
@REM cmd /c scripts\spec-gen-server.bat
@REM cmd /c scripts\spec-gen-remote.bat


@REM @REM this generates the GO API code in 
@REM oapi-codegen -generate types       -o internal/api/server/oapi-gen/types.gen.go  -package oapi typespec/out/spec-server.yaml
@REM oapi-codegen -generate echo-server -o internal/api/server/oapi-gen/server.gen.go -package oapi typespec/out/spec-server.yaml
@REM oapi-codegen -generate spec        -o internal/api/server/oapi-gen/spec.gen.go   -package oapi typespec/out/spec-server.yaml

@REM oapi-codegen -generate types       -o internal/api/remote/oapi-gen/types.gen.go  -package oapi typespec/out/spec-remote.yaml
@REM oapi-codegen -generate echo-server -o internal/api/remote/oapi-gen/server.gen.go -package oapi typespec/out/spec-remote.yaml
@REM oapi-codegen -generate spec        -o internal/api/remote/oapi-gen/spec.gen.go   -package oapi typespec/out/spec-remote.yaml


cmd /c scripts\update-api-server.bat
cmd /c scripts\update-api-remote.bat
