package swagger

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v5"

	// echoSwagger "github.com/swaggo/echo-swagger"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func RegisterSwagger(e *echo.Echo, spec *openapi3.T) {
	// from spec.gen.go
	e.GET("/swagger/doc.yaml", registerOpenapiSpec(formatYAML, spec))
	e.GET("/swagger/doc.json", registerOpenapiSpec(formatJSON, spec))

	// Serve Swagger UI (served from /swagger/)
	// e.Static("/swagger", "swaggerui")

	// serving using echoSwagger
	// e.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.URL("/swagger/doc.json")))
	// e.GET("/swagger", func(c *echo.Context) error {
	// 	// return c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/swagger/index.html")
	// 	return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	// })
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// e.GET("/swagger/index.html", echoSwagger.WrapHandler)
}

type format int

const (
	formatJSON format = iota
	formatYAML
)

func registerOpenapiSpec(f format, spec *openapi3.T) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data []byte
		var err error

		if f == formatJSON {
			data, err = json.Marshal(spec)
		}

		if f == formatYAML {
			data, err = yaml.Marshal(spec)
		}

		if err != nil {
			return err
		}

		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Response().Header().Set("Pragma", "no-cache")
		c.Response().Header().Set("Expires", "0")
		c.Response().Header().Set("Surrogate-Control", "no-store")

		return c.Blob(http.StatusOK, "application/json", data)
	}
}
