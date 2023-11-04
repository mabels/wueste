/* eslint-disable @typescript-eslint/no-explicit-any */
import { Result } from "./result";
import { Payload } from "./payload";

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

export interface WuestenAttributeFactory {
  readonly formatsHandler: WuestenFormatsHandler;
  readonly encoder: (payload: unknown) => Result<unknown>;
  readonly decoder: (payload: unknown) => Result<unknown>;
}

export interface WuestenAttributeParam<T, I> {
  readonly base: string;
  readonly varname: string;
  readonly jsonname: string;
  default?: T | I;
  format?: string;
}

export type WuestenAttributeParameter<T, I> = WuestenAttributeParam<T, I> & Partial<WuestenAttributeFactory>;

export type WuestenAttributeBase<T, I> = WuestenAttributeParam<T, I> & WuestenAttributeFactory;

export type WuestenFactoryParam<T, I> = Partial<WuestenAttributeParameter<T, I>>;

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

export interface WuestenAttribute<G, I = G> {
  readonly param: WuestenAttributeBase<G, I>;
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

export type WuestenGetterFn = (level: WuestenReflection[], value: unknown) => void;

export function WuestenRecordGetter(
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  fn: WuestenGetterFn,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  level: WuestenReflection[],
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  v: unknown,
) {
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

export interface WuestenGeneratorFunctions<G, I> {
  readonly coerce: (t: I) => Result<G>;
  // readonly reflection?: WuestenReflection;
}

function coerceAttribute<T, I>(val: unknown, param: WuestenAttributeParameter<T, I>, coerce: (t: I) => Result<T>): Result<T> {
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

export function WuestenAttributeName<T, I>(param: WuestenAttributeParameter<T, I>): string {
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
  readonly param: WuestenAttributeBase<G, I>;
  readonly _fnParams: WuestenGeneratorFunctions<G, I>;
  constructor(param: WuestenAttributeParameter<G, I>, fnParams: WuestenGeneratorFunctions<G, I>) {
    this.param = WuestenFactoryAttributeMerge(param);
    const formatHandler = this.param.formatsHandler.get(param.format);
    if (formatHandler) {
      this._fnParams = {
        coerce: (t: I) => {
          return fnParams.coerce(formatHandler(t) as unknown as I);
        },
      };
    } else {
      this._fnParams = fnParams;
    }
    let def: G | undefined = undefined;
    const result = this._fnParams.coerce(param.default as I);
    if (result.is_ok()) {
      def = result.unwrap() as G;
    }
    this.param = { ...this.param, default: def };
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
  readonly param: WuestenAttributeBase<T, C>;
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
  readonly param: WuestenAttributeBase<T | undefined, I>;
  _value: T;
  _idxs: number[] = [];

  constructor(attr: WuestenAttribute<T | undefined, I | undefined>) {
    this._attr = attr;
    this.param = WuestenFactoryAttributeMerge(attr.param);
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

export interface WuestenBuilder<T, I> extends WuestenAttribute<T, I> {
  Get(): Result<T>;
}

export interface WuestenNames {
  readonly id: string;
  readonly title: string;
  readonly names: string[];
  readonly varname: string;
}

export abstract class WuestenFactory<T, I, O> {
  readonly _params: WuestenAttributeBase<T, I>;
  constructor(param: WuestenAttributeParameter<T, I>) {
    this._params = WuestenFactoryAttributeMerge(param);
  }
  abstract Names(): WuestenNames;
  abstract Builder(param?: WuestenFactoryParam<T, I>): WuestenBuilder<T, I>;
  abstract FromPayload(val: Payload, decoder?: WuestenDecoder): Result<T>;
  abstract ToPayload(typ: T, encoder?: WuestenEncoder): Result<Payload>;
  abstract ToObject(typ: T): O; // WuestenObject; keys are json notation
  abstract Clone(typ: T): Result<T>;
  abstract Schema(): WuestenReflection;
  abstract Getter(typ: T, base: WuestenReflection[]): WuestenGetterBuilder;
  AddFormat(name: string, fn: (recv: unknown) => unknown): WuestenFactory<T, I, O> {
    this._params.formatsHandler.add(name, fn);
    return this;
  }
}

export type WuestenObject = Record<string, unknown>;

export type WuestenFNGetBuilder<T> = (b: T | undefined) => unknown;
export class WuestenObjectBuilder implements WuestenBuilder<WuestenObject, WuestenObject> {
  readonly param: WuestenAttributeBase<WuestenObject, WuestenObject>;
  constructor(param?: WuestenFactoryParam<WuestenObject, WuestenObject>) {
    this.param = WuestenFactoryAttributeMerge({
      base: "WuestenObjectBuilder",
      varname: "WuestenObjectBuilder",
      jsonname: "WuestenObjectBuilder",
      ...param,
    });
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

export class WuestenObjectFactoryImpl extends WuestenFactory<WuestenObject, WuestenObject, WuestenObject> {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Builder(
    param?: WuestenAttributeParameter<WuestenObject, WuestenObject> | undefined,
  ): WuestenBuilder<WuestenObject, WuestenObject> {
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
  Getter(typ: WuestenObject, base: WuestenReflection[]): WuestenGetterBuilder {
    throw new Error("WuestenObjectFactoryImpl:Getter not implemented.");
  }
}

function stringCoerce(value?: unknown): Result<string> {
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
    const rnum = numberCoerce((a) => parseFloat(a as string))(value);
    if (rnum.is_ok()) {
      if (isNaN(rnum.unwrap()) || !rnum.unwrap()) {
        return Result.Ok(false);
      }
      return Result.Ok(true);
    }
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
  constructor(param: WuestenAttributeParameter<T, I>, factory: WuestenFactory<T, I, O>) {
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
    return Result.Ok(this.param.default! as T);
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
  AttributeString: (
    def: WuestenAttributeParameter<string, WuesteCoerceTypestring>,
  ): WuestenAttribute<string, WuesteCoerceTypestring> => {
    return new WuestenAttr(def, { coerce: stringCoerce });
  },
  AttributeStringOptional: (
    def: WuestenAttributeParameter<string, WuesteCoerceTypestring>,
  ): WuestenAttribute<string | undefined, WuesteCoerceTypestring | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: stringCoerce }));
  },

  AttributeDateTime: (
    def: WuestenAttributeParameter<WuesteCoerceTypeDate, WuesteCoerceTypeDate>,
  ): WuestenAttribute<Date, WuesteCoerceTypeDate> => {
    return new WuestenAttr(def, { coerce: dateTimeCoerce }) as WuestenAttribute<Date, WuesteCoerceTypeDate>;
  },
  AttributeDateTimeOptional: (
    def: WuestenAttributeParameter<WuesteCoerceTypeDate | undefined, WuesteCoerceTypeDate | undefined>,
  ): WuestenAttribute<Date | undefined, WuesteCoerceTypeDate | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: dateTimeCoerce })) as WuestenAttribute<
      Date | undefined,
      WuesteCoerceTypeDate | undefined
    >;
  },

  AttributeInteger: (
    def: WuestenAttributeParameter<number, WuesteCoerceTypenumber>,
  ): WuestenAttribute<number, WuesteCoerceTypenumber> => {
    return new WuestenAttr(def, { coerce: numberCoerce((a) => parseInt(a as string, 10)) });
  },
  AttributeIntegerOptional: (
    def: WuestenAttributeParameter<number, WuesteCoerceTypenumber>,
  ): WuestenAttribute<number | undefined, WuesteCoerceTypenumber | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: numberCoerce((a) => parseInt(a as string, 10)) }));
  },

  AttributeNumber: (
    def: WuestenAttributeParameter<number, WuesteCoerceTypenumber>,
  ): WuestenAttribute<number, WuesteCoerceTypenumber> => {
    return new WuestenAttr(def, { coerce: numberCoerce((a) => parseFloat(a as string)) });
  },
  AttributeNumberOptional: (
    def: WuestenAttributeParameter<number, WuesteCoerceTypenumber>,
  ): WuestenAttribute<number | undefined, WuesteCoerceTypenumber | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: numberCoerce((a) => parseFloat(a as string)) }));
  },

  AttributeBoolean: (
    def: WuestenAttributeParameter<boolean, WuesteCoerceTypeboolean>,
  ): WuestenAttribute<boolean, WuesteCoerceTypeboolean> => {
    return new WuestenAttr(def, { coerce: booleanCoerce });
  },
  AttributeBooleanOptional: (
    def: WuestenAttributeParameter<boolean, WuesteCoerceTypeboolean>,
  ): WuestenAttribute<boolean | undefined, WuesteCoerceTypeboolean | undefined> => {
    return new WuestenAttrOptional(new WuestenAttr(def, { coerce: booleanCoerce }));
  },

  AttributeObject: <E, I, O>(def: WuestenAttributeParameter<E, I>, factory: WuestenFactory<E, I, O>): WuestenAttribute<E, I> => {
    return new WuestenAttributeObject<E, I, O>(def, factory);
  },

  AttributeObjectOptional: <E, I, O>(
    def: WuestenAttributeParameter<E, I>,
    factory: WuestenFactory<E, I, O>,
  ): WuestenAttribute<E | undefined, I | undefined> => {
    return new WuestenAttrOptional<E, I>(new WuestenAttributeObject<E, I, O>(def, factory));
  },
};

export class WuestenTypeRegistryImpl {
  readonly _registry: Record<string, WuestenFactory<unknown, unknown, unknown>> = {};
  readonly _attributeBase: WuestenAttributeFactory = {
    formatsHandler: new WuestenFormatsHandler(),
    encoder: WuestenJSONPassThroughEncoder,
    decoder: WuestenJSONPassThroughDecoder,
  };

  Register<T extends WuestenFactory<A, B, C>, A, B, C>(factory: T): T {
    for (const name of factory.Names().names) {
      this._registry[name] = factory;
    }
    return factory;
  }

  AddFormat(name: string, fn: (recv: unknown) => unknown) {
    this._attributeBase.formatsHandler.add(name, fn);
    return this;
  }

  _deleteUndefined(merge?: Partial<WuestenAttributeFactory>): Partial<WuestenAttributeFactory> {
    const m = { ...merge };
    if (!m.decoder) {
      delete m.decoder;
    }
    if (!m.encoder) {
      delete m.encoder;
    }
    if (!m.formatsHandler) {
      delete m.formatsHandler;
    }
    return m;
  }

  cloneAttributeBase<T, I>(...merges: (WuestenAttributeParameter<unknown, unknown> | undefined)[]): WuestenAttributeBase<T, I> {
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
    } as WuestenAttributeBase<T, I>;
  }
}

export const WuestenTypeRegistry = new WuestenTypeRegistryImpl();

export const WuestenObjectFactory = new WuestenObjectFactoryImpl(
  WuestenFactoryAttributeMerge({
    base: "WuestenObject",
    varname: "WuestenObject",
    jsonname: "WuestenObject",
  }),
);

export function WuestenFactoryAttributeMerge<T, I>(
  ...params: (WuestenAttributeParameter<unknown, unknown> | undefined)[]
): WuestenAttributeBase<T, I> {
  return WuestenTypeRegistry.cloneAttributeBase(...params);
}
