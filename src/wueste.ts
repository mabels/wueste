/* eslint-disable @typescript-eslint/no-explicit-any */
import { Result } from "./result";
import { Payload } from "./payload";
import { walk } from "./helper";

export type WuestePayload = Payload;

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

export function WuestenJSONPassThroughEncoder(m: unknown) {
  return Result.Ok(m as unknown);
}
export function WuestenJSONPassThroughDecoder(m: unknown) {
  return Result.Ok(m as unknown);
}

export class WuestenFormatsHandler {
  readonly _formatsHandler: Record<string, (recv: unknown) => unknown>;
  constructor(formatsHandler?: Record<string, (recv: unknown) => unknown>) {
    this._formatsHandler = { ...formatsHandler };
  }

  add(name: string, fn: (recv: unknown) => unknown): WuestenFormatsHandler {
    this._formatsHandler[name] = fn;
    return this;
  }

  clone(): WuestenFormatsHandler {
    return new WuestenFormatsHandler(this._formatsHandler);
  }
  get(key?: string): undefined | ((recv: unknown) => unknown) {
    return key ? this._formatsHandler[key] : undefined;
  }
}

export interface WuestenFormatFactory {
  readonly formatFactory: WuestenFormatsHandler;
  readonly encoder: (payload: unknown) => Result<unknown>;
  readonly decoder: (payload: unknown) => Result<unknown>;
}

export interface WuestenAttributeParam<C> {
  readonly base: string;
  readonly varname: string;
  readonly jsonname: string;
  default?: C;
  format?: string;
}

export type WuestenAttributeParameter<C> = WuestenAttributeParam<C> & Partial<WuestenFormatFactory>;

export type WuestenAttributeBase<C> = WuestenAttributeParam<C> & WuestenFormatFactory;

export type WuestenConstructionParams<C> = Partial<WuestenAttributeParam<C>> & WuestenFormatFactory;

export type WuestenFactoryParam<C> = Partial<WuestenAttributeParameter<C>>;

export type SchemaTypes = "string" | "number" | "integer" | "boolean" | "object" | "array" | "objectitem";

export type WuestenReflection =
  | WuestenReflectionObject
  | WuestenReflectionArray
  | WuestenReflectionLiteral
  | WuestenReflectionObjectItem;

export interface WuestenReflectionBase {
  readonly type: SchemaTypes;
  readonly ref?: string;
}

export interface WuestenReflectionLiteral extends WuestenReflectionBase {
  readonly type: "string" | "number" | "integer" | "boolean";
  readonly format?: string;
}

export interface WuestenReflectionObjectItem {
  readonly type: "objectitem";
  readonly name: string;
  readonly property: WuestenReflection;
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
  readonly id: string;
  readonly type: "array";
  readonly items: WuestenReflection;
}

export interface WuestenFormatterIf<T, C, O> {
  Coerce(value: C): Result<T>;
  ToObject(value: T): Result<O>;
}

export interface WuestenBuilder<T, C, O> {
  readonly param: WuestenAttributeBase<C>;
  Get(): Result<T>;
  CoerceAttribute(val: unknown): Result<T>;
  Coerce(value: C): WuestenBuilder<T, C, O>;
  ToObject(): Result<O>;
  ToPayload(encoder: WuestenEncoder): Result<Payload>;
}

// export interface WuestenAttribute<T, C, O> extends WuestenBuilder<T, C, O>{
//   // CoerceAttribute(val: unknown): Result<T>;
//   // Get(): Result<T>;
//   // Coerce(value: C): Result<T>;
//   // ToObject(): Result<O>;
// }

export class WuestenRetValType<O> {
  constructor(readonly Val: O) {}
}

export function WuestenRetVal<O>(val: O): WuestenRetValType<O> {
  return new WuestenRetValType(val);
}

export type WuestenGetterFn = (level: WuestenReflection[], value: unknown) => void;

export function WuestenRecordGetter(fn: WuestenGetterFn, level: WuestenReflection[], v: unknown) {
  if (Array.isArray(v)) {
    for (let i = 0; i < v.length; ++i) {
      WuestenRecordGetter(
        fn,
        [
          ...level,
          {
            id: `[${i}]`,
            type: "array",
            items: undefined as unknown as WuestenReflection,
          },
        ],
        v[i],
      );
    }
  } else if (v instanceof Date) {
    fn(level, v.toISOString());
  } else if (typeof v === "object" && v !== null) {
    for (const k of Object.keys(v).sort()) {
      const val = (v as Record<string, unknown>)[k];
      const myl: WuestenReflection[] = [
        ...level,
        {
          type: "objectitem",
          name: k,
          property: undefined as unknown as WuestenReflection,
        },
      ];
      WuestenRecordGetter(fn, myl, k);
      WuestenRecordGetter(fn, myl, val);
    }
  } else if (typeof v === "boolean") {
    fn(level, v);
  } else if (typeof v === "string") {
    fn(level, v);
  } else if (typeof v === "number") {
    fn(level, v);
  }
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

// export interface WuestenGeneratorFunctions<G, I> {
//   readonly coerce: (t: I) => Result<G>;
//   // readonly reflection?: WuestenReflection;
// }

function coerceAttribute<T, C>(val: unknown, param: WuestenAttributeParameter<C>, coerce: (t: C) => Result<T>): Result<T> {
  const rec = val as WuestenObject;
  for (const key of [param.jsonname, param.varname]) {
    if (rec[key] === undefined || rec[key] === null) {
      continue;
    }
    const my = coerce(rec[key] as C);
    return my;
  }
  if (param.default !== undefined) {
    return coerce(param.default as C);
  }
  return Result.Err(`not found:${param.jsonname}`);
}

export function WuestenAttributeName<C>(param: WuestenAttributeParameter<C>): string {
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

export class WuestenAttr<T, C, O> implements WuestenBuilder<T, C, O> {
  _value?: T;
  _default?: T;
  // _idxs: number[] = [];
  readonly param: WuestenAttributeBase<C>;
  // readonly _fnParams: WuestenFormatter<T, C, O>;
  constructor(param: WuestenAttributeParameter<C>) {
    this.param = WuestenAttributeFactory<T, C, O>(param);
    // const formatHandler = this.param.formatsHandler.get(param.format);
    // if (formatHandler) {
    //   throw new Error("WuestenAttr:formatHandler not implemented.");
    //   // this._fnParams = {
    //   //   coerce: (t: I) => {
    //   //     return fnParams.coerce(formatHandler(t) as unknown as I);
    //   //   },
    //   // };
    // } else {
    //   this._fnParams = fnParams;
    // }
    if (param.default !== undefined) {
      const result = this.param.Coerce(param.default);
      if (result.is_ok()) {
        this._default = result.unwrap();
      }
    }
    // this.param = { ...this.param, default: def };
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ToObject(): Result<O> {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ToPayload(encoder: WuestenEncoder): Result<Payload, Error> {
    throw new Error("Method not implemented.");
  }
  // Reflection(): WuestenReflection {
  //   throw new Error("Reflection:Method not implemented.");
  // }
  // SetNameSuffix(...idxs: number[]): void {
  //   this._idxs = idxs;
  // }
  CoerceAttribute(val: unknown): Result<T> {
    if (!(typeof val === "object" && val !== null)) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is not an object:` + val);
    }
    const res = coerceAttribute<T, C>(val, this.param, this.Coerce.bind(this));
    if (res.is_err()) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] ${res.unwrap_err().message}`);
    }
    return res;
  }
  Coerce(value: C): WuestenBuilder<T, C, O> {
    const result = this._fnParams.coerce(value);
    if (result.is_ok()) {
      this._value = result.unwrap();
    } else {
      this._value = undefined as T;
    }
    // return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is ${result.unwrap_err().message}`);
    return this;
  }
  Get(): Result<T> {
    if (this.param.default === undefined && this._value === undefined) {
      return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is required`);
    }
    if (this._value !== undefined) {
      return Result.Ok(this._value);
    }
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return Result.Ok(this.param.default! as T);
  }
}

export function WuestenObjectOptional<F>(f: F): F {
  return f;
}

// export class WuestenObjectOptional<T, C, O> implements WuestenBuilder<T, C, O> {
//   readonly typ: WuestenBuilder<T, C, O>;
//   readonly param: WuestenAttributeBase<C>;
//   _value: T;
//   constructor(typ: WuestenBuilder<T, C, O>) {
//     this.typ = typ;
//     this.param = typ.param;
//     this._value = typ.param.default as T;
//   }
//   ToObject(): Result<O, Error> {
//     throw new Error("Method not implemented.");
//   }
//   // eslint-disable-next-line @typescript-eslint/no-unused-vars
//   ToPayload(encoder: WuestenEncoder): Result<Payload> {
//     throw new Error("Method not implemented.");
//   }
//   // Reflection(): WuestenReflection {
//   //   throw new Error("Reflection:Method not implemented.");
//   // }
//   // eslint-disable-next-line @typescript-eslint/no-unused-vars
//   CoerceAttribute(val: unknown): Result<T> {
//     if (!(typeof val === "object" && val !== null)) {
//       return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is not an object:` + val);
//     }
//     const res = coerceAttribute(val, this.param, this.Coerce.bind(this));
//     if (res.is_ok()) {
//       this._value = res.unwrap() as T;
//       return res;
//     }
//     return Result.Ok(this.param.default as T);
//   }
//   // eslint-disable-next-line @typescript-eslint/no-unused-vars
//   Coerce(value: C): WuestenBuilder<T, C, O> {
//     if (value === undefined || value === null) {
//       this._value = undefined as T;
//       // return Result.Ok(this._value);
//       return this
//     }
//     const res = this.typ.Coerce(value);
//     if (res.is_ok()) {
//       this._value = res.unwrap() as unknown as T;
//       return Result.Ok(this._value);
//     }
//     return Result.Err(res.unwrap_err());
//   }
//   Get(): Result<T, Error> {
//     return Result.Ok(this._value);
//   }
// }

export class WuestenAttrOptional<T, C, O> implements WuestenFactory<T | undefined, C | undefined, O | undefined> {
  readonly _attr: WuestenAttribute<T | undefined, C | undefined>;
  readonly param: WuestenAttributeBase<C | undefined>;
  _value: T;
  _idxs: number[] = [];

  constructor(attr: WuestenAttribute<T | undefined, C | undefined>) {
    this._attr = attr;
    this.param = WuestenAttributeFactoryOptional<T | undefined, C | undefined>(attr.param);
    this._value = attr.param.default as T;
  }
  // Reflection(): WuestenReflection {
  //   throw new Error("Reflection:Method not implemented.");
  // }
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

  Coerce(value: C): Result<T> {
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

export interface WuestenNames {
  readonly id: string;
  readonly title: string;
  readonly names: string[];
  readonly varname: string;
}

// export type WuestenBuilderResult<R extends WuestenBuilder<T, C, O>, T, C, O> = R

export abstract class WuestenFactory<T, C, O> {
  readonly _params: WuestenAttributeBase<C>;
  constructor(param: WuestenAttributeBase<C>) {
    this._params = WuestenAttributeFactory<T, C, O>(param);
  }
  abstract Names(): WuestenNames;
  abstract Builder(base: WuestenAttributeBase<unknown>, params: WuestenAttributeParameter<C>): WuestenBuilder<T, C, O>;
  abstract FromPayload(val: Payload, decoder?: WuestenDecoder): Result<WuestenBuilder<T, C, O>>;
  // abstract ToPayload(typ: T, encoder?: WuestenEncoder): Result<Payload>;
  // abstract ToObject(typ: T): O; // WuestenObject; keys are json notation
  abstract Clone(typ: T): Result<T>;
  abstract Schema(): WuestenReflection;
  abstract Getter(typ: T, base: WuestenReflection[]): WuestenGetterBuilder;
  AddFormat(name: string, fn: (recv: unknown) => unknown): WuestenFactory<T, C, O> {
    this._params.formatFactory.add(name, fn);
    return this;
  }
}

export type WuestenObject = Record<string, unknown>;

export type WuestenFNGetBuilder<C, O> = (b: C | undefined) => WuestenRetValType<O> | unknown;
export class WuestenObjectBuilder implements WuestenBuilder<WuestenObject, WuestenObject, WuestenObject> {
  readonly param: WuestenAttributeBase<WuestenObject>;
  constructor(param?: WuestenFactoryParam<WuestenObject>) {
    this.param = WuestenAttributeFactory({
      base: "WuestenObjectBuilder",
      varname: "WuestenObjectBuilder",
      jsonname: "WuestenObjectBuilder",
      ...param,
    });
  }
  Get(): Result<WuestenObject> {
    throw new Error("WuestenObjectBuilder:Get Method not implemented.");
  }
  // // eslint-disable-next-line @typescript-eslint/no-unused-vars
  // CoerceAttribute(val: unknown): Result<WuestenObject, Error> {
  //   throw new Error("WuestenObjectBuilder:CoerceAttribute Method not implemented.");
  // }
  ToObject(): Result<WuestenObject> {
    throw new Error("WuestenObjectBuilder:ToObject Method not implemented.");
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ToPayload(encoder: WuestenEncoder): Result<Payload, Error> {
    throw new Error("WuestenObjectBuilder:ToObject Method not implemented.");
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  CoerceAttribute(val: unknown): Result<WuestenObject, Error> {
    throw new Error("WuestenObjectBuilder:CoerceAttribute Method not implemented.");
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Coerce(value: WuestenObject): WuestenObjectBuilder {
    // return Result.Ok(value);
    throw new Error("WuestenObjectBuilder:Coerce Method not implemented.");
  }
}

export class WuestenObjectFactoryImpl extends WuestenFactory<WuestenObject, WuestenObject, WuestenObject> {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars

  Builder(param?: WuestenAttributeParameter<WuestenObject>): WuestenObjectBuilder {
    return new WuestenObjectBuilder(param);
  }
  Names(): WuestenNames {
    throw new Error("WuestenObjectFactoryImpl:Names Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  FromPayload(val: Payload, decoder?: WuestenDecoder): Result<WuestenBuilder<WuestenObject, WuestenObject, WuestenObject>> {
    throw new Error("WuestenObjectFactoryImpl:FromPayload Method not implemented.");
  }
  // // eslint-disable-next-line @typescript-eslint/no-unused-vars
  // ToPayload(typ: WuestenObject, encoder?: WuestenEncoder): Result<Payload, Error> {
  //   throw new Error("WuestenObjectFactoryImpl:ToPayload Method not implemented.");
  // }
  // // eslint-disable-next-line @typescript-eslint/no-unused-vars
  // ToObject(typ: WuestenObject): WuestenObject {
  //   throw new Error("WuestenObjectFactoryImpl:ToObject Method not implemented.");
  // }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Clone(typ: WuestenObject): Result<WuestenObject, Error> {
    throw new Error("WuestenObjectFactoryImpl:Clone Method not implemented.");
  }
  Schema(): WuestenReflection {
    throw new Error("WuestenObjectFactoryImpl:Schema not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Getter(typ: WuestenObject, base: WuestenReflection[]): WuestenGetterBuilder {
    throw new Error("WuestenObjectFactoryImpl:Getter not implemented.");
  }
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

// export type WuesteCoerceTypeDate = Date | string;
// export type WuesteCoerceTypeboolean = boolean | string | number;
// export type WuesteCoerceTypenumber = number | string;
// export type WuesteCoerceTypestring = string | boolean | number | { toString: () => string };

// export const wuesten = {
//   AttributeString: (
//     def: WuestenAttributeParameter<StringFormatter.CoerceType>,
//   ): WuestenAttribute<string, WuesteCoerceTypestring> => {
//     return new WuestenAttr(def, { coerce: stringCoerce });
//   },
//   AttributeStringOptional: (
//     def: WuestenAttributeParameter<string, WuesteCoerceTypestring>,
//   ): WuestenAttribute<string | undefined, WuesteCoerceTypestring | undefined> => {
//     return new WuestenAttrOptional(new WuestenAttr(def, { coerce: stringCoerce }));
//   },

//   AttributeDateTime: (
//     def: WuestenAttributeParameter<WuesteCoerceTypeDate, WuesteCoerceTypeDate>,
//   ): WuestenAttribute<Date, WuesteCoerceTypeDate> => {
//     return new WuestenAttr(def, { coerce: dateTimeCoerce }) as WuestenAttribute<Date, WuesteCoerceTypeDate>;
//   },
//   AttributeDateTimeOptional: (
//     def: WuestenAttributeParameter<WuesteCoerceTypeDate | undefined, WuesteCoerceTypeDate | undefined>,
//   ): WuestenAttribute<Date | undefined, WuesteCoerceTypeDate | undefined> => {
//     return new WuestenAttrOptional(new WuestenAttr(def, { coerce: dateTimeCoerce })) as WuestenAttribute<
//       Date | undefined,
//       WuesteCoerceTypeDate | undefined
//     >;
//   },

//   AttributeInteger: (
//     def: WuestenAttributeParameter<number, WuesteCoerceTypenumber>,
//   ): WuestenAttribute<number, WuesteCoerceTypenumber> => {
//     return new WuestenAttr(def, { coerce: numberCoerce((a) => parseInt(a as string, 10)) });
//   },
//   AttributeIntegerOptional: (
//     def: WuestenAttributeParameter<number, WuesteCoerceTypenumber>,
//   ): WuestenAttribute<number | undefined, WuesteCoerceTypenumber | undefined> => {
//     return new WuestenAttrOptional(new WuestenAttr(def, { coerce: numberCoerce((a) => parseInt(a as string, 10)) }));
//   },

//   AttributeNumber: (
//     def: WuestenAttributeParameter<number, WuesteCoerceTypenumber>,
//   ): WuestenAttribute<number, WuesteCoerceTypenumber> => {
//     return new WuestenAttr(def, { coerce: numberCoerce((a) => parseFloat(a as string)) });
//   },
//   AttributeNumberOptional: (
//     def: WuestenAttributeParameter<number, WuesteCoerceTypenumber>,
//   ): WuestenAttribute<number | undefined, WuesteCoerceTypenumber | undefined> => {
//     return new WuestenAttrOptional(new WuestenAttr(def, { coerce: numberCoerce((a) => parseFloat(a as string)) }));
//   },

//   AttributeBoolean: (
//     def: WuestenAttributeParameter<boolean, WuesteCoerceTypeboolean>,
//   ): WuestenAttribute<boolean, WuesteCoerceTypeboolean> => {
//     return new WuestenAttr(def, { coerce: booleanCoerce });
//   },
//   AttributeBooleanOptional: (
//     def: WuestenAttributeParameter<boolean, WuesteCoerceTypeboolean>,
//   ): WuestenAttribute<boolean | undefined, WuesteCoerceTypeboolean | undefined> => {
//     return new WuestenAttrOptional(new WuestenAttr(def, { coerce: booleanCoerce }));
//   },

//   AttributeObject: <E, I, O>(def: WuestenAttributeParameter<E, I>, factory: WuestenFactory<E, I, O>): WuestenAttribute<E, I> => {
//     return new WuestenAttributeObject<E, I, O>(def, factory);
//   },

//   AttributeObjectOptional: <E, I, O>(
//     def: WuestenAttributeParameter<E, I>,
//     factory: WuestenFactory<E, I, O>,
//   ): WuestenAttribute<E | undefined, I | undefined> => {
//     return new WuestenAttrOptional<E, I>(new WuestenAttributeObject<E, I, O>(def, factory));
//   },
// };

export class WuestenArrayFactory<T, C, O, IT, IC, IO> implements WuesteCreateBuilder<T, C, O> {
  readonly item: WuesteCreateBuilder<IT, IC, IO>;
  constructor(item: WuesteCreateBuilder<IT, IC, IO>) {
    this.item = item;
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Builder(factory: WuestenAttributeBase<unknown>, params: WuestenFactoryParam<C>): WuestenArrayBuilder<T, C, O, IT, IC, IO> {
    return new WuestenArrayBuilder<T, C, O, IT, IC, IO>(this.item, factory, params);
  }
}

export class WuestenArrayFactoryOptional<T, C, O, IT, IC, IO> extends WuestenArrayFactory<T, C, O, IT, IC, IO> {}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export function WuestenArrayGetterWalk<T>(v: T, base: WuestenReflection[] = []): WuestenGetterBuilder {
  throw new Error("WuestenArrayGetterWalk:Method not implemented.");
}

export class WuestenArrayBuilder<T, C, O, IT, IC, IO> implements WuestenBuilder<T, C, O> {
  readonly param: WuestenAttributeBase<C>;
  private readonly _itemFactory: WuesteCreateBuilder<IT, IC, IO>;
  private readonly _values: Result<IT>[] = [];
  // readonly _value: unknown[] = [];
  // readonly _errors: Error[] = [];
  private readonly _factory: WuestenAttributeBase<unknown>;
  private readonly _params: WuestenFactoryParam<C>;
  constructor(
    itemFactory: WuesteCreateBuilder<IT, IC, IO>,
    factory: WuestenAttributeBase<unknown>,
    params: WuestenFactoryParam<C>,
  ) {
    // this.param = item.param;
    this.param = undefined as unknown as WuestenAttributeBase<C>;
    // this._format = {} as WuestenFormatterIf<T, C, O>;
    this._itemFactory = itemFactory;
    this._params = params;
    this._factory = factory;
  }
  Get(): Result<T, Error> {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  CoerceAttribute(val: unknown): Result<T, Error> {
    throw new Error("Method not implemented.");
  }
  Coerce(value: C): WuestenBuilder<T, C, O> {
    // const ret: unknown[] = []
    walk(value, (item) => {
      if (Array.isArray(item)) {
        elems.push([]);
        elems = elems[elems.length - 1] as unknown[];
        return;
      }
      const attr = this.#itemFactory.New();
      elems.push(attr.Coerce(item as C).Get());
    });
    return this;
  }
  ToObject(): Result<O, Error> {
    const errs: Error[] = [];
    const ret: unknown[] = [];
    let elems = ret;
    walk(this._value, (item) => {
      if (Array.isArray(item)) {
        elems.push([]);
        elems = elems[elems.length - 1] as unknown[];
        return;
      }
      const r = this._format.ToObject(item as T);
      if (r.is_err()) {
        errs.push(r.unwrap_err());
      } else {
        elems.push(r.unwrap());
      }
    });
    if (errs.length > 0) {
      return Result.Err(errs.join("\n"));
    }
    return Result.Ok(ret as O);
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ToPayload(encoder: WuestenEncoder): Result<Payload, Error> {
    throw new Error("Method not implemented.");
  }
}

export class WuestenArrayBuilderOptional<T, C, O, IT, IC, IO> extends WuestenArrayBuilder<T, C, O, IT, IC, IO> {}

export class WuestenTypeRegistryImpl {
  readonly _registry: Record<string, WuestenFactory<unknown, unknown, unknown>> = {};
  readonly _attributeBase: WuestenFormatFactory = {
    formatFactory: new WuestenFormatsHandler(),
    encoder: WuestenJSONPassThroughEncoder,
    decoder: WuestenJSONPassThroughDecoder,
  };

  Register<T>(factory: T): T {
    const wf = factory as unknown as WuestenFactory<unknown, unknown, unknown>;
    wf.Names().names.forEach((name) => {
      this._registry[name] = wf;
    });
    return factory;
  }

  AddFormat(name: string, fn: (recv: unknown) => unknown) {
    this._attributeBase.formatFactory.add(name, fn);
    return this;
  }

  _deleteUndefined(merge?: Partial<WuestenFormatFactory>): Partial<WuestenFormatFactory> {
    const m = { ...merge };
    if (!m.decoder) {
      delete m.decoder;
    }
    if (!m.encoder) {
      delete m.encoder;
    }
    if (!m.formatFactory) {
      delete m.formatFactory;
    }
    return m;
  }

  cloneAttributeBase<C>(...merges: (WuestenAttributeParameter<unknown> | undefined)[]): WuestenAttributeBase<C> {
    return {
      ...merges.reduce(
        (acc, cur) => {
          return { ...acc, ...cur };
        },
        {
          ...this._attributeBase,
        },
      ),
      formatsHandler: this._attributeBase.formatsHandler.clone(),
    } as WuestenAttributeBase<C>;
  }
}

export function WuestenMergeAttributeBase<C>(
  factory: WuestenFormatFactory,
  ...params: (undefined | WuestenAttributeParameter<unknown>)[]
): WuestenAttributeBase<C> {
  return params.reduce(
    (acc, cur) => {
      return { ...acc, ...cur };
    },
    { ...factory },
  ) as WuestenAttributeBase<C>;
}

export function WuestenCoerceHelper<T, C>(val: unknown, param: WuestenAttributeBase<C>, coerce: (t: C) => Result<T>): Result<T> {
  const rec = val as WuestenObject;
  for (const key of [param.jsonname, param.varname]) {
    if (rec[key] === undefined || rec[key] === null) {
      continue;
    }
    const my = coerce(rec[key] as C);
    return my;
  }
  if (param.default !== undefined) {
    return coerce(param.default as C);
  }
  return Result.Err(`not found:${param.jsonname}`);
}

export const WuestenTypeRegistry = new WuestenTypeRegistryImpl();

// export type WuestenObject = Record<string, unknown>;
export const WuestenObjectFactory = new WuestenObjectFactoryImpl(
  WuestenAttributeFactory<WuestenObject, WuestenObject | unknown, WuestenObject>({
    base: "WuestenObject",
    varname: "WuestenObject",
    jsonname: "WuestenObject",
  }),
);

export interface WuesteCreateBuilder<T, C, O> {
  Builder(base: WuestenAttributeBase<unknown>, params: WuestenFactoryParam<C>): WuestenBuilder<T, C, O>;
}

export function WuestenAttributeFactory<T, C, O>(
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ...params: (WuestenAttributeParameter<unknown> | undefined)[]
): WuesteCreateBuilder<T, C, O> {
  // return WuestenTypeRegistry.cloneAttributeBase(...params);
  return undefined as unknown as WuesteCreateBuilder<T, C, O>;
}

export function WuestenAttributeFactoryOptional<T, C, O>(
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ...params: (WuestenAttributeParameter<unknown> | undefined)[]
): WuesteCreateBuilder<T | undefined, C | undefined, O | undefined> {
  // return WuestenTypeRegistry.cloneAttributeBase(...params);
  return undefined as unknown as WuesteCreateBuilder<T, C, O>;
}
