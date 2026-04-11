oapi-codegen -generate types       -o internal/api/remote/oapi-gen/types.gen.go  -package oapi typespec/out/spec-remote.yaml
oapi-codegen -generate echo-server -o internal/api/remote/oapi-gen/server.gen.go -package oapi typespec/out/spec-remote.yaml
oapi-codegen -generate spec        -o internal/api/remote/oapi-gen/spec.gen.go   -package oapi typespec/out/spec-remote.yaml