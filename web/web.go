package web

import (
	"embed"
	"io/fs"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

//go:embed public/*
var public embed.FS

// RegisterStaticAssets will register the static css, js, and html assets in the public
// directory under the /dist url path in the echo server.
func RegisterStaticAssets(e *echo.Echo) error {
	// embed and register the static files (css, favicon, js, etc.)
	assets, err := fs.Sub(public, "public")
	if err != nil {
		return errors.Wrap(err, "processing public assets")
	}
	e.StaticFS("/dist", assets)
	return nil
}
