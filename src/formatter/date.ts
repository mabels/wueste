import { Result } from "../result";
import { WuestenFormatterIf } from "../wueste";

export type Type = Date;
export type CoerceType = Date | string | number;
export type ObjectType = string;

class DateFormatterImpl implements WuestenFormatterIf<Type, CoerceType, ObjectType> {
  ToObject(value: Date): Result<string> {
    return Result.Ok(value.toISOString());
  }
  Coerce(value: CoerceType): Result<Type> {
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
}

export const Formatter = new DateFormatterImpl();
