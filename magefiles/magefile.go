package main

import (
	"context"
	"os"
	"time"

	"github.com/grindlemire/gothem-stack/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func init() {
	log.InitGlobal()
}

// Config is the identifying information pulled out of the environment to execute
// different mage commands
type Config struct {
	Env  string   `envconfig:"env" required:"true"`
	Args []string `envconfig:"args"    default:""`
}

func Install() (err error) {
	defer func(now time.Time) {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(now, err)
	}(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ignore the first two args since they are "mage" and "init"
	return install(WithConfig(ctx, os.Args[2:]...))
}

// Tidy will run go mod tidy
func Tidy() (err error) {
	defer func(now time.Time) {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(now, err)
	}(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ignore the first two args since they are "mage" and "tidy"
	return tidy(WithConfig(ctx, os.Args[2:]...))
}

// Build will build a new binary
func Build() (err error) {
	defer func(now time.Time) {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(now, err)
	}(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ignore the first two args since they are "mage" and "build"
	return build(WithConfig(ctx, os.Args[2:]...))
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

// Templ will run the templ command and generate the go code for the templates
func Templ() (err error) {
	defer func(now time.Time) {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(now, err)
	}(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ignore the first two args since they are "mage" and "templ"
	return templ(WithConfig(ctx, os.Args[2:]...))
}

// Deploy will deploy the static site to Firebase hosting
func Deploy() (err error) {
	defer func(now time.Time) {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(now, err)
	}(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ignore the first two args since they are "mage" and "deploy"
	return deploy(WithConfig(ctx, os.Args[2:]...))
}

// Auth will authenticate with cloud services (gcloud or firebase)
func Auth() (err error) {
	defer func(now time.Time) {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(now, err)
	}(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ignore the first two args since they are "mage" and "auth"
	return auth(WithConfig(ctx, os.Args[2:]...))
}

// Release will deploy a specific version of the service to Cloud Run
func Release() (err error) {
	defer func(now time.Time) {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(now, err)
	}(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ignore the first two args since they are "mage" and "release"
	return release(WithConfig(ctx, os.Args[2:]...))
}

// Destroy will delete all remote cloud infrastructure created during deploy
func Destroy() (err error) {
	defer func(now time.Time) {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(now, err)
	}(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ignore the first two args since they are "mage" and "destroy"
	return destroy(WithConfig(ctx, os.Args[2:]...))
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
