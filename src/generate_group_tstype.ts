import { WuestenReflection, WuestenReflectionObject, WuestenReflectionObjectItem } from "./wueste";
import { JSCodeWriter, sanitize, typeName, isOptional } from "./js_code_writer";
import { GenerateGroupTypeParams, prepareGenerateGroup } from "./generate_group";

function importFileName(typ: WuestenReflection[], suffix?: string): string {
  return importTypeName(typ, suffix).toLowerCase();
}

function importTypeName(typ: WuestenReflection[], suffix?: string): string {
  const oNames = importedType(typ);
  return [...oNames.slice(0, oNames.length - 1), oNames[oNames.length - 1] + (suffix ? suffix : "")].join("$");
}

function importedType(typ: WuestenReflection[]): string[] {
  return typ.filter((o) => o.type === "object").map((o) => typeName(o as WuestenReflectionObject));
}

export async function generateGroupTSType(iFile: string, opts: GenerateGroupTypeParams) {
  const { inputFile, includePath, oc, outputDir } = await prepareGenerateGroup(iFile, opts);

  console.log(
    "generate from: " + opts.fs.relative(inputFile),
    "filtered by:",
    opts.cfg.filter.x_key,
    "/",
    opts.cfg.filter.x_value,
    "not selected:",
    opts.cfg.not_selected,
  );

  for (const typ of Array.from(oc.objects.values())) {
    const obj = typ[0][typ[0].length - 2] as WuestenReflectionObject;
    if (obj.type !== "object") {
      opts.log.Error().Str("type", obj.type).Msg("object expected");
      continue;
    }
    const resultFname = importFileName(typ[0], opts.cfg.type_name) + ".ts";
    const out = await opts.fs.create(opts.fs.join(outputDir, resultFname));
    console.log("  creating file: " + opts.fs.relative(out.name));
    // const log = ctx.log.With().Str("type", typeName(obj)).Logger();
    const w = new JSCodeWriter({
      resultStream: out,
      fileEngine: opts.fs,
    });

    const typeName = `${opts.cfg.type_name}Type`;

    w.block(["export", "interface", importTypeName(typ[0], typeName)], (w) => {
      for (const attr of typ) {
        const oi = attr[attr.length - 1] as WuestenReflectionObjectItem;
        w.writeLn(w.line(w.readonly(w.declareType(oi.property.type, sanitize(oi.name), isOptional(obj, oi)))));
      }
    });
    w.writeLn();
    w.block(["export", "class", importTypeName(typ[0], opts.cfg.type_name)], (w) => {
      w.block(
        w.static(w.declareType(importTypeName(typ[0], typeName), w.call("Coerce", w.declareType(importTypeName(typ[0]), "val")))),
        (w) => {
          w.block(
            "return",
            (w) => {
              for (const attr of typ) {
                const iFile = opts.fs.join(includePath, importFileName(attr));
                w.import(importTypeName(attr), iFile);
                const oi = attr[attr.length - 1] as WuestenReflectionObjectItem;
                w.writeLn(w.tuple(sanitize(oi.name), w.deref("val", sanitize(oi.name))));
              }
            },
            {
              close: "};",
            },
          );
        },
      );
    });
    await w.close();
  }
}
