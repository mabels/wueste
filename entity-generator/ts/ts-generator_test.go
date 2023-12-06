package ts

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"testing"

	eg "github.com/mabels/wueste/entity-generator"
	"github.com/stretchr/testify/assert"
)

func runCmd(cmdStr string) error {
	split := strings.Split(cmdStr, " ")
	cmd := exec.Command(split[0], split[1:]...)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	go func() {
		merged := io.MultiReader(stderr, stdout)
		scanner := bufio.NewScanner(merged)
		for scanner.Err() == nil && scanner.Scan() {
			msg := scanner.Text()
			fmt.Printf("JS: %s\n", msg)
		}
	}()
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func getConfig() *eg.GeneratorConfig {
	return &eg.GeneratorConfig{
		OutputDir: "../../src/generated/go",
		EntityCfg: eg.Config{
			Indent: "  ",
			// PackageName: "test",
			FromWueste: "../../wueste",
			FromResult: "../../result",
		},
	}
}

func TestTypescript(t *testing.T) {
	cfg := getConfig()
	sl := eg.NewTestContext()

	tfs := eg.TestFlatSchema(sl).Ok()

	tfsObj := tfs.(eg.PropertyObject)
	for _, pi := range tfsObj.Items() {
		if pi.Name() == "sub" {
			pis := pi.Property().(eg.PropertyObject).Items()
			assert.Equal(t, pis[1].Name(), "opt-Test")
			assert.Equal(t, pis[1].Optional(), true)
		}
	}

	TsGenerator(cfg, tfs, sl)
	TsGenerator(cfg, eg.TestSchema(sl), sl)
	// for _, prop := range g.includes.ActiveTypes() {
	// 	if prop.property.IsSome() {
	// 		TsGenerator(cfg, prop.property.Value(), sl)
	// 	}
	// }
	// for _, p := range sl.Registry.Items() {
	// 	// if !p.Written() {
	// 	// 	continue
	// 	// }
	// 	TsGenerator(cfg, p.Property(), sl)
	// }
	err := runCmd("npm run build:js")
	if err != nil {
		t.Fatal(err)
	}
	err = runCmd("npm run test:js ./typescript.test.ts")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMainAction(t *testing.T) {
	cfg := getConfig()
	MainAction([]string{
		"--write-test-schema", "true",
		"--include-dir", "../../src/generated",
		"--input-file", "go/base.schema.json",
		"--input-file", "../../src/generated/go/simple_type.schema.json",
		"--input-file", "../generated/go/nested_type.schema.json",
		"--eg-from-wueste", cfg.EntityCfg.FromWueste,
		"--eg-from-result", cfg.EntityCfg.FromResult,
		"--output-dir", "../../src/generated/go",
	})
}
