package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/grindlemire/gothem-stack/pkg/log"
	"github.com/grindlemire/gothem-stack/web/pages/home"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type HomeHandler struct {
	// you could put a database handle here or any dependencies you want
}

func NewHomeHandler() (h *HomeHandler, err error) {
	return &HomeHandler{}, nil
}

// RegisterRoutes registers all the subroutes for the home handler to manage
func (h *HomeHandler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.RenderHomepage)
	g.GET("/random-string", h.GetRandomString)
}

func (h *HomeHandler) RenderHomepage(c echo.Context) error {
	return render(c, home.Page())
}

func (h *HomeHandler) GetRandomString(c echo.Context) error {
	time.Sleep(750 * time.Millisecond)

	err := DoThing()
	if err != nil {
		zap.L().Info("example error", log.Callers(err)...)
		// zap.L().Info("example error with stacktrace", log.Callers(err, log.WithStack())...)
	}

	return render(c, home.RandomString(uuid.NewString()))
}

func DoThing() error {
	return DoSubThing()
}

func DoSubThing() error {
	// wrap third party errors at the callsite to get nice stack traces
	return errors.Wrap(ThirdPartyError(), "wrapped third party error")
}

func ThirdPartyError() error {
	return errors.New("third party error")
}
