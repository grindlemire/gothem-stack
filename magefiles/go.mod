module github.com/grindlemire/htmx-templ-template/magefiles

go 1.22

require (
	github.com/grindlemire/htmx-templ-template v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.27.0
)

require (
	github.com/fatih/color v1.9.0 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/princjef/mageutil v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
)

replace github.com/grindlemire/htmx-templ-template => ./..
