.SILENT:

build:
	go mod tidy && go build -o ./bin/profiles

run: build
	./bin/profiles

test:
	go test ./internal/api... ./internal/database... ./internal/validation/...
