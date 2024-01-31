import { FileService } from "@adviser/cement";
import { WuestenReflection, WuestenReflectionArray, WuestenReflectionObject } from "./wueste";

export function fileSystemResolver(fs: FileService) {
  return async (f: string, outerRef?: string): Promise<Record<string, unknown>> => {
    if (f.startsWith("file://")) {
      f = f.slice("file://".length);
    }
    if (outerRef && outerRef.startsWith("file://")) {
      outerRef = fs.dirname(outerRef.slice("file://".length));
    }
    // console.log(`resolver: ${f} ${outerRef}`);
    if (!fs.isAbsolute(f)) {
      f = fs.join(outerRef || "", f);
    }
    const obj = JSON.parse(await fs.readFileString(f));
    obj.$fileref = `file://${f}`;
    obj.id = obj.id || fs.basename(f);
    return obj;
  };
}

export async function jsonSchema2Reflection(
  iSchema: unknown,
  resolver: (f: string, outerRef?: string) => Promise<Record<string, unknown>>,
): Promise<WuestenReflection> {
  if (typeof iSchema !== "object" || iSchema === null) {
    throw new Error("schema must be an object");
  }
  const schema = iSchema as Record<string, unknown> & {
    $ref?: string;
    $fileref?: string;
  };
  // console.log("jsonSchema2Reflection: " + schema.$ref, schema.$fileref);
  if (schema.$ref && !schema.$fileref) {
    const rschema = await resolver(schema.$ref, schema.$fileref);
    // console.log("resolved: " + rschema.$ref, rschema.$fileref);
    return jsonSchema2Reflection(rschema, resolver);
  }
  switch (schema.type) {
    case "object":
      return {
        type: "object",
        id: schema["$id"] as string,
        title: schema.title as string,
        required: schema.required as string[],
        description: schema.description as string,
        properties: await Promise.all(
          Object.entries(schema.properties || []).map(async ([key, val]) => {
            return {
              type: "objectitem",
              name: key,
              optional: ((schema.required || []) as string[]).indexOf(key) === -1,
              property: await jsonSchema2Reflection(val as Record<string, unknown>, (f) => {
                return resolver(f, schema.$fileref);
              }),
            };
          }),
        ),
      } as WuestenReflectionObject;
    case "array":
      return {
        id: schema["$id"] as string,
        type: "array",
        items: await jsonSchema2Reflection(schema.items, (f) => {
          return resolver(f, schema.$fileref);
        }),
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
