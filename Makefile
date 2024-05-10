.SILENT:

build:
	go mod tidy && go build -o profiles

run: build
	./profiles

test:
	go test ./internal/api... ./internal/database...
