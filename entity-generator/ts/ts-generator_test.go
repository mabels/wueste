package ts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	eg "github.com/mabels/wueste/entity-generator"
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
		OutputDir: "../../src-generated/go",
		EntityCfg: eg.Config{
			Indent: "  ",
			// PackageName: "test",
			FromWueste: "../../src/wueste",
			FromResult: "../../src/result",
		},
	}
}

func writeSchema() string {
	cfg := getConfig()
	bytes, _ := json.MarshalIndent(eg.TestJsonFlatSchema(), "", "  ")
	schemaFile := path.Join(cfg.OutputDir, "simple_type.schema.json")
	os.WriteFile(schemaFile, bytes, 0644)
	return schemaFile
}

func TestTypescript(t *testing.T) {
	cfg := getConfig()
	sl := eg.NewTestSchemaLoader()

	writeSchema()

	TsGenerator(cfg, eg.TestFlatSchema(sl), sl)
	TsGenerator(cfg, eg.TestSchema(sl), sl)
	for _, p := range sl.SchemaRegistry().Items() {
		if !p.Written() {
			continue
		}
		TsGenerator(cfg, p.PropertItem().Property(), sl)
	}
	err := runCmd("npm install")
	if err != nil {
		t.Fatal(err)
	}
	err = runCmd("npm run build")
	if err != nil {
		t.Fatal(err)
	}
	err = runCmd("npx jest ./typescript.test.ts")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMainAction(t *testing.T) {
	cfg := getConfig()

	MainAction([]string{
		"--input-file", writeSchema(),
		"--eg-from-wueste", cfg.EntityCfg.FromWueste,
		"--eg-from-result", cfg.EntityCfg.FromResult,
		"--output-dir", "../../src-generated/wasm",
	})
}
