module github.com/grindlemire/gothem-stack/magefiles

go 1.22

require (
	github.com/grindlemire/gothem-stack v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/magefile/mage v1.15.0
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.27.0
)

require go.uber.org/multierr v1.11.0 // indirect

replace github.com/grindlemire/gothem-stack => ./..
