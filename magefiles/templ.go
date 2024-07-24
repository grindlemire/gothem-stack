package main

import (
	"context"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"

	"go.uber.org/zap"
)

func templ(ctx context.Context) error {
	config, err := GetConfig(ctx)
	if err != nil {
		return err
	}

	zap.S().Infof("mage running with config: %+v", config)

	err = cmd.Run(ctx,
		cmd.WithCMD(
			"templ",
			"generate",
		),
	)
	if err != nil {
		zap.S().Errorf("Error running air server: %v", err)
	}
	return err

}
