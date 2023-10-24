export class Result<T, E = Error> {
  static Ok<T>(t: T): Result<T, Error> {
    return new ResultOK(t);
  }
  static Err<T extends Error = Error>(t: T | string): Result<never, T> {
    if (typeof t === "string") {
      return new ResultError(new Error(t) as T);
    }
    return new ResultError(t);
  }
  static Is<T>(t: unknown): t is Result<T> {
    return t instanceof ResultOK || t instanceof ResultError;
  }

  isOk(): boolean {
    return this.is_ok();
  }
  isErr(): boolean {
    return this.is_ok();
  }

  Ok(): T {
    return this.unwrap();
  }
  Err(): E {
    return this.unwrap_err();
  }

  is_ok(): boolean {
    throw new Error("Not implemented");
  }
  is_err(): boolean {
    throw new Error("Not implemented");
  }
  unwrap(): T {
    throw new Error("Not implemented");
  }
  unwrap_err(): E {
    throw new Error("Not implemented");
  }
}

export class ResultOK<T> extends Result<T, Error> {
  private _t: T;
  constructor(t: T) {
    super();
    this._t = t;
  }
  is_ok(): boolean {
    return true;
  }
  is_err(): boolean {
    return false;
  }
  unwrap_err(): Error {
    throw new Error("Result is Ok");
  }
  unwrap(): T {
    return this._t;
  }
}

export class ResultError<T extends Error> extends Result<never, T> {
  private _error: T;
  constructor(t: T) {
    super();
    this._error = t;
  }
  is_ok(): boolean {
    return false;
  }
  is_err(): boolean {
    return true;
  }
  unwrap(): never {
    throw new Error(`Result is Err: ${this._error}`);
  }
  unwrap_err(): T {
    return this._error;
  }
}

export function IsResult<T>(t: unknown): t is Result<T> {
  return t instanceof ResultOK || t instanceof ResultError;
}
