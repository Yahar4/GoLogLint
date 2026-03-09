run:
	go run cmd/gologlint/main.go

build:
	go build -o gologlint

build-plugin:
	go build -buildmode=plugin plugin/gologlint.go
