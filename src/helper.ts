import {
  WuestenAttr,
  WuestenGetterBuilder,
  WuestenReflection,
  WuestenReflectionObject,
  WuestenReflectionObjectItem,
  WuestenReflectionValue,
} from "./wueste";

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

export function asDottedPath(path: WuestenReflectionValue[]): string {
  return path
    .map(({ schema: r }) => {
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
  ref.Apply((path) => {
    const ref = path[path.length - 1].value;
    const dotted = asDottedPath(path);
    for (const ex of exclude) {
      if (typeof ex === "string" && ex === dotted) {
        return;
      }
      if (ex instanceof RegExp && ex.test(dotted)) {
        return;
      }
    }
    const last = path[path.length - 1].schema;
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

export function toPathValue(a: WuestenReflectionValue[]): unknown {
  if (!Array.isArray(a) || a.length === 0) {
    return undefined;
  }
  const my = a[a.length - 1];
  return my.value;
}

export interface Group {
  readonly path: string;
  readonly schema: WuestenReflection;
  readonly ref: unknown;
}
export interface Groups {
  [key: string]: Group[];
}

export type WalkPathFilterFN<T = unknown> = (key: string, val: T) => T | undefined;

export function xFilter(xName = "x-groups", grp?: string, notSelected = false): WalkPathFilterFN {
  xName = xName.toLocaleLowerCase();
  return (key: string, val: unknown) => {
    if (key.toLocaleLowerCase() === xName && (!grp || (Array.isArray(val) ? val : [val]).includes(grp))) {
      return notSelected ? undefined : val;
    }
    return notSelected ? val : undefined;
  };
}

export function getValueByAttrName<T = unknown>(prop: unknown, fn: WalkPathFilterFN): T | undefined {
  if (typeof prop === "object" && prop !== null) {
    const x = Object.entries(prop).find(([k, v]) => fn(k, v));
    if (x) {
      return x[1];
    }
  }
  return undefined;
}

export type WalkPathFn = (path: WuestenReflection[]) => void;

export function walkSchemaFilter(fn: WalkPathFilterFN, walkFn: WalkPathFn = () => {}): WalkPathFn {
  return (path: WuestenReflection[]) => {
    const rln = path[path.length - 1];
    if (rln.type === "objectitem") {
      const unkProp = rln.property as unknown as Record<string, unknown>;
      const val = getValueByAttrName(unkProp, fn);
      if (val) {
        walkFn(path);
        return rln;
      }
    }
  };
}

export class WalkSchemaObjectCollector {
  readonly objects: Map<string, WuestenReflection[][]> = new Map();

  readonly add = (path: WuestenReflection[]) => {
    if (path.length < 2) {
      throw new Error("path too short");
    }
    const object = path[path.length - 2] as WuestenReflectionObject;
    const objectitem = path[path.length - 1] as WuestenReflectionObjectItem;
    if (object.type !== "object" || objectitem.type !== "objectitem") {
      throw new Error("not an object");
    }
    const id = object.title || object.id;
    if (!id) {
      throw new Error("no id");
    }
    let rfn = this.objects.get(id);
    if (!rfn) {
      rfn = [];
      this.objects.set(id, rfn);
    }
    if (rfn.find((x) => (x[x.length - 1] as WuestenReflectionObjectItem).name === objectitem.name)) {
      return;
    }
    rfn.push(path);
  };
}

export function walkSchema(
  reflection: WuestenReflection,
  walkFn: (path: WuestenReflection[]) => unknown,
  path: WuestenReflection[] = [],
) {
  path = path.concat(reflection);
  switch (reflection.type) {
    case "object":
      // console.log("object", reflection.properties)
      walkFn(path);
      (reflection.properties || [])
        .map((p: WuestenReflection) => walkSchema(p, walkFn, path))
        .filter((r) => r)
        .map((r) => r);
      break;
    case "array":
      walkFn(path);
      walkSchema(reflection.items, walkFn, path);
      break;
    case "arrayitem":
      walkFn(path);
      walkSchema(reflection.item, walkFn, path);
      break;
    case "objectitem":
      walkFn(path);
      walkSchema(reflection.property, walkFn, path);
      break;
  }
}

export function groups(ref: WuestenGetterBuilder, xName = "x-groups"): Groups {
  const groups: Groups = {};
  ref.Apply((path) => {
    const last = path[path.length - 1];
    const ref = last.value;
    const property = ((last.schema as WuestenReflectionObjectItem).property || last.schema) as WuestenReflection;
    const _Groups = getValueByAttrName(property, xFilter(xName));
    if (!_Groups) {
      return;
    }
    const xGroups = Array.isArray(_Groups) ? _Groups : [_Groups];
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
  });
  return groups;
}
