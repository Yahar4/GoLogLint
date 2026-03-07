package checklogs

import (
	"go/ast"

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

		var functionName string
		switch function := call.Fun.(type) {
		case *ast.Ident:
			functionName = function.Name
		case *ast.SelectorExpr:
			functionName = function.Sel.Name
		default:
			return
		}

		if isLogFunction(functionName) {
			return
		}

	})
	return nil, nil
}

func isLogFunction(functionName string) bool {
	loggerFunctions := map[string]bool{
		"Info":   true,
		"Debug":  true,
		"Warn":   true,
		"Fatal":  true,
		"Fatalf": true,
		"Panic":  true,
		"DPanic": true,
	}

	return loggerFunctions[functionName]
}
