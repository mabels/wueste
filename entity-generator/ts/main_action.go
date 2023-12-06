package ts

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"

	eg "github.com/mabels/wueste/entity-generator"
)

func MainAction(args []string, version string, gitCommit string) {
	var cfg eg.GeneratorConfig
	pflag.StringArrayVar(&cfg.IncludeDirs, "include-dir", []string{}, "include directories")
	pflag.StringVar(&cfg.OutputDir, "output-dir", "./", "output directory")
	pflag.StringArrayVar(&cfg.InputFiles, "input-file", []string{}, "input files")
	pflag.BoolVar(&cfg.WriteTestSchema, "write-test-schema", false, "write test schema")
	pflag.BoolVar(&cfg.Version, "version", false, "write version")
	eg.FromArgs("eg-", &cfg.EntityCfg)

	pflag.CommandLine.Parse(args)

	if cfg.Version {
		if version == "" {
			version = "dev"
		}
		if gitCommit == "" {
			gitCommit = "dev"
		}
		fmt.Printf("Version: %s:%s\n", version, gitCommit)
	}

	err := os.MkdirAll(cfg.OutputDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.WriteTestSchema {
		eg.WriteTestSchema(&cfg)
	}

	sl := eg.PropertyCtx{
		Registry: eg.NewSchemaRegistry(eg.NewSchemaLoaderImpl(cfg.IncludeDirs...)),
	}
	for _, file := range cfg.InputFiles {
		prop := eg.NewJSONDict()
		prop.Set("$ref", "file://"+file)
		schema := eg.NewPropertiesBuilder(sl).FromJson(prop).Build()
		if schema.IsErr() {
			log.Fatalf("File:%s with %v", file, schema.Err())
		}
		TsGenerator(&cfg, schema.Ok(), sl)
	}
}
