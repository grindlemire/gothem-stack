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

	// Render the page to get the HTML string
	html, err := renderStaticPage()
	if err != nil {
		zap.S().Errorf("Error rendering page: %v", err)
		return err
	}

	// Save the HTML string to a file
	err = os.WriteFile("dist/static.html", []byte(html), 0644)
	if err != nil {
		zap.S().Errorf("Error writing HTML to file: %v", err)
		return err
	}

	zap.S().Info("Successfully generated static HTML at dist/static.html")
	return nil
}

func renderStaticPage() (string, error) {
	f, err := os.Open("./web/public/styles.min.css")
	if err != nil {
		return "", errors.Wrap(err, "opening styles")
	}
	defer f.Close()
	var s strings.Builder

	b, err := io.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, "reading styles")
	}

	err = home.StaticPage(string(b)).Render(context.Background(), &s)
	if err != nil {
		return "", errors.Wrap(err, "rendering static page")
	}

	return strings.ReplaceAll(s.String(), "{REPLACE_ME}", string(b)), nil
}
