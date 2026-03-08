package checklogs

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "gologlint",
	Doc:      "Checks if logging function are correctly written",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		if !isLogFunction(call, pass) {
			return
		}

		if len(call.Args) == 0 {
			return
		}

		message, ok := extractString(pass, call)
		if !ok {
			return
		}
		checkLogMessage(pass, call.Pos(), message)
	})

	return nil, nil
}

func checkLogMessage(pass *analysis.Pass, pos token.Pos, message string) {
	// check if is not empty
	if len(message) > 0 {
		pass.Reportf(pos, "log-message cant be empty")
	}

	// check if lowercase
	if !unicode.IsLower(rune(message[0])) {
		pass.Reportf(pos, "log-message must be in lowercase")
	}

	// check if written in english
	for _, r := range message {
		if unicode.Is(unicode.Cyrillic, r) {
			pass.Reportf(pos, "log-message must be in english")
			return
		}
	}

	// check if has special symbols
	specialChars := "!?.:;,~@#$%^&*(){}[]<>/|+-*="
	if strings.ContainsAny(message, specialChars) {
		pass.Reportf(pos, "log-message cant contain any special symbols")
	}

	// check if has any sensitive data
	potentialSensitiveData := []string{
		"password", "passwd", "pwd",
		"token", "api_key", "apikey", "secret",
		"key", "credential", "auth",
		"jwt", "access_token", "refresh_token",
		"private_key", "public_key", "certificate",
	}
	for _, sensitive := range potentialSensitiveData {
		if strings.Contains(message, sensitive) {
			pass.Reportf(pos, "log-message cant contain any sensitive data %s", sensitive)
			break
		}
	}
}

func extractString(pass *analysis.Pass, expression ast.Expr) (string, bool) {
	lit, ok := expression.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}

	return strings.Trim(lit.Value, "\""), true
}

func isLogFunction(call *ast.CallExpr, pass *analysis.Pass) bool {
	function := call.Fun

	switch expression := function.(type) {
	case *ast.SelectorExpr:
		isSelectorExpr(expression, pass)
	case *ast.Ident:
		isIdentExpr(expression, pass)
	default:
		return false
	}

	return false
}

func isSelectorExpr(expression *ast.SelectorExpr, pass *analysis.Pass) bool {
	exprTypeX := pass.TypesInfo.TypeOf(expression.X)
	if exprTypeX == nil {
		return false
	}

	if !isLoggerType(exprTypeX) {
		return false
	}

	return isLogMethod(expression.Sel.Name)
}

func isIdentExpr(expression *ast.Ident, pass *analysis.Pass) bool {
	obj := pass.TypesInfo.ObjectOf(expression)
	if obj == nil {
		return false
	}

	function, ok := obj.(*types.Func)
	if !ok {
		return false
	}

	pkg := function.Pkg()
	if pkg == nil {
		return false
	}

	switch pkg.Path() {
	case "log/slog":
		return isLogMethod(expression.Name)
	case "go.uber/zap":
		return isLogMethod(expression.Name)
	default:
		return false
	}
}

func isLoggerType(inputType types.Type) bool {
	for {
		switch t := inputType.(type) {
		case *types.Pointer:
			inputType = t.Elem()
			continue
		case *types.Named:
			inputType = t
		}

		break
	}

	namedType, ok := inputType.(*types.Named)
	if !ok {
		return false
	}

	obj := namedType.Obj()
	if obj == nil {
		return false
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}

	switch pkg.Path() {
	case "log/slog":
		return obj.Name() == "Logger"
	case "go.uber/zap":
		return obj.Name() == "Logger" || obj.Name() == "SugaredLogger"
	default:
		return false
	}
}

// All functions are from official docs of log/slog and go.uber/zap
func isLogMethod(methodName string) bool {
	loggerFunctions := map[string]bool{
		"Info":   true,
		"Debug":  true,
		"Warn":   true,
		"Fatal":  true,
		"Fatalf": true,
		"Panic":  true,
		"DPanic": true,
	}

	return loggerFunctions[methodName]
}
