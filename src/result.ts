export class ResultOK<T> implements Result<T, Error> {
  private _t: T;
  constructor(t: T) {
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

export class ResultError<T extends Error> implements Result<never, T> {
  private _error: T;
  constructor(t: T) {
    this._error = t;
  }
  is_ok(): boolean {
    return false;
  }
  is_err(): boolean {
    return true;
  }
  unwrap(): never {
    throw new Error("Result is Err");
  }
  unwrap_err(): T {
    return this._error;
  }
}

export abstract class Result<T, E = Error> {
  static Ok<T>(t: T): Result<T, Error> {
    return new ResultOK(t);
  }
  static Err<T extends Error = Error>(t: T | string): Result<never, T> {
    if (typeof t === "string") {
      return new ResultError(new Error(t) as T);
    }
    return new ResultError(t);
  }
  abstract is_ok(): boolean;
  abstract is_err(): boolean;
  // abstract err(): E;
  // abstract unwrap(): T;
  abstract unwrap(): T;
  abstract unwrap_err(): E;
}
