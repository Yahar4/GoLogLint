# gologlint

## A Go linter that checks log messages with predefined rules:
- Must start with a lowercase letter
- Must be in English only
- No special characters or emojis
- No sensitive data

## Works as:
- **standalone app**
- **golangci-lint** plugin

## How it works

`gologlint` is a static analyzer built on Go's [`go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis) framework. It works by:

1. **Parsing Go source code** into an Abstract Syntax Tree (AST)
2. **Finding all function calls** that look like logging (`log.Info`, `slog.Error`, `zap.Debug`, etc.) (supports `log/slog` and `go.uber.org/zap`)
3. **Extracting the log message** from the first argument
4. **Applying rules** to check the message format
5. **Reporting violations** with precise file and line locations

### What happens under the hood

When you run `gologlint` (or `golangci-lint` with the plugin), it:

- Scans all `.go` files in your project
- Identifies logging calls from supported loggers
- Checks each log message against the rules
- Returns diagnostics

### Examples and scenarios

Let's say you have this code:

```go
package main

import "log/slog"

func main() {
    slog.Info("User logged in successfully!")      // uppercase first letter + exclamation mark
    slog.Info("user logged in successfully")       // correct
    
    password := "secret123"
    slog.Info("user password: " + password)        // sensitive data exposed
    
    slog.Error("ошибка подключения к базе данных") // russian language
}
```
After running the linter, you'll see:

```go
main.go:5:2: log message must be in lowercase
main.go:5:2: log message cant contain any special symbols
main.go:8:2: log message cant contain any sensitive data
main.go:10:2: log message must be in english
```

## Installation

### Create a plugin from this linter
1. Download source code
2. From the root directory run `make build-plugin`
3. Copy the generated `gologlint.so` file into your project or to some other known location of your choice.

`make build-plugin` simply runs function from official docs of golangci-lint: 
```Go
go build -buildmode=plugin plugin/gologlint.go
```

### Or you can download it using `go install` and verify the installation

```bash
go install github.com/Yahar4/GoLogLint/cmd/gologlint@latest
gologlint --help 
```
