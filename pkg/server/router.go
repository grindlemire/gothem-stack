package server

import (
	"context"
	"net/http"

	"github.com/grindlemire/gothem-stack/pkg/auth"
	"github.com/grindlemire/gothem-stack/pkg/handler"
	"github.com/grindlemire/gothem-stack/web"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

func NewRouter(ctx context.Context) (h http.Handler, err error) {
	e := echo.New()

	e.Use(
		// recover from panics and create errors from them
		middleware.Recover(),
		// TODO: other global middleware goes here
	)

	// register the customer pages and components
	homeHandler, err := handler.NewHomeHandler()
	if err != nil {
		return h, err
	}
	homeHandler.RegisterRoutes(
		e.Group("", auth.Middleware()),
	)

	// register the static assets like the favicon and the css
	err = web.RegisterStaticAssets(e)
	if err != nil {
		return h, err
	}

	// all other routes should return not found. This should be the last registered route in the list
	e.HTTPErrorHandler = handler.Error
	e.Add(echo.RouteNotFound, "/*", echo.HandlerFunc(func(c echo.Context) error {
		return echo.ErrNotFound.SetInternal(errors.Errorf("not found | uri=[%s]", c.Request().RequestURI))
	}), []echo.MiddlewareFunc{}...)

	return e.Server.Handler, nil
}
