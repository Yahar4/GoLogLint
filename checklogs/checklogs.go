package checklogs

import (
	"go/ast"
	"go/types"

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

		// message check WIP
	})
	return nil, nil
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
