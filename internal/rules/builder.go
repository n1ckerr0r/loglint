package rules

import "github.com/n1ckerr0r/loglint/internal/config"

// Для удобного создания правил
func Build(cfg config.Config) []Rule {

	var ruleSet []Rule

	if cfg.EnableLowercase {
		ruleSet = append(ruleSet, NewLowercaseRule())
	}

	if cfg.EnableEnglish {
		ruleSet = append(ruleSet, NewEnglishRule())
	}

	if cfg.EnableSpecial {
		ruleSet = append(ruleSet, NewSpecialCharsRule())
	}

	if cfg.EnableSensitive {

		sensitiveRule := NewSensitiveRule()

		if len(cfg.SensitiveKeywords) > 0 {

			sensitiveRule.CustomizeKeywords(
				cfg.SensitiveKeywords,
				cfg.CaseSensitive,
			)
		}

		ruleSet = append(ruleSet, sensitiveRule)
	}

	return ruleSet
}
