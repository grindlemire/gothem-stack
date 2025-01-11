package main

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"
	"github.com/grindlemire/gothem-stack/web/pages/home"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func static(ctx context.Context) error {
	// Generate CSS using TailwindCSS
	err := cmd.Run(ctx,
		cmd.WithDir("./web"),
		cmd.WithCMD(
			"node_modules/.bin/tailwindcss",
			"-i", "tailwind.css",
			"-o", "public/styles.min.css",
		),
	)
	if err != nil {
		zap.S().Errorf("Error generating CSS: %v", err)
		return err
	}

	// Generate HTML using templ
	err = cmd.Run(ctx,
		cmd.WithCMD(
			"templ",
			"generate",
		),
	)
	if err != nil {
		zap.S().Errorf("Error generating HTML: %v", err)
		return err
	}

	// delete dist/public/dist if it exists so we can recreate it
	if _, err := os.Stat("dist/public"); err == nil {
		if err := os.RemoveAll("dist/public"); err != nil {
			return errors.Wrap(err, "deleting dist/public")
		}
	}

	// Create dist directories
	if err := os.MkdirAll("dist/public/dist", 0755); err != nil {
		return errors.Wrap(err, "creating dist directories")
	}

	// Copy static assets
	if err := copyStaticAssets(); err != nil {
		return errors.Wrap(err, "copying static assets")
	}

	// Render the page to HTML
	html, err := renderStaticPage()
	if err != nil {
		zap.S().Errorf("Error rendering page: %v", err)
		return err
	}

	// Save the HTML to file
	err = os.WriteFile("dist/public/index.html", []byte(html), 0644)
	if err != nil {
		zap.S().Errorf("Error writing HTML to file: %v", err)
		return err
	}

	zap.S().Info("Successfully generated static site in dist/public")
	return nil
}

func copyStaticAssets() error {
	// Copy favicon
	if err := copyFile(
		"./web/public/favicon.ico",
		"dist/public/dist/favicon.ico",
	); err != nil {
		return errors.Wrap(err, "copying favicon")
	}

	// Copy CSS
	if err := copyFile(
		"./web/public/styles.min.css",
		"dist/public/dist/styles.min.css",
	); err != nil {
		return errors.Wrap(err, "copying styles")
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "opening source file %s", src)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "creating destination file %s", dst)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func renderStaticPage() (string, error) {
	// Render the page with a link to the external CSS file
	var s strings.Builder
	err := home.Page().Render(context.Background(), &s)
	if err != nil {
		return "", errors.Wrap(err, "rendering static page")
	}

	return s.String(), nil
}
