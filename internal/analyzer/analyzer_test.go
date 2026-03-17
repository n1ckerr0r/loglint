package analyzer

import (
	"testing"

	"github.com/n1ckerr0r/loglint/internal/config"
	"github.com/n1ckerr0r/loglint/internal/rules"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {

	ruleSet := []rules.Rule{
		rules.NewLowercaseRule(),
	}

	a := New(ruleSet, config.Config{})

	analysistest.Run(
		t,
		analysistest.TestData(),
		a,
		"lowercase",
	)
}

func TestAnalyzer_PanicRule(t *testing.T) {

	a := New(
		[]rules.Rule{
			rules.PanicRule{},
		},
		config.Config{},
	)

	analysistest.Run(
		t,
		analysistest.TestData(),
		a,
		"panicrule",
	)
}

func TestAnalyzer_NoRules(t *testing.T) {

	a := New(nil, config.Config{})

	analysistest.Run(
		t,
		analysistest.TestData(),
		a,
		"norules",
	)
}
