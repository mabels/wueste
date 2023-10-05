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
	pflag.BoolVar(&cfg.WriteTestSchema, "write-test-schema", false, "write test schema")
	eg.FromArgs("eg-", &cfg.EntityCfg)

	pflag.CommandLine.Parse(args)

	// uuid := uuid.New().String()
	err := os.MkdirAll(cfg.OutputDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.WriteTestSchema {
		eg.WriteTestSchema(&cfg)
	}
	// defer os.RemoveAll(dir)

	sl := eg.PropertyCtx{
		Registry: eg.NewSchemaRegistry(),
	}
	for _, file := range cfg.InputFiles {
		prop := eg.NewJSONProperty()
		prop.Set("$ref", "file://"+file)
		schema := eg.NewPropertiesBuilder(sl).FromJson(prop).Build()
		if schema.IsErr() {
			log.Fatalf("File:%s with %v", file, schema.Err())
		}
		TsGenerator(&cfg, schema.Ok(), sl)
	}
}
