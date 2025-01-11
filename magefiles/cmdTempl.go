package main

import (
	"context"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "generating templ files")
	}

	return nil
}
