import { FileService, Logger } from "@adviser/cement";
import { GenerateGroupConfig } from "./generated/generategroupconfig";
import { fileSystemResolver, jsonSchema2Reflection } from "./json_schema_2_reflection";
import { WalkSchemaObjectCollector, walkSchema, walkSchemaFilter, xFilter } from "./helper";
import { WuestenReflection } from "./wueste";

export interface GenerateGroupTypeParams {
  readonly fs: FileService;
  readonly log: Logger;
  readonly cfg: GenerateGroupConfig;
}

export interface PreparedGenerateGroup {
  readonly inputFile: string;
  readonly outputDir: string;
  readonly includePath: string;
  readonly schema: WuestenReflection;
  readonly oc: WalkSchemaObjectCollector;
}

export async function prepareGenerateGroup(iFile: string, opts: GenerateGroupTypeParams) {
  const inputFile = opts.fs.abs(iFile);
  const outputDir = opts.fs.abs(opts.cfg.output_dir);
  const includePath = opts.fs.abs(opts.cfg.include_path);
  const schema = await jsonSchema2Reflection({ $ref: inputFile }, fileSystemResolver(opts.fs));
  const oc = new WalkSchemaObjectCollector();
  walkSchema(schema, walkSchemaFilter(xFilter(opts.cfg.filter.x_key, opts.cfg.filter.x_value, opts.cfg.not_selected), oc.add));
  return {
    inputFile,
    outputDir,
    includePath,
    schema,
    oc,
  };
}
