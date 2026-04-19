cmd /c pnpm -C web --filter central build

@REM go:embed all:dist
@REM go:embed dist/index.html
@REM go:embed routeTree.gen.ts

if exist internal\embed\central\dist rmdir /s /q internal\embed\central\dist
if exist internal\embed\central\routeTree.gen.ts rmdir /s /q internal\embed\central\routeTree.gen.ts

xcopy web\apps\central\dist internal\embed\central\dist /E /I /Y
copy web\apps\central\src\routeTree.gen.ts internal\embed\central\routeTree.gen.ts