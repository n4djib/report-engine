oapi-codegen -generate types       -o internal/api/central/oapi-gen/types.gen.go  -package oapi typespec/out/spec-central.yaml
oapi-codegen -generate echo-server -o internal/api/central/oapi-gen/server.gen.go -package oapi typespec/out/spec-central.yaml
oapi-codegen -generate spec        -o internal/api/central/oapi-gen/spec.gen.go   -package oapi typespec/out/spec-central.yaml