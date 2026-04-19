package frontendembed

import (
	"embed"
	"log"

	"github.com/labstack/echo/v4"
	pathextractor "github.com/n4djib/report-engine/internal/embed"
)

var (
	//go:embed all:dist
	dist embed.FS
	//go:embed dist/index.html
	indexHTML embed.FS
	//go:embed routeTree.gen.ts
	routeTreeFS embed.FS

	distDirFS     = echo.MustSubFS(dist, "dist")
	distIndexHtml = echo.MustSubFS(indexHTML, "dist")
)

func RegisterHandlers(e *echo.Echo) {
	file := "routeTree.gen.ts"

	routes, err := pathextractor.ExtractAndTransformRoutes(routeTreeFS, file)
	if err != nil {
		log.Fatal("couldn't extract routes from ", file, "\n", err)
	}

	for _, r := range routes {
		e.FileFS(r, "index.html", distIndexHtml)
	}

	e.StaticFS("/", distDirFS)
}
