package main

import (
	"github.com/Yahar4/GoLogLint/pkg/checklogs"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{checklogs.Analyzer}, nil
}
