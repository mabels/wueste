import { WuestenAttr, WuestenReflection } from "./wueste";

type Builder<T, P, O> = WuestenAttr<T, Partial<T> | Partial<P> | Partial<O>>;

export function fromEnv<T, P, O>(builder: Builder<T, P, O>, env: Record<string, string>): Builder<T, P, O> {
  const result = builder.Reflection();
  walk(result, [], env);
  return builder;
}

function walk(reflection: WuestenReflection, path: string[], env: Record<string, string>) {
  switch (reflection.type) {
    case "object": {
      if (!reflection.properties) {
        throw new Error("no properties");
      }
      for (const item of reflection.properties) {
        walk(item.property, path.concat(item.name), env);
      }
      break;
    }
    case "array": {
      throw new Error("not implemented");
    }
    case "string":
    case "number":
    case "boolean":
    case "integer": {
      reflection.coerceFromString(env[path.map((p) => p.toUpperCase()).join("_")]);
      break;
    }
    default:
      throw new Error("unknown type");
  }
}
