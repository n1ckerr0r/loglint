package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/n1ckerr0r/loglint/internal/analyzer"
	"github.com/n1ckerr0r/loglint/internal/config"
	"github.com/n1ckerr0r/loglint/internal/rules"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type pluginSettings struct {
	CheckLowercase    bool     `json:"check-lowercase"`
	CheckEnglish      bool     `json:"check-english"`
	CheckSpecialChars bool     `json:"check-specialchars"`
	CheckSensitive    bool     `json:"check-sensitive"`
	SuggestedFix      bool     `json:"suggested-fix"`
	SensitiveCase     bool     `json:"sensitive-case"`
	SensitiveKeywords []string `json:"sensitive-keywords"`
}

type Plugin struct {
	cfg config.Config
}

func New(conf any) (register.LinterPlugin, error) {
	if conf == nil {
		conf = map[string]any{}
	}
	s, err := register.DecodeSettings[pluginSettings](conf)
	if err != nil {
		return nil, err
	}
	cfg := config.Config{
		EnableLowercase:    s.CheckLowercase,
		EnableEnglish:      s.CheckEnglish,
		EnableSpecial:      s.CheckSpecialChars,
		EnableSensitive:    s.CheckSensitive,
		EnableSuggestedFix: s.SuggestedFix,
		CaseSensitive:      s.SensitiveCase,
		SensitiveKeywords:  s.SensitiveKeywords,
	}
	return &Plugin{cfg: cfg}, nil
}

func (plugin *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	ruleSet := rules.Build(plugin.cfg)

	return []*analysis.Analyzer{
		analyzer.New(ruleSet, plugin.cfg),
	}, nil
}

func (Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
