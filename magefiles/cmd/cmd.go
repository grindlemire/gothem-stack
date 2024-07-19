package cmd

import (
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type cmdconf struct {
	name string
	args []string
	dir  string
	env  []string

	logCMD  bool
	infoLog io.Writer
	errLog  io.Writer
	in      io.Reader
}

type cmdopt func(*cmdconf) error

// WithDir specifies the directory to run the subprocess in
func WithDir(dir string) cmdopt {
	return func(conf *cmdconf) (err error) {
		conf.dir = dir
		return nil
	}
}

// WithCMD specifies the command to run and the args to pass
func WithCMD(name string, args ...string) cmdopt {
	return func(conf *cmdconf) (err error) {
		conf.name = name
		if conf.args == nil {
			conf.args = []string{}
		}
		if len(args) > 0 {
			// prepend the command args before the others
			conf.args = append(args, conf.args...)
		}
		return nil
	}
}

// WithSilent discards the subprocess stdout
func WithSilent() cmdopt {
	return func(conf *cmdconf) (err error) {
		conf.infoLog = io.Discard
		conf.errLog = io.Discard
		return nil
	}
}

// WithArgs adds additional arguments for the command
func WithArgs(args ...string) cmdopt {
	return func(conf *cmdconf) (err error) {
		if conf.args == nil {
			conf.args = []string{}
		}
		conf.args = append(conf.args, args...)
		return nil
	}
}

// WithEnv provides additional env vars for the subprocess
func WithEnv(args ...string) cmdopt {
	return func(conf *cmdconf) (err error) {
		if conf.env == nil {
			conf.env = []string{}
		}
		conf.env = append(conf.env, args...)
		return nil
	}
}

// WithLog logs the command that is run and gives the ability to copy
func WithLog() cmdopt {
	return func(conf *cmdconf) (err error) {
		conf.logCMD = true
		return nil
	}
}

// WithLogger hooks up the external zap logger with the stdout and stderr of the subprocess
func WithLogger() cmdopt {
	return func(conf *cmdconf) (err error) {
		infoLog, err := zap.NewStdLogAt(zap.L(), zap.InfoLevel)
		if err != nil {
			return errors.Wrap(err, "wrapping info level zap logger")
		}

		errLog, err := zap.NewStdLogAt(zap.L(), zap.ErrorLevel)
		if err != nil {
			return errors.Wrap(err, "wrapping error level zap logger")
		}
		conf.infoLog = infoLog.Writer()
		conf.errLog = errLog.Writer()
		return nil
	}
}

// CMD creates a command but does not run it. Pass it to Run to run the command.
func CMD(ctx context.Context, opts ...cmdopt) *exec.Cmd {
	conf := &cmdconf{
		infoLog: os.Stdout,
		errLog:  os.Stderr,
		in:      os.Stdin,
	}
	for _, opt := range opts {
		err := opt(conf)
		if err != nil {
			zap.S().Fatalf("creating command config: %v", err)
		}
	}
	cmd := exec.CommandContext(ctx, conf.name, conf.args...)
	cmd.Dir = conf.dir

	cmd.Stderr = conf.errLog
	cmd.Stdout = conf.infoLog
	cmd.Stdin = conf.in
	cmd.Env = append(os.Environ(), conf.env...)

	if conf.logCMD {
		zap.S().Infof("ENV: %s", conf.env)
		zap.S().Infof("DIR: %s", cmd.Dir)
		zap.S().Infof("CMD: %s", cmd.String())
		zap.S().Infof("COPY: pushd %s; %s %s; popd", cmd.Dir, strings.Join(conf.env, " "), cmd.String())
	}
	return cmd
}

// Run a command and wait for it to return
func Run(ctx context.Context, opts ...cmdopt) (err error) {
	cmd := CMD(ctx, opts...)
	// Create a signal channel
	sigCh := make(chan os.Signal, 1)

	// Register a signal handler for SIGINT
	signal.Notify(sigCh, os.Interrupt)

	// start the command and return the result
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "starting subprocess")
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	for {
		select {
		case err = <-done:
			return err
		case <-sigCh:
			zap.S().Debugf("killing subprocess %s", cmd.String())
			err := cmd.Process.Signal(os.Interrupt)
			if err != nil {
				return errors.Wrap(err, "sending cancel signal")
			}
		case <-ctx.Done():
			zap.S().Infof("context cancelled, killing subprocess %s", cmd.String())
			cmd.Process.Signal(os.Interrupt)
			_, err := cmd.Process.Wait()
			return err
		}
	}
}
