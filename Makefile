.SILENT:

build:
	go mod tidy && go build -o ./bin/profiles

run: build
	./profiles

test:
	go test ./internal/api... ./internal/database... ./internal/validation/...
