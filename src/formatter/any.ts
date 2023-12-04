import { Payload } from "../payload";
import { Result } from "../result";
import { WuestenAttributeBase, WuestenAttributeParameter, WuestenBuilder, WuestenEncoder, WuestenFormatterIf } from "../wueste";

export type Type = Record<string, unknown>;
export type CoerceType = Record<string, unknown>;
export type ObjectType = Record<string, unknown>;

export class Builder implements WuestenBuilder<Type, CoerceType, ObjectType> {
  readonly param: WuestenAttributeBase<CoerceType>;
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  constructor(param: WuestenAttributeBase<CoerceType>, ...params: WuestenAttributeParameter<CoerceType>[]) {
    this.param = param;
  }
  Get(): Result<Type, Error> {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  CoerceAttribute(val: unknown): Result<Type, Error> {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Coerce(value?: CoerceType): WuestenBuilder<Type, CoerceType, ObjectType> {
    throw new Error("Method not implemented.");
  }
  ToObject(): Result<ObjectType, Error> {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ToPayload(encoder: WuestenEncoder = this.param.encoder): Result<Payload, Error> {
    throw new Error("Method not implemented.");
  }
}

class AnyFormatterImpl implements WuestenFormatterIf<Type, CoerceType, ObjectType> {
  Coerce(value: CoerceType): Result<Type> {
    throw new Error("Method not implemented.");
    return Result.Ok(value);
    // if (typeof value === "string") {
    //     return Result.Ok(value);
    //   }
    //   if (typeof value === "number") {
    //     return Result.Ok("" + value);
    //   }
    //   if (typeof value === "boolean") {
    //     return Result.Ok(value ? "true" : "false");
    //   }
    //   if (
    //     (typeof value === "object" || typeof value === "function") &&
    //     value !== null &&
    //     typeof (value as { toString: () => string })["toString"] === "function"
    //   ) {
    //     return this.Coerce((value as { toString: () => string }).toString());
    //   }
    //   try {
    //     return Result.Err("not a string: " + value);
    //   } catch (err) {
    //     return Result.Err("not a string: " + err);
    //   }
  }
  ToObject(value: Type): Result<ObjectType> {
    return Result.Ok(value);
  }
}

export const Formatter = new AnyFormatterImpl();
