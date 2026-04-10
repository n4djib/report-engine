== Initialize GO ==
> go mod init github.com/n4djib/report-engine
    initialize mod file
> go mod tidy
    tidy up mod


== GIT
> git config --global user.name "Your Name"
> git config --global user.email "your.email@example.com"
> git init


== Folder Structure ==
/cmd      : contains the apps
/pkg      : is for code you are OK with other projects importing.
/internal : is for code that must stay private inside your repo
/rest     : is for rest api tests
/tmp      : air build output path
/web      : contains frontend code


== Air ==
> go install github.com/air-verse/air@latest
> air init
    this generates .air.toml


== Run using Makefile.cmd
> Makefile.cmd            display help
> Makefile.cmd dc         run Services in docker-compose
> Makefile.cmd sd         run Server with AIR
> Makefile.cmd rd         run Remote with AIR


== nk (nkey key pair generator)
> go install github.com/nats-io/nkeys/nk@latest
vnk -gen user -pubout

we can generate the key pair using "github.com/nats-io/nkeys"


== pnpm
> npm install -g pnpm@latest-10
> pnpm setup


== install frontend
> pnpx create-tsrouter-app@latest


== run frontend
> pnpm dev
> pnpm build


== install echo
> go get github.com/labstack/echo/v4


== install typespec
> pnpm add -g @typespec/compiler
> mkdir typespec
> cd typespec
typespec> tsp init
    initialize typespec inside typespec folder
typespec> cd ..

change:
openapi-versions:
      - 3.0.0

> tsp compile typespec
> tsp compile typespec --watch


== installe echoSwagger
> go get github.com/swaggo/echo-swagger
> go get gopkg.in/yaml.v2


== install oapi-codegen
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest


== rfontend client code generator

pnpm add -D @openapitools/openapi-generator-cli

pnpx @openapitools/openapi-generator-cli generate \
  -i typespec/spec.yaml \
  -g typescript-axios \
  -o src/api-generated

how about  openapi-zod-client



