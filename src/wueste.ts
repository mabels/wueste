import { Result, ResultError, ResultOK, WithoutResult } from "@adviser/result";
import { Payload } from "./payload";
import { isArrayOrObject } from "./helper";

export class WuesteResult<T, E = Error> extends Result<T, E> {}

export class WuesteResultOK<T> extends ResultOK<T> {}
export class WuesteResultError<T extends Error> extends ResultError<T> {}

export type WuesteWithoutResult<T> = WithoutResult<T>;

export type WuestePayload = Omit<Payload, "Data"> & { Data: unknown };

export type WuestenEncoder = (payload: unknown) => Result<unknown>;

const txtEncoder = new TextEncoder();
export function WuesteJsonBytesEncoder(payload: unknown): Result<unknown> {
  return Result.Ok(txtEncoder.encode(JSON.stringify(payload)));
}

export type WuestenDecoder = (payload: unknown) => Result<unknown>;
const txtDecoder = new TextDecoder();
export function WuesteJsonBytesDecoder(payload: unknown): Result<unknown> {
  try {
    const str = txtDecoder.decode(payload as Uint8Array);
    return Result.Ok(JSON.parse(str));
  } catch (err) {
    return Result.Err(err as Error);
  }
}

export const WuestenJSONPassThroughEncoder = (m: unknown) => Result.Ok(m as unknown);
export const WuestenJSONPassThroughDecoder = (m: unknown) => Result.Ok(m as unknown);

export interface WuestenAttributeParameter<T> {
  readonly base: string;
  readonly varname: string;
  readonly jsonname: string;
  default?: T;

  // setError?: (err : string | Error) => void;
  // format?: string // date-time
}

export type SchemaTypes = "string" | "number" | "integer" | "boolean" | "object" | "array" | "objectitem" | "arrayitem";

export type WuestenReflection =
  | WuestenReflectionObject
  | WuestenReflectionArray
  | WuestenReflectionLiteralNumber
  | WuestenReflectionLiteralInteger
  | WuestenReflectionLiteralBoolean
  | WuestenReflectionLiteralString
  | WuestenReflectionObjectItem
  | WuestenReflectionArrayItem;

// export type WuestenXKeyedMap = Record<string, unknown>;
// export type WuestenXKeyedMap<T extends string= any> = { [P in keyof T]: string extends `x-${T}` ? string : never };
export type WuestenXKeyedMap = Partial<{
  readonly "x-groups": string[];
}>;

export interface WuestenReflectionBase extends WuestenXKeyedMap {
  readonly type: SchemaTypes;
  readonly description?: string;
  readonly ref?: string;
  readonly default?: unknown;
}

// export type WuestenReflectionBase = WuestenReflectionForSchema | WuestenXKeyedMap

export interface WuestenReflectionLiteralInteger extends WuestenReflectionBase {
  readonly type: "integer";
  readonly format?: string;
  readonly default?: number;
}

export interface WuestenReflectionLiteralNumber extends WuestenReflectionBase {
  readonly type: "number";
  readonly format?: string;
  readonly default?: number;
}

export interface WuestenReflectionLiteralBoolean extends WuestenReflectionBase {
  readonly type: "boolean";
  readonly format?: string;
  readonly default?: boolean;
}

export interface WuestenReflectionLiteralString extends WuestenReflectionBase {
  readonly type: "string";
  readonly default?: string;
  readonly format?: string;
}

export interface WuestenReflectionObjectItem {
  readonly type: "objectitem";
  readonly name: string;
  readonly optional: boolean;
  readonly property: WuestenReflection;
  readonly key?: string; // only used for pur object
}

export interface WuestenReflectionArrayItem {
  readonly type: "arrayitem";
  readonly name: string;
  readonly idx: number;
  readonly item: WuestenReflection;
}

export interface WuestenReflectionObject extends WuestenReflectionBase {
  readonly type: "object";
  readonly id?: string;
  readonly title?: string;
  readonly schema?: string;
  readonly properties?: WuestenReflectionObjectItem[];
  readonly required?: string[];
}
export interface WuestenReflectionArray extends WuestenReflectionBase {
  readonly id?: string;
  readonly type: "array";
  readonly items: WuestenReflection;
}

export interface WuestenAttribute<G, I = G> {
  readonly param: WuestenAttributeParameter<G>;
  // SetNameSuffix(...idxs: number[]): void;
  // Reflection(): WuestenReflection;
  CoerceAttribute(val: unknown): Result<G>;
  Coerce(value: I): Result<G>;
  Get(): Result<G>;
}

export class WuestenRetValType {
  constructor(readonly Val: unknown) {}
}

export function WuestenRetVal(val: unknown): WuestenRetValType {
  return new WuestenRetValType(val);
}

export interface WuestenReflectionValue {
  readonly schema: WuestenReflection;
  readonly value: unknown;
}

export type WuestenGetterFn = (path: WuestenReflectionValue[]) => void;

function isOptional(obj: { readonly required?: string[] }, key: string): boolean {
  return obj.required?.indexOf(key) !== -1;
}

function setValue(fn: WuestenGetterFn, v: unknown, path: WuestenReflectionValue[]) {
  const last = path[path.length - 1] as { value: unknown };
  if (Array.isArray(v)) {
    WuestenRecordGetter(fn, path, false);
  } else if (v instanceof Date) {
    last.value = v.toISOString();
  } else if (typeof v === "object" && v !== null) {
    WuestenRecordGetter(fn, path, false);
  } else if (typeof v === "boolean") {
    last.value = v;
  } else if (typeof v === "string") {
    last.value = v;
  } else if (typeof v === "number") {
    last.value = v;
  }
}

export function WuestenRecordGetter(fn: WuestenGetterFn, path: WuestenReflectionValue[], toplevel = true) {
  if (path.length == 0) {
    return;
  }
  const v = path[path.length - 1].value;
  if (Array.isArray(v)) {
    const alevel: WuestenReflectionValue[] = [
      ...path,
      {
        schema: {
          id: "[]",
          type: "array",
        } as WuestenReflectionArray,
        value: v,
      },
    ];
    if (v.length === 0 || toplevel || path[path.length - 1].schema.type === "arrayitem") {
      fn(alevel);
    }
    for (let i = 0; i < v.length; ++i) {
      const myl = [
        ...alevel,
        {
          schema: {
            name: `[${i}]`,
            idx: i,
            type: "arrayitem",
          } as WuestenReflectionArrayItem,
          value: v[i],
        },
      ];
      fn(myl);
      setValue(fn, v[i], myl);
    }
    return;
  } else if (typeof v === "object" && v !== null) {
    const olevel: WuestenReflectionValue[] = [
      ...path,
      {
        schema: {
          id: "{}",
          type: "object",
        },
        value: v,
      },
    ];
    const keys = Object.keys(v).sort();
    if (keys.length === 0 || toplevel || path[path.length - 1].schema.type === "arrayitem") {
      fn(olevel);
    }
    for (const k of keys) {
      const val = (v as Record<string, unknown>)[k];
      const myl: WuestenReflectionValue[] = [
        ...olevel,
        {
          schema: {
            type: "objectitem",
            name: `[${k}]`,
            optional: isOptional(v, k),
            key: k,
          } as WuestenReflectionObjectItem,
          value: val,
        },
      ];
      fn(myl);
      isArrayOrObject(val) && WuestenRecordGetter(fn, myl, false);
    }
    return;
  }
  console.warn(`WuestenRecordGetter: never reached: ${v}`);
  return;
}

export class WuestenGetterBuilder {
  readonly _getterAction: (wgf: WuestenGetterFn) => void;
  constructor(fn: (wgf: WuestenGetterFn) => void) {
    this._getterAction = fn;
  }
  Apply(wgf: WuestenGetterFn) {
    this._getterAction(wgf);
  }
}

export interface WuestenGeneratorFunctions<G, I> {
  readonly coerce: (t: I) => Result<G>;
  // readonly reflection?: WuestenReflection;
}

function coerceAttribute<T, I>(val: unknown, param: WuestenAttributeParameter<T>, coerce: (t: I) => Result<T>): Result<T> {
  const rec = val as WuestenObject;
  for (const key of [param.jsonname, param.varname]) {
    if (rec[key] === undefined || rec[key] === null) {
      continue;
    }
    const my = coerce(rec[key] as I);
    return my;
  }
  if (param.default !== undefined) {
    return coerce(param.default as I);
  }
  return Result.Err(`not found:${param.jsonname}`);
}

export function WuestenAttributeName<T>(param: WuestenAttributeParameter<T>): string {
  const names = [];
  if (param.base) {
    names.push(param.base);
  }
  names.push(param.jsonname);
  return names.join(".");
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export function WuestenCoerceAttribute<T>(val: unknown): Result<T> {
  throw new Error("WuestenCoerceAttribute:Method not implemented.");
  // if (!(typeof val === "object" && val !== null)) {
  //   return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is not an object:` + val);
  // }
  // const res = coerceAttribute<G, I>(val, this.param, this.Coerce.bind(this));
  // if (res.is_err()) {
  //   return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] ${res.unwrap_err().message}`);
  // }
  // return res;
}

export class WuestenAttr<G, I = G> implements WuestenAttribute<G, I> {
  _value?: G;
  // _idxs: number[] = [];
  readonly param: WuestenAttributeParameter<G>;
  readonly _fnParams: WuestenGeneratorFunctions<G, I>;
  constructor(param: WuestenAttributeParameter<I>, fnParams: WuestenGeneratorFunctions<G, I>) {
    let def: G | undefined = undefined;
    this._fnParams = fnParams;
    const result = fnParams.coerce(param.default as I);
    if (result.is_ok()) {
      def = result.unwrap() as G;
    }
    this.param = {
      ...param,
      default: def,
    };
  }
  // Reflection(): WuestenReflection {
  //   throw new Error("Reflection:Method not implemented.");
  // }
  // SetNameSuffix(...idxs: number[]): void {
  //   this._idxs = idxs;
  // }
  CoerceAttribute(val: unknown): Result<G> {
    if (!(typeof val === "object" && val !== null)) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is not an object:` + val);
    }
    const res = coerceAttribute<G, I>(val, this.param, this.Coerce.bind(this));
    if (res.is_err()) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] ${res.unwrap_err().message}`);
    }
    return res;
  }
  Coerce(value: I): Result<G> {
    const result = this._fnParams.coerce(value);
    if (result.is_ok()) {
      this._value = result.unwrap();
      return result;
    }
    return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is ${result.unwrap_err().message}`);
  }
  Get(): Result<G> {
    if (this.param.default === undefined && this._value === undefined) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is required`);
    }
    if (this._value !== undefined) {
      return Result.Ok(this._value);
    }
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return Result.Ok(this.param.default! as G);
  }
}

export class WuestenObjectOptional<B extends WuestenAttribute<T, C>, T, C> implements WuestenAttribute<T, C> {
  readonly typ: B;
  readonly param: WuestenAttributeParameter<T>;
  _value: T;
  constructor(typ: B) {
    this.typ = typ;
    this.param = typ.param;
    this._value = typ.param.default as T;
  }
  Reflection(): WuestenReflection {
    throw new Error("Reflection:Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  CoerceAttribute(val: unknown): Result<T, Error> {
    if (!(typeof val === "object" && val !== null)) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is not an object:` + val);
    }
    const res = coerceAttribute(val, this.param, this.Coerce.bind(this));
    if (res.is_ok()) {
      this._value = res.unwrap() as T;
      return res;
    }
    return Result.Ok(this.param.default as T);
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Coerce(value: C): Result<T, Error> {
    if (value === undefined || value === null) {
      this._value = undefined as T;
      return Result.Ok(this._value);
    }
    const res = this.typ.Coerce(value);
    if (res.is_ok()) {
      this._value = res.unwrap() as unknown as T;
      return Result.Ok(this._value);
    }
    return Result.Err(res.unwrap_err());
  }
  Get(): Result<T, Error> {
    return Result.Ok(this._value);
  }
}

export class WuestenAttrOptional<T, I = T> implements WuestenAttribute<T | undefined, I | undefined> {
  readonly _attr: WuestenAttribute<T | undefined, I | undefined>;
  readonly param: WuestenAttributeParameter<T | undefined>;
  _value: T;
  _idxs: number[] = [];

  constructor(attr: WuestenAttribute<T | undefined, I | undefined>) {
    this._attr = attr;
    this.param = {
      ...attr.param,
      default: attr.param.default as T,
    };
    this._value = attr.param.default as T;
  }
  Reflection(): WuestenReflection {
    throw new Error("Reflection:Method not implemented.");
  }
  // SetNameSuffix(...idxs: number[]): void {
  //   this._idxs = idxs;
  // }
  CoerceAttribute(val: unknown): Result<T | undefined> {
    if (!(typeof val === "object" && val !== null)) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is not an object:` + val);
    }
    const res = coerceAttribute(val, this.param, this.Coerce.bind(this));
    if (res.is_ok()) {
      this._value = res.unwrap() as T;
      return res;
    }
    return Result.Ok(this.param.default as T);
  }

  Coerce(value: I): Result<T> {
    if (value === undefined || value === null) {
      this._value = undefined as T;
      return Result.Ok(this._value);
    }
    const res = this._attr.Coerce(value);
    if (res.is_ok()) {
      this._value = res.unwrap() as unknown as T;
      return Result.Ok(this._value);
    }
    return Result.Err(res.unwrap_err());
  }
  Get(): Result<T> {
    return Result.Ok(this._value);
  }
}

// export interface WuestenSchema {
//   readonly Id: string;
//   readonly Schema: string;
//   readonly Title: string;
// }

export interface WuestenBuilder<T, I> extends WuestenAttribute<T, I> {
  Get(): Result<T>;
}

export interface WuestenFactory<T, I, O> {
  readonly T: T;
  readonly I: I;
  readonly O: O;
  Names(): WuestenNames;
  Builder(param?: WuestenAttributeParameter<I>): WuestenBuilder<T, I>;
  FromPayload(val: WuestePayload, decoder?: WuestenDecoder): Result<T>;
  ToPayload(typ: T, encoder?: WuestenEncoder): Result<WuestePayload>;
  ToObject(typ: T): O; // WuestenObject; keys are json notation
  Clone(typ: T): Result<T>;
  Schema(): WuestenReflection;
  Getter(typ: T, base: WuestenReflectionValue[]): WuestenGetterBuilder;
}
export type WuestenFactoryInferT<F extends WuestenFactory<unknown, unknown, unknown>> =
  F extends WuestenFactory<infer T, unknown, unknown> ? T : never;
export type WuestenFactoryInferI<F extends WuestenFactory<unknown, unknown, unknown>> =
  F extends WuestenFactory<unknown, infer I, unknown> ? I : never;
export type WuestenFactoryInferO<F extends WuestenFactory<unknown, unknown, unknown>> =
  F extends WuestenFactory<unknown, unknown, infer O> ? O : never;

export type WuestenObject = Record<string, unknown>;

export type WuestenFNGetBuilder<T> = (b: T | undefined) => unknown;
export class WuestenObjectBuilder implements WuestenBuilder<WuestenObject, WuestenObject> {
  readonly param: WuestenAttributeParameter<WuestenObject>;
  constructor(param?: WuestenAttributeParameter<WuestenObject>) {
    this.param = param || {
      base: "WuestenObjectBuilder",
      varname: "WuestenObjectBuilder",
      jsonname: "WuestenObjectBuilder",
    };
  }
  Reflection(): WuestenReflection {
    throw new Error("Reflection:Method not implemented.");
  }

  Get(): Result<WuestenObject, Error> {
    throw new Error("WuestenObjectBuilder:Get Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  CoerceAttribute(val: unknown): Result<WuestenObject, Error> {
    throw new Error("WuestenObjectBuilder:CoerceAttribute Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Coerce(value: WuestenObject): Result<WuestenObject, Error> {
    return Result.Ok(value);
    throw new Error("WuestenObjectBuilder:Coerce Method not implemented.");
  }
}

export class WuestenObjectFactoryImpl implements WuestenFactory<WuestenObject, WuestenObject, WuestenObject> {
  readonly T = undefined as unknown as WuestenObject;
  readonly I = undefined as unknown as WuestenObject;
  readonly O = undefined as unknown as WuestenObject;
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Builder(param?: WuestenAttributeParameter<WuestenObject> | undefined): WuestenBuilder<WuestenObject, WuestenObject> {
    return new WuestenObjectBuilder(param);
  }
  Names(): WuestenNames {
    throw new Error("WuestenObjectFactoryImpl:Names Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  FromPayload(val: Payload, decoder?: WuestenDecoder): Result<WuestenObject, Error> {
    throw new Error("WuestenObjectFactoryImpl:FromPayload Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ToPayload(typ: WuestenObject, encoder?: WuestenEncoder): Result<Payload, Error> {
    throw new Error("WuestenObjectFactoryImpl:ToPayload Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ToObject(typ: WuestenObject): WuestenObject {
    throw new Error("WuestenObjectFactoryImpl:ToObject Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Clone(typ: WuestenObject): Result<WuestenObject, Error> {
    throw new Error("WuestenObjectFactoryImpl:Clone Method not implemented.");
  }
  Schema(): WuestenReflection {
    throw new Error("WuestenObjectFactoryImpl:Schema not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Getter(typ: WuestenObject, base: WuestenReflectionValue[]): WuestenGetterBuilder {
    throw new Error("WuestenObjectFactoryImpl:Getter not implemented.");
  }

  // FromPayload(val: Payload, decoder?: WuestenDecoder<WuestenObject> =  WuesteJsonDecoder<Partial<WuestenObject>|Partial<WuestenObject>|Partial<WuestenObject>>)): Result<WuestenObject, Error> {
  //     if (!(val.Type === "https://NestedType" || val.Type === "NestedType")) {
  //       return Result.Err(new Error(`WuestePayload Type mismatch:[https://NestedType,NestedType] != ${val.Type}`));
  //     }
  //     const data = decoder(val.Data)
  //     if (data.is_err()) {
  //       return Result.Err(data.unwrap_err());
  //     }
  //     const builder = new NestedTypeBuilder()
  //     return builder.Coerce(data.unwrap());
  //   }
  // }
  // Clone(typ: WuestenObject): Result<WuestenObject, Error> {
  //   const ret: WuestenObject = {};
  //   for (const key in typ) {
  //       const element = typ[key];
  //       if (typeof element === "object" && element !== null) {
  //         const res = this.Clone(element as WuestenObject);
  //         if (res.is_ok()) {
  //           ret[key] = res.unwrap();
  //         } else {
  //           return res;
  //         }
  //       } else {
  //         ret[key] = element;
  //       }
  //   }
  //   return Result.Ok(ret)
  // }
}
export const WuestenObjectFactory = new WuestenObjectFactoryImpl();

function stringCoerce(value: unknown): Result<string> {
  if (typeof value === "string") {
    return Result.Ok(value);
  }
  if (typeof value === "number") {
    return Result.Ok("" + value);
  }
  if (typeof value === "boolean") {
    return Result.Ok(value ? "true" : "false");
  }
  if (
    (typeof value === "object" || typeof value === "function") &&
    value !== null &&
    typeof (value as { toString: () => string })["toString"] === "function"
  ) {
    return stringCoerce((value as { toString: () => string }).toString());
  }
  try {
    return Result.Err("not a string: " + value);
  } catch (err) {
    return Result.Err("not a string: " + err);
  }
}

function dateTimeCoerce(value: unknown): Result<Date> {
  if (typeof value === "string") {
    return Result.Ok(new Date(value));
  }
  if (typeof value === "number") {
    return Result.Ok(new Date(value));
  }
  if (value instanceof Date) {
    return Result.Ok(value);
  }
  return Result.Err("not a Date: " + value);
}

function booleanCoerce(value: unknown): Result<boolean> {
  if (typeof value === "boolean") {
    return Result.Ok(value);
  }
  if (typeof value === "string") {
    if (["true", "1", "yes", "on"].includes(value.toLowerCase())) {
      return Result.Ok(true);
    }
    if (["false", "0", "no", "off"].includes(value.toLowerCase())) {
      return Result.Ok(false);
    }
  }
  if (typeof value === "number") {
    return Result.Ok(!!value);
  }
  return Result.Err("not a boolean: " + value);
}

export class WuestenAttributeObject<T, I, O> extends WuestenAttr<T, I> {
  private readonly _builder: WuestenAttribute<T, I>;
  constructor(param: WuestenAttributeParameter<I>, factory: WuestenFactory<T, I, O>) {
    const builder = factory.Builder(param);
    super(param, { coerce: builder.Coerce.bind(builder) });
    this._builder = builder;
  }

  Coerce(value: I): Result<T> {
    const res = this._builder.Coerce(value);
    if (res.is_ok()) {
      this._value = res.unwrap();
    }
    return res;
  }

  Get(): Result<T> {
    if (this.param.default === undefined && this._value === undefined) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is required`);
    }
    if (this._value !== undefined) {
      return Result.Ok(this._value);
    }
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return Result.Ok(this.param.default!);
  }
}

function numberCoerce(parse: (i: unknown) => number): (value: unknown) => Result<number> {
  return (value: unknown): Result<number> => {
    const val = parse(value as string);
    if (isNaN(val)) {
      return Result.Err(`not a number: ${value}`);
    }
    return Result.Ok(val);
  };
}

export interface WuesteIteratorNext<T> {
  readonly done?: boolean;
  readonly idx: number;
  readonly value: T;
}
export interface WuesteIterator<T> {
  next(): WuesteIteratorNext<T>;
}

export class WuesteIteratorArray<T> implements WuesteIterator<T> {
  readonly _array: ArrayLike<T>;
  _idx = 0;
  constructor(_array: ArrayLike<T>) {
    this._array = _array;
  }

  next(): WuesteIteratorNext<T> {
    if (this._idx < this._array.length) {
      const idx = this._idx;
      this._idx++;
      return { value: this._array[idx], idx };
    }
    return { done: true, idx: this._idx, value: undefined as unknown as T };
  }
}

export class WuesteIteratorGenerator<T> implements WuesteIterator<T> {
  readonly _iter: IterableIterator<T>;
  _idx = 0;
  constructor(_iter: unknown) {
    this._iter = _iter as IterableIterator<T>;
  }

  next(): WuesteIteratorNext<T> {
    const res = this._iter.next();
    if (!res.done) {
      const idx = this._idx;
      this._idx++;
      return { value: res.value, idx };
    }
    return { done: true, idx: this._idx, value: undefined as unknown as T };
  }
}

export function WuesteToIterator<T>(obj: unknown): Result<WuesteIterator<T>> {
  if (Array.isArray(obj)) {
    return Result.Ok(new WuesteIteratorArray(obj));
  }
  if (typeof obj === "function") {
    obj = obj();
  }
  if (typeof obj === "object" && obj !== null) {
    if (Symbol.iterator in (obj as { [Symbol.iterator]: unknown })) {
      return Result.Ok(new WuesteIteratorGenerator(obj));
    }
    return Result.Ok(new WuesteIteratorArray(Object.values(obj)));
  }
  return Result.Err("not iterable");
}

export type WuesteCoerceTypeDate = Date | string;
export type WuesteCoerceTypeboolean = boolean | string | number;
export type WuesteCoerceTypenumber = number | string;
export type WuesteCoerceTypestring = string | boolean | number | { toString: () => string };

export const wuesten = {
  AttributeString: (def: WuestenAttributeParameter<WuesteCoerceTypestring>): WuestenAttribute<string, WuesteCoerceTypestring> => {
    return new WuestenAttr(def, { coerce: stringCoerce });
  },
  AttributeStringOptional: (
    def: WuestenAttributeParameter<WuesteCoerceTypestring>,
  ): WuestenAttribute<string | undefined, WuesteCoerceTypestring | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: stringCoerce }));
  },

  AttributeDateTime: (def: WuestenAttributeParameter<WuesteCoerceTypeDate>): WuestenAttribute<Date, WuesteCoerceTypeDate> => {
    return new WuestenAttr(def, { coerce: dateTimeCoerce });
  },
  AttributeDateTimeOptional: (
    def: WuestenAttributeParameter<Date | string>,
  ): WuestenAttribute<Date | undefined, WuesteCoerceTypeDate | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: dateTimeCoerce }));
  },

  AttributeInteger: (def: WuestenAttributeParameter<WuesteCoerceTypenumber>): WuestenAttribute<number, WuesteCoerceTypenumber> => {
    return new WuestenAttr(def, { coerce: numberCoerce((a) => parseInt(a as string, 10)) });
  },
  AttributeIntegerOptional: (
    def: WuestenAttributeParameter<WuesteCoerceTypenumber>,
  ): WuestenAttribute<number | undefined, WuesteCoerceTypenumber | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: numberCoerce((a) => parseInt(a as string, 10)) }));
  },

  AttributeNumber: (def: WuestenAttributeParameter<WuesteCoerceTypenumber>): WuestenAttribute<number, WuesteCoerceTypenumber> => {
    return new WuestenAttr(def, { coerce: numberCoerce((a) => parseFloat(a as string)) });
  },
  AttributeNumberOptional: (
    def: WuestenAttributeParameter<WuesteCoerceTypenumber>,
  ): WuestenAttribute<number | undefined, WuesteCoerceTypenumber | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: numberCoerce((a) => parseFloat(a as string)) }));
  },

  AttributeBoolean: (
    def: WuestenAttributeParameter<WuesteCoerceTypeboolean>,
  ): WuestenAttribute<boolean, WuesteCoerceTypeboolean> => {
    return new WuestenAttr(def, { coerce: booleanCoerce });
  },
  AttributeBooleanOptional: (
    def: WuestenAttributeParameter<WuesteCoerceTypeboolean>,
  ): WuestenAttribute<boolean | undefined, WuesteCoerceTypeboolean | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: booleanCoerce }));
  },

  AttributeObject: <E, I, O>(def: WuestenAttributeParameter<I>, factory: WuestenFactory<E, I, O>): WuestenAttribute<E, I> => {
    return new WuestenAttributeObject<E, I, O>(def, factory);
  },
  // AttributeObjectOptional: <E, I, O extends WuestenAttribute<E | undefined, I | undefined>>(o: O): WuestenAttribute<E | undefined, I | undefined> => {
  //   return new WuestenAttrOptional<E, I>(o);
  // },

  AttributeObjectOptional: <E, I, O>(
    def: WuestenAttributeParameter<I>,
    factory: WuestenFactory<E, I, O>,
  ): WuestenAttribute<E | undefined, I | undefined> => {
    return new WuestenAttrOptional<E, I>(new WuestenAttributeObject<E, I, O>(def, factory));
  },

  //   AttributeArray: <T>(): WuestenAttribute<T> => {
  //     return new WuestenAttributeType([] as unknown as T);
  //   },

  //   AttributeArrayOptional: <T>(): WuestenAttribute<T|undefined> => {
  //     return new WuestenAttributeType([] as unknown as T);
  //   }
};

export interface WuestenNames {
  readonly id: string;
  readonly title: string;
  readonly names: string[];
  readonly varname: string;
}

export class WuestenTypeRegistryImpl {
  readonly #registry: Map<string, WuestenFactory<unknown, unknown, unknown>> = new Map();

  Register<T extends WuestenFactory<unknown, unknown, unknown>>(wf: T): T {
    wf.Names().names.forEach((name) => {
      this.#registry.set(name, wf);
    });
    return wf;
  }

  RegisteredNames(): string[] {
    return Array.from(this.#registry.keys());
  }

  GetByName<T extends WuestenFactory<unknown, unknown, unknown>>(name: string): T | undefined {
    return this.#registry.get(name) as T;
  }
}

export const WuestenTypeRegistry = new WuestenTypeRegistryImpl();
