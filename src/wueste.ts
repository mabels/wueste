import { Result } from "./result";

export interface WuestenAttributeParameter<T> {
    readonly base: string
    readonly varname: string
    readonly jsonname: string
    default?: T

    // setError?: (err : string | Error) => void;
    // format?: string // date-time
}

export interface WuestenAttribute<G, I = G> {
    readonly param: WuestenAttributeParameter<G>
    SetNameSuffix(...idxs: number[]): void
    CoerceAttribute(val: unknown): Result<G>
    Coerce(value: I): Result<G>;
    Get(): Result<G>
}

function coerceAttribute<T, I>(val: unknown, param: WuestenAttributeParameter<T>, coerce: (t: I) => Result<T>): Result<T> {
    const rec = val as Record<string, unknown>
    for (const key of [param.jsonname, param.varname]) {
        if (rec[key] === undefined || rec[key] === null) {
            continue
        }
        const my = coerce(rec[key] as I)
        return my
    }
    if (param.default !== undefined) {
        return coerce(param.default as I)
    }
    return Result.Err(`not found:${param.jsonname}`)
}

export function WuestenAttributeName<T>(param: WuestenAttributeParameter<T>): string {
    const names = []
    if (param.base) {
        names.push(param.base)
    }
    names.push(param.jsonname)
    return names.join(".")
}

export class WuestenAttr<G, I = G> implements WuestenAttribute<G, I> {
    _value?: G
    _idxs: number[] = []
    readonly param: WuestenAttributeParameter<G>
    readonly _coerce: (t: I) => Result<G>
    constructor(param: WuestenAttributeParameter<I>, coerce: (t: I) => Result<G>) {
        let def: G | undefined = undefined
        this._coerce = coerce
        const result = coerce(param.default as I)
        if (result.is_ok()) {
            def = result.unwrap() as G
        }
        this.param = {
            ...param,
            default: def
        }
    }
    SetNameSuffix(...idxs: number[]): void {
        this._idxs = idxs
    }
    CoerceAttribute(val: unknown): Result<G> {
        if (!(typeof val === "object" && val !== null)) {
            return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is not an object:` + val)
        }
        const res = coerceAttribute<G, I>(val, this.param, this.Coerce.bind(this))
        if (res.is_err()) {
            return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] ${res.unwrap_err().message}`)
        }
        return res
    }
    Coerce(value: I): Result<G> {
        const result = this._coerce(value)
        if (result.is_ok()) {
            this._value = result.unwrap()
            return result
        }
        return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is ${result.unwrap_err().message}`)
    }
    Get(): Result<G> {
        if (this.param.default === undefined && this._value === undefined) {
            return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is required`)
        }
        if (this._value !== undefined) {
            return Result.Ok(this._value)
        }
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        return Result.Ok(this.param.default! as G)
    }
}

export class WuestenAttrOptional<G, T = G | undefined, I=G> implements WuestenAttribute<T, I> {
    readonly _attr: WuestenAttribute<T, I>
    readonly param: WuestenAttributeParameter<T>
    _value: T
    _idxs: number[] = []

    constructor(attr: WuestenAttribute<T, I>) {
        this._attr = attr
        this.param = {
            ...attr.param,
            default: attr.param.default as T
        }
        this._value = attr.param.default as T
    }
    SetNameSuffix(...idxs: number[]): void {
        this._idxs = idxs
    }
    CoerceAttribute(val: unknown): Result<T> {
        if (!(typeof val === "object" && val !== null)) {
            return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is not an object:` + val)
        }
        const res = coerceAttribute(val, this.param, this.Coerce.bind(this))
        if (res.is_ok()) {
            this._value = res.unwrap() as T
            return res;
        }
        return Result.Ok(this.param.default as T)
    }

    Coerce(value: I): Result<T> {
        if (value === undefined || value === null) {
            this._value = undefined as T
            return Result.Ok(this._value)
        }
        const res = this._attr.Coerce(value)
        if (res.is_ok()) {
            this._value = res.unwrap() as unknown as T
            return Result.Ok(this._value)
        }
        return Result.Err(res.unwrap_err())
    }
    Get(): Result<T> {
        return Result.Ok(this._value)
    }
}

export interface WuestenFactory<B extends WuestenAttribute<T, I>, T, I = T> {
    Builder(param?: WuestenAttributeParameter<I>): B;
    // Coerce FromObject(object: unknown): Result<T>;
    ToObject(typ: T): unknown;
    Clone(typ: T): Result<T>;
}

function stringCoerce(value: unknown): Result<string> {
    if (typeof value === "string") {
        return Result.Ok(value)
    }
    if (typeof value === "number") {
        return Result.Ok('' + value)
    }
    if (typeof value === "boolean") {
        return Result.Ok(value ? "true" : "false")
    }
    if (((typeof value === "object" || typeof value === "function") && value !== null) &&
        typeof (value as { toString: () => string })['toString'] === "function") {
        return stringCoerce((value as { toString: () => string }).toString())
    }
    try {
        return Result.Err("not a string: " + value)
    } catch (err) {
        return Result.Err("not a string: " + err)
    }
}



function dateTimeCoerce(value: unknown): Result<Date> {
    if (typeof value === "string") {
        return Result.Ok(new Date(value))
    }
    if (typeof value === "number") {
        return Result.Ok(new Date(value))
    }
    if (value instanceof Date) {
        return Result.Ok(value)
    }
    return Result.Err("not a Date: " + value)
}


function booleanCoerce(value: unknown): Result<boolean> {
    if (typeof value === "boolean") {
        return Result.Ok(value)
    }
    if (typeof value === "string") {
        if (["true", "1", "yes", "on"].includes(value.toLowerCase())) {
            return Result.Ok(true)
        }
        if (["false", "0", "no", "off"].includes(value.toLowerCase())) {
            return Result.Ok(false)
        }
    }
    if (typeof value === "number") {
        return Result.Ok(!!value)
    }
    return Result.Err("not a boolean: " + value)
}

export class WuestenAttributeObject<B extends WuestenAttribute<T, I>, T, I> extends WuestenAttr<T, I> {
    private readonly _builder: WuestenAttribute<T, I>;
    constructor(param: WuestenAttributeParameter<I>, factory: WuestenFactory<B, T, I>) {
        const builder = factory.Builder(param)
        super(param, builder.Coerce.bind(builder))
        this._builder = builder
    }

    Coerce(value: I): Result<T> {
        const res = this._builder.Coerce(value)
        if (res.is_ok()) {
            this._value = res.unwrap()
        }
        return res
    }

    Get(): Result<T> {
        if (this.param.default === undefined && this._value === undefined) {
            return Result.Err(`Attribute[${WuestenAttributeName(this.param)}] is required`)
        }
        if (this._value !== undefined) {
            return Result.Ok(this._value)
        }
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        return Result.Ok(this.param.default!)
    }
}

function numberCoerce(parse: (i: unknown) => number): (value: unknown) => Result<number> {
    return (value: unknown): Result<number> => {
        const val = parse(value as string)
        if (isNaN(val)) {
            return Result.Err(`not a number: ${value}`)
        }
        return Result.Ok(val)
    }
}

export function WuesteIterable<T>(obj: unknown): IterableIterator<T> | undefined {
    if (Array.isArray(obj)) {
        const range = {
            [Symbol.iterator]() { // (1)
                return {
                    current: 0,
                    next() { // (2)
                        const aobj = obj as ArrayLike<T>
                        if (this.current < aobj.length) {
                            return { done: false, value: aobj[this.current++] };
                        } else {
                            return { done: true };
                        }
                    }
                };
            }
        };
        return range as unknown as IterableIterator<T>
    }
    if (typeof obj === "function") {
        obj = obj()
    }
    if (typeof obj === "object") {
        if (Symbol.iterator in (obj as { [Symbol.iterator]: unknown })) {
            // const iter = (obj as unknown as Iterable<unknown>)[Symbol.iterator]()
            // const range = {
            //     [Symbol.iterator]() { // (1)
            //         return {
            //             current: 0,
            //             next() { // (2)
            //                 return iter.next()
            //             }
            //         };
            //     }
            // };
            return obj as unknown as IterableIterator<T>
        }
        // if (Symbol.asyncIterator in (obj as { [Symbol.asyncIterator]: unknown })) {
        //     return obj as unknown as AsyncIterableIterator<unknown>
        // }
        if (obj !== null) {
            const vobj = Object.values(obj)
            const range = {
                [Symbol.iterator]() { // (1)
                    return {
                        current: 0,
                        next() { // (2)
                            if (this.current < vobj.length) {
                                return { done: false, value: vobj[this.current++] };
                            } else {
                                return { done: true };
                            }
                        }
                    };
                }
            };
            return range as unknown as IterableIterator<T>
        }
    }
    return undefined
}


export type WuesteCoerceTypeDate = Date|string
export type WuesteCoerceTypeboolean = boolean|string|number
export type WuesteCoerceTypenumber = number|string
export type WuesteCoerceTypestring = string|boolean|number|{toString: () => string}

export const wuesten = {
    AttributeString: (def: WuestenAttributeParameter<string>): WuestenAttribute<string, WuesteCoerceTypestring> => {
        return new WuestenAttr(def, stringCoerce);
    },
    AttributeStringOptional: (def: WuestenAttributeParameter<string>): WuestenAttribute<string | undefined, WuesteCoerceTypestring|undefined> => {
        return new WuestenAttrOptional(new WuestenAttr(def, stringCoerce));
    },

    AttributeDateTime: (def: WuestenAttributeParameter<Date | string>): WuestenAttribute<Date, WuesteCoerceTypeDate> => {
        return new WuestenAttr(def, dateTimeCoerce);
    },
    AttributeDateTimeOptional: (def: WuestenAttributeParameter<Date | string>): WuestenAttribute<Date | undefined, WuesteCoerceTypeDate|undefined> => {
        return new WuestenAttrOptional(new WuestenAttr(def, dateTimeCoerce));
    },

    AttributeInteger: (def: WuestenAttributeParameter<number>): WuestenAttribute<number, WuesteCoerceTypenumber> => {
        return new WuestenAttr(def, numberCoerce((a) => parseInt(a as string, 10)));
    },
    AttributeIntegerOptional: (def: WuestenAttributeParameter<number>): WuestenAttribute<number | undefined, WuesteCoerceTypenumber|undefined> => {
        return new WuestenAttrOptional(new WuestenAttr(def, numberCoerce((a) => parseInt(a as string, 10))));
    },

    AttributeNumber: (def: WuestenAttributeParameter<number>): WuestenAttribute<number, WuesteCoerceTypenumber> => {
        return new WuestenAttr(def, numberCoerce((a) => parseFloat(a as string)));
    },
    AttributeNumberOptional: (def: WuestenAttributeParameter<number>): WuestenAttribute<number | undefined, WuesteCoerceTypenumber|undefined> => {
        return new WuestenAttrOptional(new WuestenAttr(def, numberCoerce((a) => parseFloat(a as string))));
    },

    AttributeBoolean: (def: WuestenAttributeParameter<boolean>): WuestenAttribute<boolean, WuesteCoerceTypeboolean> => {
        return new WuestenAttr(def, booleanCoerce);
    },
    AttributeBooleanOptional: (def: WuestenAttributeParameter<boolean>): WuestenAttribute<boolean | undefined, WuesteCoerceTypeboolean|undefined> => {
        return new WuestenAttrOptional(new WuestenAttr(def, booleanCoerce));
    },

    AttributeObject: <A extends WuestenAttribute<E, I>, E, I>(def: WuestenAttributeParameter<I>, factory: WuestenFactory<A, E, I>): WuestenAttribute<E, I> => {
        return new WuestenAttributeObject<A, E, I>(def, factory);
    },
    AttributeObjectOptional: <A extends WuestenAttribute<E, I>, E, I>(def: WuestenAttributeParameter<I>, factory: WuestenFactory<A, E, I>): WuestenAttribute<E | undefined, I | undefined> => {
        return new WuestenAttrOptional<A, E, I>(new WuestenAttributeObject<A, E, I>(def, factory));
    },

    //   AttributeArray: <T>(): WuestenAttribute<T> => {
    //     return new WuestenAttributeType([] as unknown as T);
    //   },

    //   AttributeArrayOptional: <T>(): WuestenAttribute<T|undefined> => {
    //     return new WuestenAttributeType([] as unknown as T);
    //   }

}
