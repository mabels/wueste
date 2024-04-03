import { generateGroupJSONSchema } from "./generate_group_jsonschema";
import { GenerateGroupConfig, GenerateGroupConfigFactory } from "./generated/generategroupconfig";
import { LoggerImpl } from "@adviser/cement";
import { MockFileService } from "@adviser/cement/node";

it("test generated json-schema", async () => {
  const cfg: GenerateGroupConfig = GenerateGroupConfigFactory.Builder()
    .Coerce({
      input_files: ["src/generate_group_type.schema.json"],
      output_dir: "src/generated/keys",
      output_format: "JSchema",
      type_name: "Xxx",
      include_path: "src/generated/wasm",
      filter: {
        x_key: "x-groups",
        x_value: "primary-key",
      },
    })
    .Ok();
  const fs = new MockFileService();

  await generateGroupJSONSchema(cfg.input_files[0], {
    log: new LoggerImpl(),
    fs,
    cfg,
  });
  expect(Object.keys(fs.files).length).toBe(2);
  expect(Object.values(fs.files)[0].name).toBe(fs.abs("src/generated/keys/generategroupconfigxxx.schema.json"));
  expect(JSON.parse(Object.values(fs.files)[0].content)).toEqual({
    $id: "GenerateGroupConfigXxx",
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
    title: "GenerateGroupConfigXxx",
    type: "object",
  });

  expect(Object.values(fs.files)[1].name).toBe(fs.abs("src/generated/keys/generategroupconfig$filterxxx.schema.json"));
  expect(JSON.parse(Object.values(fs.files)[1].content)).toEqual({
    $id: "GenerateGroupConfig$FilterXxx",
    $schema: "http://json-schema.org/draft-07/schema#",
    properties: {
      "x-key": {
        default: "x-groups",
        type: "string",
        "x-groups": ["primary-key", "sub-key"],
      },
    },
    required: ["x-key"],
    title: "GenerateGroupConfig$FilterXxx",
    type: "object",
  });
});

it("test generated not-selected json-schema", async () => {
  const cfg: GenerateGroupConfig = GenerateGroupConfigFactory.Builder()
    .Coerce({
      input_files: ["src/generate_group_type.schema.json"],
      output_dir: "src/generated/keys",
      output_format: "JSchema",
      type_name: "Xxx",
      include_path: "src/generated/wasm",
      not_selected: true,
      filter: {
        x_key: "x-groups",
      },
    })
    .Ok();
  const fs = new MockFileService();

  await generateGroupJSONSchema(cfg.input_files[0], {
    log: new LoggerImpl(),
    fs,
    cfg,
  });
  expect(Object.keys(fs.files).length).toBe(2);
  expect(Object.values(fs.files)[0].name).toBe(fs.abs("src/generated/keys/generategroupconfigxxx.schema.json"));
  expect(JSON.parse(Object.values(fs.files)[0].content)).toEqual({
    $id: "GenerateGroupConfigXxx",
    $schema: "http://json-schema.org/draft-07/schema#",
    properties: {
      filter: {
        id: "Filter",
        properties: [
          {
            name: "x-key",
            optional: false,
            property: {
              default: "x-groups",
              type: "string",
              "x-groups": ["primary-key", "sub-key"],
            },
            type: "objectitem",
          },
          {
            name: "x-value",
            optional: false,
            property: {
              default: "primary-key",
              type: "string",
            },
            type: "objectitem",
          },
        ],
        required: ["x-key", "x-value"],
        title: "Filter",
        type: "object",
        "x-groups": ["xcluded"],
      },
      filters: {
        items: {
          id: "Filter",
          properties: [
            {
              name: "x-key",
              optional: false,
              property: {
                default: "x-groups",
                type: "string",
                "x-groups": ["primary-key", "sub-key"],
              },
              type: "objectitem",
            },
            {
              name: "x-value",
              optional: false,
              property: {
                default: "primary-key",
                type: "string",
              },
              type: "objectitem",
            },
          ],
          required: ["x-key", "x-value"],
          title: "Filter",
          type: "object",
        },
        type: "array",
      },
      "include-path": {
        default: "./",
        type: "string",
        "x-groups": ["xcluded"],
      },
      "input-files": {
        items: {
          type: "string",
        },
        type: "array",
      },
      "not-selected": {
        default: false,
        description: "use all which is not filtered",
        type: "boolean",
      },
      "output-format": {
        default: "TS",
        description: "format TS for Typescript, JSchema for JSON Schema",
        type: "string",
      },
      "type-name": {
        default: "Key",
        description: "name of the type",
        type: "string",
      },
    },
    required: ["type-name", "output-format", "include-path", "input-files", "filter"],
    title: "GenerateGroupConfigXxx",
    type: "object",
  });

  expect(Object.values(fs.files)[1].name).toBe(fs.abs("src/generated/keys/generategroupconfig$filterxxx.schema.json"));
  expect(JSON.parse(Object.values(fs.files)[1].content)).toEqual({
    $id: "GenerateGroupConfig$FilterXxx",
    $schema: "http://json-schema.org/draft-07/schema#",
    properties: {
      "x-value": {
        default: "primary-key",
        type: "string",
      },
    },
    required: ["x-value"],
    title: "GenerateGroupConfig$FilterXxx",
    type: "object",
  });
});
