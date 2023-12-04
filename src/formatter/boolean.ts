import { Result } from "../result";
import { WuestenFormatterIf } from "../wueste";
import { Formatter as NumberFormatter } from "./number";

export type Type = boolean;
export type CoerceType = boolean | string | number | unknown;
export type ObjectType = boolean;

class BooleanFormatterImpl implements WuestenFormatterIf<Type, CoerceType, ObjectType> {
  Coerce(value: CoerceType): Result<Type> {
    if (typeof value === "boolean") {
      return Result.Ok(value);
    }
    if (typeof value === "string") {
      if (["true", "yes", "on"].includes(value.toLowerCase())) {
        return Result.Ok(true);
      }
      if (["false", "no", "off"].includes(value.toLowerCase())) {
        return Result.Ok(false);
      }
      const rnum = NumberFormatter.Coerce(value);
      if (rnum.is_ok()) {
        return this.Coerce(rnum.unwrap());
      }
    }
    return Result.Ok(!!value);
  }
  ToObject(value: Type): Result<ObjectType> {
    return Result.Ok(value);
  }
}

export const Formatter = new BooleanFormatterImpl();
