package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ServerConfig is configuration for the server parsed from the env.
type ServerConfig struct {
	Port       int  `envconfig:"PORT"              default:"4433"`
	LocalCerts bool `envconfig:"LOCAL_CERTS"       default:"false" split_words:"true"`
}

// Run runs the server. The context will be cancelled if we receive a SIGTERM (ctrl-c)
func Run(ctx context.Context) error {
	// parse the env config
	var config ServerConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return errors.Wrap(err, "loading environment")
	}
	addr := fmt.Sprintf(":%d", config.Port)

	if config.LocalCerts {
		if !hasCerts() {
			_, err := generateCerts()
			if err != nil {
				return err
			}
		}
	}

	// create the top level http router
	httpRouter := http.NewServeMux()

	// create our echo router and match all routes to it
	webMux, err := NewRouter(ctx)
	if err != nil {
		return err
	}
	httpRouter.Handle("/", webMux)
	server := &http.Server{Addr: addr, Handler: httpRouter}

	// run the listeners in their own goroutine, this is so we can properly propagate signals
	// and cleanup everything since there may be other signals that need to be cleaned up.
	errCh := make(chan error, 1)
	go func() {
		zap.S().Infof("started listening on %s", addr)
		if config.LocalCerts {
			zap.S().Debug(ctx, "listening with tls")
			err := server.ListenAndServeTLS(publicKeyFile, privateKeyFile)
			errCh <- errors.Wrap(err, "starting server")
			return
		}
		err := server.ListenAndServe()
		errCh <- errors.Wrap(err, "starting server")
	}()

	// wait for either the context to be cancelled indicating we should shutdown
	// or for the servers to fail
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	}
}
