package config

// В текущей реализации настраивается из golanci.yml
type Config struct {
	EnableLowercase bool `mapstructure:"check-lowercase"`
	EnableEnglish   bool `mapstructure:"check-english"`
	EnableSpecial   bool `mapstructure:"check-specialchars"`
	EnableSensitive bool `mapstructure:"check-sensitive"`

	SensitiveKeywords []string `mapstructure:"sensitive-keywords"`
	CaseSensitive     bool     `mapstructure:"sensitive-case"`
}
