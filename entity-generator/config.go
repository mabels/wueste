package entity_generator

import "github.com/spf13/pflag"

type Config struct {
	Language    string
	Indent      string
	PackageName string
	FromWueste  string
	FromResult  string
}

type GeneratorConfig struct {
	IncludeDirs     []string
	OutputDir       string
	InputFiles      []string
	EntityCfg       Config
	WriteTestSchema bool
}

func FromArgs(prefix string, cfg *Config) *Config {
	pflag.StringVar(&cfg.Language, prefix+"language", "go", "Language to generate entity for")
	pflag.StringVar(&cfg.Indent, prefix+"indent", "  ", "one indent level")
	pflag.StringVar(&cfg.PackageName, prefix+"package", "please_set_this", "Package name")
	pflag.StringVar(&cfg.FromWueste, prefix+"from-wueste", "wueste/wueste", "Path to wueste")
	pflag.StringVar(&cfg.FromResult, prefix+"from-result", "wueste/result", "Path to result")
	return cfg
}
