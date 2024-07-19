package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type configKey string

var key = configKey("mageconfigkey")

// GetConfig retrieves the mage config from the context
func GetConfig(ctx context.Context) (config Config, err error) {
	config, ok := ctx.Value(key).(Config)
	if !ok {
		return config, errors.Errorf("config not found in mage context")
	}
	return config, nil
}

// WithConfig adds the mage config to the context
func WithConfig(ctx context.Context, args ...string) context.Context {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		zap.S().Fatalf("unable to parse environment config: %v", err)
	}

	godotenv.Load(fmt.Sprintf("%s.env", config.Env))

	// we have to process the env config again to get the env vars loaded from the env file
	err = envconfig.Process("", &config)
	if err != nil {
		zap.S().Fatalf("unable to parse environment config: %v", err)
	}
	config.Args = args

	return context.WithValue(ctx, key, config)
}
