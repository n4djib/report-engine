@REM @REM @REM this generates typespec/spec.yaml
@REM cmd /c tsp compile typespec/main-server.tsp --config typespec/tspconfig-server.yaml
@REM cmd /c tsp compile typespec/main-remote.tsp --config typespec/tspconfig-remote.yaml


@REM this generates the GO API code in 
oapi-codegen -generate types       -o internal/api/server/oapi-gen/types.gen.go  -package oapi typespec/out/spec-server.yaml
oapi-codegen -generate echo-server -o internal/api/server/oapi-gen/server.gen.go -package oapi typespec/out/spec-server.yaml
oapi-codegen -generate spec        -o internal/api/server/oapi-gen/spec.gen.go   -package oapi typespec/out/spec-server.yaml

oapi-codegen -generate types       -o internal/api/remote/oapi-gen/types.gen.go  -package oapi typespec/out/spec-remote.yaml
oapi-codegen -generate echo-server -o internal/api/remote/oapi-gen/server.gen.go -package oapi typespec/out/spec-remote.yaml
oapi-codegen -generate spec        -o internal/api/remote/oapi-gen/spec.gen.go   -package oapi typespec/out/spec-remote.yaml


@REM cmd /c npm -C web run prebuild
