package checklogs

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLowercaseRule(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "lowercase")
}

func TestRussianRule(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "russian")
}

func TestSpecialCharsRule(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "specialchars")
}

func TestSensitiveDataRule(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "sensitive")
}
