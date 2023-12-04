import { Result } from "../result";
import { WuestenFormatterIf } from "../wueste";

export type Type = string;
export type CoerceType = string | boolean | number | { toString: () => string };
export type ObjectType = string;

class StringFormatterImpl implements WuestenFormatterIf<Type, CoerceType, ObjectType> {
  Coerce(value: CoerceType): Result<Type> {
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
      return this.Coerce((value as { toString: () => string }).toString());
    }
    try {
      return Result.Err("not a string: " + value);
    } catch (err) {
      return Result.Err("not a string: " + err);
    }
  }
  ToObject(value: Type): Result<ObjectType> {
    return Result.Ok(value);
  }
}

export const Formatter = new StringFormatterImpl();
