import { WuestenAttr, WuestenGetterBuilder, WuestenReflection } from "./wueste";

import { hmac } from "@noble/hashes/hmac";
import { sha1 } from "@noble/hashes/sha1";

// type Builder<T, P, O> = WuestenAttr<T, Partial<T> | Partial<P> | Partial<O>>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export function fromEnv<T, P>(builder: WuestenAttr<T, P>, env: Record<string, string>): WuestenAttr<T, P> {
  // const ref = builder.Reflection();
  // if (ref.type !== "object") {
  //   throw new Error("reflection top type must be object");
  // }
  // walk(ref, [ref.title || ref.id || ""], (path, ref) => {
  //   if (ref.type === "object" || ref.type === "array") {
  //     return;
  //   }
  //   // ref.coerceFromString(env[path.map((p) => p.toUpperCase().replace(/[^A-Za-z0-9]/, "_")).join("_")]);
  // });
  return builder;
}

function asDottedPath(path: WuestenReflection[]): string {
  return path
    .map((r) => {
      switch (r.type) {
        case "object":
          return r.title || r.id || "_";
        case "array":
          return undefined;
        case "objectitem":
          return r.name;
        case "string":
        case "number":
        case "integer":
        case "boolean":
          return undefined;
        default:
          throw new Error("invalid type");
      }
    })
    .filter((i) => i)
    .join(".");
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export async function toHash(ref: WuestenGetterBuilder, exclude: (string | RegExp)[] = []): Promise<Uint8Array> {
  const mac = hmac.create(sha1, "");
  const enc = new TextEncoder();
  ref.Apply((path, ref) => {
    const dotted = asDottedPath(path);
    for (const ex of exclude) {
      if (typeof ex === "string" && ex === dotted) {
        return;
      }
      if (ex instanceof RegExp && ex.test(dotted)) {
        return;
      }
    }
    let val: string | undefined = undefined;
    if (typeof ref === "boolean") {
      val = ref ? "true" : "false";
    } else if (typeof ref === "number") {
      val = ref.toExponential(15);
    } else if (typeof ref === "string") {
      val = ref;
    } else if (ref instanceof Date) {
      val = ref.toISOString();
    }
    // console.log(">>>>>>", dotted, val);
    val && mac.update(enc.encode(val));
  });
  return mac.digest();
}
