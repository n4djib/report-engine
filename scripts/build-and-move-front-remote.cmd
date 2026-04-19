cmd /c pnpm -C web --filter remote build

@REM go:embed all:dist
@REM go:embed dist/index.html
@REM go:embed routeTree.gen.ts

if exist internal\embed\remote\dist rmdir /s /q internal\embed\remote\dist
if exist internal\embed\remote\routeTree.gen.ts rmdir /s /q internal\embed\remote\routeTree.gen.ts

xcopy web\apps\remote\dist internal\embed\remote\dist /E /I /Y
copy web\apps\remote\src\routeTree.gen.ts internal\embed\remote\routeTree.gen.ts