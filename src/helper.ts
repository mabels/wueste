import { WuestenAttr, WuestenGetterBuilder, WuestenReflection, WuestenReflectionObjectItem, WuestenXKeyedMap } from "./wueste";

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
          return r.title || r.id || "{}";
        case "array":
          return r.id || "[]";
        case "arrayitem":
          return r.name;
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

export function walk<T>(a: T, strategy: (x: unknown) => unknown): unknown {
  if (Array.isArray(a)) {
    const b = strategy(a);
    if (!Array.isArray(b)) {
      return walk(b, strategy);
    }
    const x = [...b];
    for (let i = 0; i < b.length; ++i) {
      x[i] = walk(b[i], strategy); // (sanitize(a[i], strategy, a[i]));
    }
    return x;
  }
  if (typeof a === "object" && a !== null) {
    const b = strategy(a);
    if (Array.isArray(b) || !(typeof b === "object" && b !== null)) {
      return walk(b, strategy);
    }
    if (b === null) {
      return null;
    }
    const y: Record<string, unknown> = { ...b } as Record<string, unknown>;
    for (const k of Object.keys(b)) {
      y[k] = walk((b as Record<string, unknown>)[k] as unknown, strategy);
    }
    return y;
  }
  const b = strategy(a);
  if (typeof b !== typeof a) {
    return walk(b, strategy);
  }
  return b;
}

export function isArrayOrObject(v: unknown): boolean {
  if (v instanceof Date) {
    return false;
  }
  if (Array.isArray(v)) {
    return true;
  }
  if (typeof v === "object" && v !== null) {
    return true;
  }
  return false;
}

const enc = new TextEncoder();
// eslint-disable-next-line @typescript-eslint/no-unused-vars
export function toHash(ref: WuestenGetterBuilder, exclude: (string | RegExp)[] = []): Uint8Array {
  const mac = hmac.create(sha1, "");
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
    const last = path[path.length - 1] as WuestenReflection;
    if (last.type === "object" || last.type === "array") {
      // skip object and array only hash primitives
      return;
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
    if (val) {
      if (last.type === "objectitem" && last.key) {
        mac.update(enc.encode(last.key));
      }
      mac.update(enc.encode(val));
    }
  });
  return mac.digest();
}

export interface Group {
  readonly path: string;
  readonly schema: WuestenReflection;
  readonly ref: unknown;
}
export interface Groups {
  [key: string]: Group[];
}

export function groups(ref: WuestenGetterBuilder) {
  const groups: Groups = {};
  ref.Apply((path, ref) => {
    const last = path[path.length - 1] as WuestenXKeyedMap;
    const property = ((last as WuestenReflectionObjectItem).property || last) as WuestenReflection;
    const xGroups = (property as WuestenXKeyedMap)["x-groups"];
    if (Array.isArray(xGroups)) {
      for (const g of xGroups) {
        if (!groups[g]) {
          groups[g] = [];
        }
        groups[g].push({
          path: asDottedPath(path),
          schema: property,
          ref,
        });
      }
    }
  });
  return groups;
}
