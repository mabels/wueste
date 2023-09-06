package ts

import (
	"log"
	"os"

	"github.com/spf13/pflag"

	eg "github.com/mabels/wueste/entity-generator"
	"github.com/mabels/wueste/entity-generator/rusty"
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
		schema := eg.NewPropertiesBuilder(sl).Resolve(eg.PropertyRuntime{},
			eg.NewProperty(eg.PropertyParam{
				Ref: rusty.Some("file://" + file),
			}))
		if schema.IsErr() {
			log.Fatal(schema.Err())
		}
		TsGenerator(&cfg, schema.Ok(), sl)
	}
}
