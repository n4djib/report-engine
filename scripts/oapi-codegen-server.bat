oapi-codegen -generate types       -o internal/api/server/oapi-gen/types.gen.go  -package oapi typespec/out/spec-server.yaml
oapi-codegen -generate echo-server -o internal/api/server/oapi-gen/server.gen.go -package oapi typespec/out/spec-server.yaml
oapi-codegen -generate spec        -o internal/api/server/oapi-gen/spec.gen.go   -package oapi typespec/out/spec-server.yaml