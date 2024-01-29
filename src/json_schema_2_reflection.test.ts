import { jsonSchema2Reflection } from "./json_schema_2_reflection";
import { WuestenReflectionObject } from "./wueste";

const filterSchema = {
  type: "object",
  $id: "Filter",
  title: "Filter",
  properties: {
    xKey: {
      type: "string",
      default: "x-groups",
      "x-groups": ["primary-key", "sub-key"],
    },
    xValue: {
      type: "string",
      default: "primary-key",
    },
  },
  required: ["xKey"],
};

const jsSchema = {
  $id: "GenerateGroupConfig",
  title: "GenerateGroupConfig",
  type: "object",
  properties: {
    debug: {
      type: "string",
      description: "this is debug",
      "x-groups": ["primary-key", "secondary-key"],
    },
    outDir: {
      type: "string",
      default: "./",
      "x-groups": ["primary-key", "top-key"],
    },
    inputFiles: {
      $id: "arrayId",
      type: "array",
      items: {
        type: "string",
      },
    },
    filter: {
      $ref: "file://./test/data/generated/Filter.json",
    },
  },
  required: ["inputFiles", "outDir"],
};

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
      name: "outDir",
      property: {
        type: "string",
        default: "./",
        "x-groups": ["primary-key", "top-key"],
      },
      optional: false,
    },
    {
      type: "objectitem",
      name: "inputFiles",
      optional: false,
      property: {
        id: "arrayId",
        type: "array",
        items: {
          type: "string",
        },
      },
    },
    {
      type: "objectitem",
      name: "filter",
      optional: true,
      property: {
        type: "object",
        id: "Filter",
        title: "Filter",
        properties: [
          {
            type: "objectitem",
            name: "xKey",
            property: {
              type: "string",
              default: "x-groups",
              "x-groups": ["primary-key", "sub-key"],
            },
            optional: false,
          },
          {
            type: "objectitem",
            name: "xValue",
            property: {
              type: "string",
              default: "primary-key",
            },
            optional: true,
          },
        ],
        required: ["xKey"],
      },
    },
  ],
  required: ["inputFiles", "outDir"],
};

it("json_schema_2_reflection", () => {
  const ref = jsonSchema2Reflection(jsSchema, (f) => {
    switch (f) {
      case "file://./test/data/generated/Filter.json":
        return filterSchema;
      default:
        throw new Error("unknown ref " + f);
    }
  });
  expect(ref).toEqual(refSchema);
});
