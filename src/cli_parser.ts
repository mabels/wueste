import { parse, ArgumentConfig, PropertyConfig, TypeConstructor } from "ts-command-line-args";
import {
  WuestenReflection,
  WuestenReflectionObject,
  WuestenReflectionObjectItem,
  WuestenReflectionBase,
  WuestenTypeRegistry,
  WuestenFactory,
  WuestenFactoryInferT,
} from "./wueste";
import { Logger, Result } from "@adviser/cement";
import { isOptional, sanitize, typeName, typePath } from "./js_code_writer";
import { WalkSchemaObjectCollector, walkSchema, walkSchemaFilter } from "./helper";

function toEnvKey(s: string) {
  return sanitize(s).toUpperCase();
}

function cliType(oi: WuestenReflection): TypeConstructor<unknown> {
  switch (oi.type) {
    case "string":
      if (oi.format === "date-time") {
        return (p: string) => new Date(Date.parse(p));
      }
      return String;
    case "integer":
      return Number;
    case "boolean":
      return Boolean;
    case "array":
      return cliType(oi.items);
    case "object":
      return () => ({});
    default:
      throw new Error("unknown type " + oi.type);
  }
}

export interface FromSystemParams {
  readonly args: string[];
  readonly env: Record<string, string>;
  readonly log: Logger;
  readonly parseOut?: typeof console;
}

export interface FromSystemResultHelp {
  readonly isHelp: true;
}

export interface FromSystemResultParsed<F extends WuestenFactory<unknown, unknown, unknown>> {
  readonly isHelp: false;
  readonly parsed: Result<WuestenFactoryInferT<F>>;
}

export type FromSystemResult<F extends WuestenFactory<unknown, unknown, unknown>> =
  | FromSystemResultHelp
  | FromSystemResultParsed<F>;

export function fromSystem<F extends WuestenFactory<unknown, unknown, unknown>>(
  fac: F,
  opts: FromSystemParams,
): FromSystemResult<F> {
  const oc = new WalkSchemaObjectCollector();
  walkSchema(
    fac.Schema(),
    walkSchemaFilter((val: WuestenReflection) => ({ container: val }), oc.add),
  );
  const coerceValue = parseCommandLine(Array.from(oc.objects.values()), opts);
  if (coerceValue.help) {
    return { isHelp: true };
  }
  // console.log("CV=>", coerceValue)
  const builder = fac.Builder();
  builder.Coerce(coerceValue);
  return {
    isHelp: false,
    parsed: builder.Get() as Result<WuestenFactoryInferT<F>>,
  };
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars

function parseCommandLine(types: WuestenReflection[][][], opts: FromSystemParams): Record<string, unknown> {
  const argCfg = {} as ArgumentConfig<Record<string, unknown>>;
  argCfg["help"] = {
    type: Boolean,
    optional: true,
    description: "Prints this usage guide",
  } as unknown as PropertyConfig<unknown>;
  let topType: string | undefined = undefined;
  for (const subTypes of types) {
    for (const subType of subTypes) {
      const obj = subType[subType.length - 2] as WuestenReflectionObject;
      if (!topType) {
        topType = typeName(obj);
      }
      const tpath = typePath(subType);
      const name = tpath.map((o) => o.name).join("-");
      const oi = tpath[tpath.length - 1] as WuestenReflectionObjectItem;

      if (oi.property.type === "object") {
        continue;
      }
      let defaultValue = undefined;
      const envKey = toEnvKey(name);
      if (oi.property.type === "array") {
        const envVal = opts.env[envKey] as string | undefined;
        if (envVal) {
          defaultValue = envVal.split(",").map((v) => v.trim());
        } else {
          const aval = [] as string[];
          for (let i = 0, envVal = opts.env[toEnvKey(`${name}_0`)]; envVal; i++, envVal = opts.env[toEnvKey(`${name}_${i}`)]) {
            aval.push(envVal);
          }
          defaultValue = aval;
        }
      } else {
        defaultValue = opts.env[envKey] || (oi.property as WuestenReflectionBase).default;
      }
      // console.log("DEFAULT", oi.property)
      // const defaultValueStr = defaultValue === undefined ? "" : ` [default: ${defaultValue}]`
      argCfg[name] = {
        type: cliType(oi.property),
        description: `${(oi.property as WuestenReflectionBase).description} [env: ${envKey}]`,
        multiple: oi.property.type === "array",
        optional: isOptional(obj, oi),
        defaultValue,
        path: subType,
      } as PropertyConfig<unknown>;
    }
  }

  const parsed = parse(
    argCfg,
    {
      logger: opts.parseOut || console,
      argv: opts.args,
      partial: true,
      stopAtFirstUnknown: true,
      allowEmpty: false,

      helpArg: "help",
      baseCommand: "node exampleConfigWithHelp",
      headerContentSections: [{ header: topType, content: `Configuration for ${topType}` }],
      // footerContentSections: [{ header: "Footer", content: `Copyright: Big Faceless Corp. inc.` }],
      prependParamOptionsToDescription: true,
    },
    false,
  );
  if (parsed.help) {
    return parsed;
  }

  const coerceType = {} as Record<string, unknown>;
  for (const key of Object.keys(parsed)) {
    const arg = argCfg[key];
    if (!argCfg[key]) {
      opts.log.Warn().Str("key", key).Msg("unknown key");
      continue;
    }
    const path = (arg as unknown as { readonly path: WuestenReflection[] }).path;
    const val = parsed[key];
    setValue(opts.log, val, path, coerceType);
  }
  return coerceType;
}

function setValue(log: Logger, value: unknown, path: WuestenReflection[], result: Record<string, unknown>) {
  const obj = path.shift() as WuestenReflectionObject;
  if (!obj) {
    log.Error().Any("path", path).Msg("can't find object");
    return;
  }
  const oi = path.shift() as WuestenReflectionObjectItem;
  if (!oi) {
    log.Error().Any("path", path).Msg("can't find objectitem");
    return;
  }
  if (oi.property.type === "object") {
    result[oi.name] = result[oi.name] || {};
    setValue(log, value, path, result[oi.name] as Record<string, unknown>);
    return;
  }
  const fac = WuestenTypeRegistry.GetByName(typeName(obj));
  if (!fac) {
    log.Error().Str("type", typeName(obj)).Msg("can't find factory");
    return;
  }
  const builder = fac.Builder();
  const recBuilder = builder as unknown as Record<string, (u: unknown) => unknown>;
  const fname = sanitize(oi.name);
  if (typeof recBuilder[fname] !== "function") {
    log.Error().Str("type", typeName(obj)).Str("fname", fname).Msg("can't find builder function");
    return;
  }
  result[fname] = value;
  return result;
}
