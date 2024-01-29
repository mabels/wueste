import { WuestenReflection, WuestenReflectionArray, WuestenReflectionObject } from "./wueste";

export function jsonSchema2Reflection(iSchema: unknown, resolver: (f: string) => Record<string, unknown>): WuestenReflection {
  if (typeof iSchema !== "object" || iSchema === null) {
    throw new Error("schema must be an object");
  }
  const schema = iSchema as Record<string, unknown>;
  if (schema.$ref) {
    const ref = resolver(schema.$ref as string);
    return jsonSchema2Reflection(ref, resolver);
  }
  switch (schema.type) {
    case "object":
      return {
        type: "object",
        id: schema["$id"] as string,
        title: schema.title as string,
        required: schema.required as string[],
        description: schema.description as string,
        properties: Object.entries(schema.properties || []).map(([key, val]) => {
          return {
            type: "objectitem",
            name: key,
            optional: ((schema.required || []) as string[]).indexOf(key) === -1,
            property: jsonSchema2Reflection(val as Record<string, unknown>, resolver),
          };
        }),
      } as WuestenReflectionObject;
    case "array":
      return {
        id: schema["$id"] as string,
        type: "array",
        items: schema.items as WuestenReflection,
      } as WuestenReflectionArray;
    case "string":
    case "number":
    case "boolean":
    case "integer":
      return schema as unknown as WuestenReflection;
    default:
      throw new Error("unknown type " + schema.type);
  }
}
