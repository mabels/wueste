import { NodeFileService } from "@adviser/cement/node";
import { fileSystemResolver, jsonSchema2Reflection } from "./json_schema_2_reflection";
import { WuestenReflectionObject } from "./wueste";

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
      name: "not-selected",
      optional: true,
      property: {
        default: false,
        description: "use all which is not filtered",
        type: "boolean",
      },
      type: "objectitem",
    },
    {
      name: "type-name",
      optional: false,
      property: {
        default: "Key",
        description: "name of the type",
        type: "string",
      },
      type: "objectitem",
    },
    {
      name: "output-format",
      optional: false,
      property: {
        default: "TS",
        description: "format TS for Typescript, JSchema for JSON Schema",
        type: "string",
      },
      type: "objectitem",
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
        "x-groups": ["xcluded"],
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
        "x-groups": ["xcluded"],
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
    {
      type: "objectitem",
      name: "filters",
      optional: true,
      property: {
        type: "array",
        items: {
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
    },
  ],
  required: ["input-files", "output-dir", "filter", "include-path", "output-format", "type-name"],
};

it("json_schema_2_reflection", async () => {
  const ref = await jsonSchema2Reflection(
    {
      $ref: "src/generate_group_type.schema.json",
    },
    fileSystemResolver(new NodeFileService()),
  );
  expect(ref).toEqual(refSchema);
});
