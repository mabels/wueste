import { generateGroupTSType } from "./generate_group_tstype";
import { GenerateGroupConfig, GenerateGroupConfigFactory } from "./generated/generategroupconfig";
import { MockFileService, LoggerImpl } from "@adviser/cement";

export function ansiRegex({ onlyFirst = false } = {}) {
  const pattern = [
    "[\\u001B\\u009B][[\\]()#;?]*(?:(?:(?:(?:;[-a-zA-Z\\d\\/#&.:=?%@~_]+)*|[a-zA-Z\\d]+(?:;[-a-zA-Z\\d\\/#&.:=?%@~_]*)*)?\\u0007)",
    "(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PR-TZcf-nq-uy=><~]))",
  ].join("|");

  return new RegExp(pattern, onlyFirst ? undefined : "g");
}

export function cleanCode(code: string): string[] {
  const reg = ansiRegex();
  return code
    .split("\n")
    .map((i) => i.trim().replace(reg, ""))
    .filter((i) => i);
}

it("test generated key types", async () => {
  const cfg: GenerateGroupConfig = GenerateGroupConfigFactory.Builder()
    .Coerce({
      input_files: ["src/generate_group_type.schema.json"],
      output_dir: "src/generated/keys",
      output_format: "TS",
      type_name: "Xxx",
      include_path: "src/generated/wasm",
      filter: {
        x_key: "x-groups",
        x_value: "primary-key",
      },
    })
    .Ok();
  const fs = new MockFileService();

  await generateGroupTSType(cfg.input_files[0], {
    log: new LoggerImpl(),
    fs,
    cfg,
  });
  expect(Object.keys(fs.files).length).toBe(2);
  expect(Object.values(fs.files)[0].name).toBe(fs.abs("src/generated/keys/generategroupconfigxxx.ts"));
  expect(Object.values(fs.files)[1].name).toBe(fs.abs("src/generated/keys/generategroupconfig$filterxxx.ts"));
  expect(cleanCode(Object.values(fs.files)[0].content)).toEqual([
    'import { GenerateGroupConfig } from "./../wasm/generategroupconfig";',
    "export interface GenerateGroupConfigXxxType {",
    "readonly debug?: string;",
    "readonly output_dir: string;",
    "}",
    "export class GenerateGroupConfigXxx {",
    "static Coerce(val: GenerateGroupConfig): GenerateGroupConfigXxxType {",
    "return {",
    "debug: val.debug,",
    "output_dir: val.output_dir,",
    "};",
    "}",
    "}",
  ]);

  expect(cleanCode(Object.values(fs.files)[1].content)).toEqual([
    'import { GenerateGroupConfig$Filter } from "./../wasm/generategroupconfig$filter";',
    "export interface GenerateGroupConfig$FilterXxxType {",
    "readonly x_key: string;",
    "}",
    "export class GenerateGroupConfig$FilterXxx {",
    "static Coerce(val: GenerateGroupConfig$Filter): GenerateGroupConfig$FilterXxxType {",
    "return {",
    "x_key: val.x_key,",
    "};",
    "}",
    "}",
  ]);
});

it("test generated not_selected key types", async () => {
  const cfg: GenerateGroupConfig = GenerateGroupConfigFactory.Builder()
    .Coerce({
      input_files: ["src/generate_group_type.schema.json"],
      output_dir: "src/generated/keys",
      output_format: "TS",
      type_name: "Xxx",
      include_path: "src/generated/wasm",
      not_selected: true,
      filter: {
        x_value: "xcluded",
      },
    })
    .Ok();
  const fs = new MockFileService();

  await generateGroupTSType(cfg.input_files[0], {
    log: new LoggerImpl(),
    fs,
    cfg,
  });
  expect(Object.keys(fs.files).length).toBe(2);
  expect(Object.values(fs.files)[0].name).toBe(fs.abs("src/generated/keys/generategroupconfigxxx.ts"));
  expect(Object.values(fs.files)[1].name).toBe(fs.abs("src/generated/keys/generategroupconfig$filterxxx.ts"));
  expect(cleanCode(Object.values(fs.files)[0].content)).toEqual([
    'import { GenerateGroupConfig } from "./../wasm/generategroupconfig";',
    "export interface GenerateGroupConfigXxxType {",
    "readonly debug?: string;",
    "readonly not_selected?: boolean;",
    "readonly type_name: string;",
    "readonly output_format: string;",
    "readonly output_dir: string;",
    "readonly input_files: array;",
    "readonly filters?: array;",
    "}",
    "export class GenerateGroupConfigXxx {",
    "static Coerce(val: GenerateGroupConfig): GenerateGroupConfigXxxType {",
    "return {",
    "debug: val.debug,",
    "not_selected: val.not_selected,",
    "type_name: val.type_name,",
    "output_format: val.output_format,",
    "output_dir: val.output_dir,",
    "input_files: val.input_files,",
    "filters: val.filters,",
    "};",
    "}",
    "}",
  ]);

  expect(cleanCode(Object.values(fs.files)[1].content)).toEqual([
    'import { GenerateGroupConfig$Filter } from "./../wasm/generategroupconfig$filter";',
    "export interface GenerateGroupConfig$FilterXxxType {",
    "readonly x_key: string;",
    "readonly x_value: string;",
    "}",
    "export class GenerateGroupConfig$FilterXxx {",
    "static Coerce(val: GenerateGroupConfig$Filter): GenerateGroupConfig$FilterXxxType {",
    "return {",
    "x_key: val.x_key,",
    "x_value: val.x_value,",
    "};",
    "}",
    "}",
  ]);
});
