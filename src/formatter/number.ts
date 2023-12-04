import { Result } from "../result";
import { WuestenFormatterIf } from "../wueste";

export type Type = number;
export type CoerceType = string | number;
export type ObjectType = number;

class NumberFormatterImpl implements WuestenFormatterIf<Type, CoerceType, ObjectType> {
  Coerce(value: CoerceType): Result<Type> {
    const val = parseFloat(value as string);
    if (isNaN(val)) {
      return Result.Err(`not a number: ${value}`);
    }
    return Result.Ok(val);
  }
  ToObject(value: Type): Result<ObjectType> {
    return Result.Ok(value);
  }
}

export const Formatter = new NumberFormatterImpl();
