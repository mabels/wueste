import { WuestenReflection, WuestenReflectionObject, WuestenReflectionObjectItem } from "./wueste";
import { typeName } from "./js_code_writer";
import { WalkSchemaObjectCollector, walkSchema, walkSchemaFilter, xFilter } from "./helper";
import { fileSystemResolver, jsonSchema2Reflection } from "./json_schema_2_reflection";
import { GenerateGroupTypeParams } from "./generate_group_tstype";

function importFileName(typ: WuestenReflection[], suffix?: string): string {
  return importTypeName(typ, suffix).toLowerCase();
}

function importTypeName(typ: WuestenReflection[], suffix?: string): string {
  const oNames = importedType(typ);
  return [...oNames.slice(0, oNames.length - 1), oNames[oNames.length - 1] + (suffix ? suffix : "")].join("$");
}

function importedType(typ: WuestenReflection[], suffix?: string): string[] {
  return typ.filter((o) => o.type === "object").map((o) => typeName(o as WuestenReflectionObject, suffix));
}

export async function generateGroupJSONSchema(iFile: string, opts: GenerateGroupTypeParams) {
  const inputFile = opts.fs.abs(iFile);
  const outDir = opts.fs.abs(opts.outDir);
  const schema = await jsonSchema2Reflection({ $ref: inputFile }, fileSystemResolver(opts.fs));
  const oc = new WalkSchemaObjectCollector();
  walkSchema(schema, walkSchemaFilter(xFilter(opts.filter.x_key, opts.filter.x_value), oc.add));

  console.log("generate from: " + opts.fs.relative(inputFile));

  for (const typ of Array.from(oc.objects.values())) {
    const obj = typ[0][typ[0].length - 2] as WuestenReflectionObject;
    if (obj.type !== "object") {
      opts.log.Error().Str("type", obj.type).Msg("object expected");
      continue;
    }
    const resultFname = importFileName(typ[0], "key") + ".schema.json";
    const out = await opts.fs.create(opts.fs.join(outDir, resultFname));
    console.log("  creating file: " + opts.fs.relative(out.name));
    // const log = ctx.log.With().Str("type", typeName(obj)).Logger();

    const jschema = {
      $schema: "http://json-schema.org/draft-07/schema#",
      type: "object",
      $id: importTypeName(typ[0], "Key"),
      title: importTypeName(typ[0], "Key"),
      properties: {} as Record<string, unknown>,
      required: [] as string[],
    };
    for (const attr of typ) {
      const oi = attr[attr.length - 1] as WuestenReflectionObjectItem;
      if (oi.optional === false) {
        jschema.required.push(oi.name);
      }
      jschema.properties[oi.name] = oi.property;
    }
    await opts.fs.writeFileString(out.name, JSON.stringify(jschema, null, 2));
  }
}
