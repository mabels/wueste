import { generateGroupJSONSchema } from "./generate_group_jsonschema";
import { GenerateGroupConfig } from "./generated/generategroupconfig";
import { MockFileService, LoggerImpl } from "@adviser/cement";

it("test generated json-schema", async () => {
  const cfg: GenerateGroupConfig = {
    input_files: ["src/generate_group_type.schema.json"],
    output_dir: "src/generated/keys",
    output_format: "JSchema",
    include_path: "src/generated/wasm",
    filter: {
      x_key: "x-groups",
      x_value: "primary-key",
    },
  };
  const fs = new MockFileService();

  await generateGroupJSONSchema(cfg.input_files[0], {
    log: new LoggerImpl(),
    fs,
    filter: cfg.filter,
    includePath: cfg.include_path,
    outDir: cfg.output_dir,
    notSelected: false,
  });
  expect(Object.keys(fs.files).length).toBe(2);
  expect(Object.values(fs.files)[0].name).toBe(fs.abs("src/generated/keys/generategroupconfigkey.schema.json"));
  expect(JSON.parse(Object.values(fs.files)[0].content)).toEqual({
    $id: "GenerateGroupConfigKey",
    $schema: "http://json-schema.org/draft-07/schema#",
    properties: {
      debug: {
        description: "this is debug",
        type: "string",
        "x-groups": ["primary-key", "secondary-key"],
      },
      "output-dir": {
        default: "./",
        type: "string",
        "x-groups": ["primary-key", "top-key"],
      },
    },
    required: ["output-dir"],
    title: "GenerateGroupConfigKey",
    type: "object",
  });

  expect(Object.values(fs.files)[1].name).toBe(fs.abs("src/generated/keys/generategroupconfig$filterkey.schema.json"));
  expect(JSON.parse(Object.values(fs.files)[1].content)).toEqual({
    $id: "GenerateGroupConfig$FilterKey",
    $schema: "http://json-schema.org/draft-07/schema#",
    properties: {
      "x-key": {
        default: "x-groups",
        type: "string",
        "x-groups": ["primary-key", "sub-key"],
      },
    },
    required: ["x-key"],
    title: "GenerateGroupConfig$FilterKey",
    type: "object",
  });
});
