package ts

import (
	"log"
	"os"

	"github.com/spf13/pflag"

	eg "github.com/mabels/wueste/entity-generator"
)

func MainAction(args []string) {
	var cfg eg.GeneratorConfig
	pflag.StringVar(&cfg.OutputDir, "output-dir", "./", "output directory")
	pflag.StringArrayVar(&cfg.InputFiles, "input-file", []string{}, "input files")
	eg.FromArgs("eg-", &cfg.EntityCfg)

	pflag.CommandLine.Parse(args)

	// uuid := uuid.New().String()
	err := os.MkdirAll(cfg.OutputDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	// defer os.RemoveAll(dir)

	sl := eg.NewSchemaLoader()
	for _, file := range cfg.InputFiles {
		schema, err := eg.LoadSchema(file, sl)
		if err != nil {
			log.Fatal(err)
		}
		TsGenerator(&cfg, schema, sl)
	}
}
