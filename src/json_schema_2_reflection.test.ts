import { fileSystemResolver, jsonSchema2Reflection } from "./json_schema_2_reflection";
import { WuestenReflectionObject } from "./wueste";
import { SimpleFileSystem } from "./simple_file_system";

const refSchema: WuestenReflectionObject = {
  id: "GenerateGroupConfig",
  title: "GenerateGroupConfig",
  type: "object",
  properties: [
    {
      type: "objectitem",
      name: "debug",
      property: {
        type: "string",
        description: "this is debug",
        "x-groups": ["primary-key", "secondary-key"],
      },
      optional: true,
    },
    {
      type: "objectitem",
      name: "output-dir",
      property: {
        type: "string",
        default: "./",
        "x-groups": ["primary-key", "top-key"],
      },
      optional: false,
    },
    {
      type: "objectitem",
      name: "include-path",
      property: {
        type: "string",
        default: "./",
      },
      optional: false,
    },
    {
      type: "objectitem",
      name: "input-files",
      optional: false,
      property: {
        // id: "arrayId",
        type: "array",
        items: {
          type: "string",
        },
      },
    },
    {
      type: "objectitem",
      name: "filter",
      optional: false,
      property: {
        type: "object",
        id: "Filter",
        title: "Filter",
        properties: [
          {
            type: "objectitem",
            name: "x-key",
            property: {
              type: "string",
              default: "x-groups",
              "x-groups": ["primary-key", "sub-key"],
            },
            optional: false,
          },
          {
            type: "objectitem",
            name: "x-value",
            property: {
              type: "string",
              default: "primary-key",
            },
            optional: false,
          },
        ],
        required: ["x-key", "x-value"],
      },
    },
  ],
  required: ["input-files", "output-dir", "filter", "include-path"],
};

it("json_schema_2_reflection", async () => {
  const ref = await jsonSchema2Reflection(
    {
      $ref: "src/generate_group_type.schema.json",
    },
    fileSystemResolver(new SimpleFileSystem()),
  );
  expect(ref).toEqual(refSchema);
});
