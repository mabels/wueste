import { generateGroupTSType } from "./generate_group_tstype";
import { GenerateGroupConfig } from "./generated/generategroupconfig";
import { MockFileService, LoggerImpl } from "@adviser/cement";

function cleanCode(code: string): string[] {
  return code
    .split("\n")
    .map((i) => i.trim())
    .filter((i) => i);
}

it("test generated key types", async () => {
  const cfg: GenerateGroupConfig = {
    input_files: ["src/generate_group_type.schema.json"],
    output_dir: "src/generated/keys",
    output_format: "TS",
    include_path: "src/generated/wasm",
    filter: {
      x_key: "x-groups",
      x_value: "primary-key",
    },
  };
  const fs = new MockFileService();

  await generateGroupTSType(cfg.input_files[0], {
    log: new LoggerImpl(),
    fs,
    filter: cfg.filter,
    includePath: cfg.include_path,
    outDir: cfg.output_dir,
    notSelected: false,
  });
  expect(Object.keys(fs.files).length).toBe(2);
  expect(Object.values(fs.files)[0].name).toBe(fs.abs("src/generated/keys/generategroupconfigkey.ts"));
  expect(Object.values(fs.files)[1].name).toBe(fs.abs("src/generated/keys/generategroupconfig$filterkey.ts"));
  expect(cleanCode(Object.values(fs.files)[0].content)).toEqual([
    'import { GenerateGroupConfig } from "./../wasm/generategroupconfig";',
    "export interface GenerateGroupConfigKeyType {",
    "readonly debug?: string;",
    "readonly output_dir: string;",
    "}",
    "export class GenerateGroupConfigKey {",
    "static Coerce(val: GenerateGroupConfig): GenerateGroupConfigKeyType {",
    "return {",
    "debug: val.debug,",
    "output_dir: val.output_dir,",
    "};",
    "}",
    "}",
  ]);

  expect(cleanCode(Object.values(fs.files)[1].content)).toEqual([
    'import { GenerateGroupConfig$Filter } from "./../wasm/generategroupconfig$filter";',
    "export interface GenerateGroupConfig$FilterKeyType {",
    "readonly x_key: string;",
    "}",
    "export class GenerateGroupConfig$FilterKey {",
    "static Coerce(val: GenerateGroupConfig$Filter): GenerateGroupConfig$FilterKeyType {",
    "return {",
    "x_key: val.x_key,",
    "};",
    "}",
    "}",
  ]);
});
