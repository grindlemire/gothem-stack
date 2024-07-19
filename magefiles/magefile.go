package main

import (
	"context"
	"os"
	"time"

	"github.com/grindlemire/htmx-templ-template/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func init() {
	log.InitGlobal()
}

// Config is the identifying information pulled out of the environment to execute
// different mage commands
type Config struct {
	Env  string   `envconfig:"env" default:"local"`
	Args []string `envconfig:"args"    default:""`
}

// Run will run a local dev server and UI
func Run() (err error) {
	defer func(now time.Time) {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(now, err)
	}(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ignore the first two args since they are "mage" and "run"
	return run(WithConfig(ctx, os.Args[2:]...))
}

func finish(start time.Time, err error) {
	zap.S().Infof("elapsed time: %s", time.Since(start))
	// This is a hack to get around the fact that mage treats command line args
	// as other targets. In a run we just want to interpret them as arugments to the binary,
	// not as other targets. So we just short circuit and tell mage to stop
	if err != nil {
		zap.S().Fatal(err)
	}
	os.Exit(0)
}
