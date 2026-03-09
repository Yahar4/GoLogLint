package main

import (
	"github.com/Yahar4/GoLogLint/pkg/checklogs"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(checklogs.Analyzer)
}
