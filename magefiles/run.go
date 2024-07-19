package main

import (
	"context"

	"github.com/grindlemire/htmx-templ-template/magefiles/cmd"

	"go.uber.org/zap"
)

func run(ctx context.Context) error {
	config, err := GetConfig(ctx)
	if err != nil {
		return err
	}

	if len(config.Args) == 1 && config.Args[0] == "dev" {
		// tailwindcss will recompute and compile the necessary styles if any of the
		// classes change in the templ files
		go func() {
			err = cmd.Run(ctx,
				cmd.WithDir("./web"),
				cmd.WithCMD(
					"npx",
					"tailwindcss",
					"-i", "tailwind.css",
					"-o", "public/styles.min.css",
					"--watch",
				),
			)
		}()

		// templ watch will watch for changes to templ files and regenerate the code
		// as necessary
		go func() {
			err = cmd.Run(ctx,
				cmd.WithCMD(
					"templ",
					"generate",
					"--watch",
					`--proxy=http://localhost:4433`,
					"--open-browser=false",
				),
			)
			if err != nil {
				zap.S().Errorf("Error running templ watch: %v", err)
			}
		}()

		// air will restart the main binary when it detects a change
		// to a templ or go file. It delays 1 second to allow for the templ
		// generation to go first
		err = cmd.Run(ctx,
			cmd.WithCMD(
				"air",
			),
		)
		if err != nil {
			zap.S().Errorf("Error running air server: %v", err)
		}
		return err
	}

	zap.S().Infof("running with config: %+v", config)
	return nil
}
