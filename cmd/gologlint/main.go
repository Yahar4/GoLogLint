package main

import (
	"github.com/Yahar4/pkg/checklogs"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(checklogs.Analyzer)
}
