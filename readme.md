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

# Project Setup

## Clone repo

> git clone https://github.com/n4djib/report-engine.git

## Download GO packages

> go mod tidy

## install pnpm globally

> npm install -g pnpm@latest-10

## install web packages

> cd web
> pnpm install

## install typespec

> pnpm add -g @typespec/compiler
> cd typespec
> pnpm install

# Run Project

> Makefile.cmd spec-gen
> Makefile.cmd dcd
