# Folder Structure

- **/cmd :** contains the apps entries and env
- **/docker :** docker config
- **/internal :** is for code that must stay private inside your repo
- **/pkg :** is for code you are OK with other projects importing
- **/rest :** is for rest api tests
- **/scripts :** contains scripts to generate code
- **/typespec :** contains typespec files and output
- **/web :** contains frontend code

**Note:** we have two monrepos inside each other

- GO monorepo in the root (/) of the project
- pnpm monorepo (workspace) in /web

# Run using Makefile.cmd

> \> **Makefile.cmd**  
> display help  
> \> **Makefile.cmd spec-gen**  
>  Generate API Spec and Code  
> \> **Makefile.cmd dcd**  
>  run Services in docker-compose (developement)  
> \> **Makefile.cmd dcp**  
>  run Services in docker-compose (Production)  
> \> **Makefile.cmd flint**  
>  Lint frontend code

== Initialize GO

> go mod init github.com/n4djib/report-engine

    initialize mod file

> go mod tidy

    tidy up mod

== GIT

> git config --global user.name "Your Name"
> git config --global user.email "your.email@example.com"
> git init

== Air ==

> go install github.com/air-verse/air@latest
> air init

    this generates .air.toml

== nk (nkey key pair generator)

> go install github.com/nats-io/nkeys/nk@latest
> vnk -gen user -pubout

we can generate the key pair using "github.com/nats-io/nkeys"

== pnpm

> npm install -g pnpm@latest-10
> pnpm setup

== install frontend

> pnpx create-tsrouter-app@latest

web> pnpm add -D -w eslint @eslint/js typescript-eslint
add eslint

== run frontend

> pnpm dev
> pnpm build

== install echo

> go get github.com/labstack/echo/v4
> go get github.com/labstack/echo/v4/middleware

== install typespec

> pnpm add -g @typespec/compiler
> mkdir typespec
> cd typespec
> typespec> tsp init

    initialize typespec inside typespec folder

typespec> cd ..

change:
openapi-versions: - 3.0.0

> tsp compile typespec
> tsp compile typespec --watch

== installe echoSwagger

> go get github.com/swaggo/echo-swagger
> go get gopkg.in/yaml.v2

== install oapi-codegen
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

== fontend client code generator with openapi-zod-client
web\server> pnpm add zod
web\server> pnpm add -D openapi-zod-client
web\server> pnpm add @zodios/core

> pnpx openapi-zod-client typespec\out\spec-server.yaml --output web/server/src/api --group-strategy=tag-file --strict-objects

== add tanstack/react-query
pnpm add @tanstack/react-query

== switicing web apps to monorepo
web> pnpm install
web> pnpm add @types/react

pnpm --filter server dev
pnpm --filter remote dev
