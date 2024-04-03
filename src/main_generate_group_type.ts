import { LoggerImpl } from "@adviser/cement";
import { NodeFileService  } from "@adviser/cement/node";
import { generateGroupTSType } from "./generate_group_tstype";
import { generateGroupJSONSchema } from "./generate_group_jsonschema";
import { GenerateGroupConfig, GenerateGroupConfigFactory } from "./generated/generategroupconfig";
import { fromSystem } from "./cli_parser";
import { GenerateGroupTypeParams } from "./generate_group";

const log = new LoggerImpl();
const rcfg = fromSystem(GenerateGroupConfigFactory, {
  log,
  env: Object.keys(process.env).reduce(
    (a, k: string) => {
      a[k] = process.env[k]!;
      return a;
    },
    {} as Record<string, string>,
  ),
  args: process.argv.slice(2),
});

if (rcfg.isHelp) {
  process.exit(0);
}

if (rcfg.parsed.isErr()) {
  log.Error().Err(rcfg.parsed.unwrap_err()).Msg("error in config");
  process.exit(1);
}
const cfg = rcfg.parsed.Ok() as GenerateGroupConfig;
const sfs = new NodeFileService();

(async () => {
  for (const f of cfg.input_files) {
    const gg: GenerateGroupTypeParams = {
      log,
      fs: sfs,
      cfg,
    };
    switch (cfg.output_format) {
      case "JSchema":
        await generateGroupJSONSchema(f, gg);
        break;
      default:
        await generateGroupTSType(f, gg);
    }
  }
})().catch((e) => console.error(e));
