import { Result } from "./result";

type PromiseOrValue<T> = Promise<T> | T;

function scopey<T>(fn: (scope: Scope) => PromiseOrValue<Result<T>>) {
  const scope = new Scope();
  try {
    return Result.Ok(fn(scope));
  } catch (err) {
    for (const bld of scope.builders.reverse() ){
      try {
        const ret = bld.catchFn?.(err as Error);
        if (ret && typeof ret.then === 'function') {
          await ret;
        }
      } catch (err) {
        console.error(err);
      }
    }
    return Result.Err(err as Error);
  } finally {
    scope.builders.reverse().forEach((bld) => {
      try {
        bld.finallyFn?.();
      } catch (err) {
        console.error(err);
      }
    })
  }
}

type WithoutPromise<T> = T extends Promise<infer U> ? U : T;

class EvalBuilder<T extends PromiseOrValue<unknown>> {
  readonly scope: Scope;
  readonly evalFn: () => PromiseOrValue<T>;
  cleanupFn?: (t: WithoutPromise<T>) => PromiseOrValue<void>;
  catchFn?: (err: Error) => PromiseOrValue<void>;
  finallyFn?: () => PromiseOrValue<void>;

  constructor(scope: Scope, fn: () => PromiseOrValue<T>) {
    this.evalFn = fn;
    this.scope = scope;
  }

  cleanup(fn: (t: WithoutPromise<T>) => PromiseOrValue<void>): this {
    this.cleanupFn = fn;
    return this;
  }
  catch(fn: (err: Error) => PromiseOrValue<void>): this {
    this.catchFn = fn;
    return this;
  }
  finally(fn: () => PromiseOrValue<void>): this {
    this.finallyFn = fn;
    return this;
  }
  do(): PromiseOrValue<T> {
    let ctx: PromiseOrValue<T> | undefined;
    try {
      ctx = this.evalFn();
      return ctx;
    } catch (err) {
      this.scope.onCatch(() => {
        this.catchFn?.(err as Error);
      });
      throw err;
    } finally {
      this.scope.onFinally(() => {
        if (ctx) {
          this.scope.onCleanup(() => {
            this.cleanupFn?.(ctx as WithoutPromise<T>);
          });
        }
        this.finallyFn?.();
      });
    }
  }
}

class Scope {
  builders: EvalBuilder<unknown>[] = [];
  eval<T>(fn: () => PromiseOrValue<T>): EvalBuilder<T> {
    const bld = new EvalBuilder<T>(this, fn);
    this.builders.push(bld as EvalBuilder<unknown>);
    return bld;
  }

  cleanups: (() => PromiseOrValue<void>)[] = [];
  onCleanup(fn: () => PromiseOrValue<void>) {
    this.cleanups.push(fn);
  }
  catchFn?: (err: Error) => PromiseOrValue<void>;
  onCatch(fn: () => PromiseOrValue<void>) {
    this.catchFn = fn;
  }
  finallys: (() => PromiseOrValue<void>)[] = [];
  onFinally(fn: () => PromiseOrValue<void>) {
    this.finallys.push(fn);
  }
}

function throwsError() {
  throw new Error("error");
}


it("a scopey is a exception in catch and finally", async () => {
})

it("a scopey is a function", async () => {
  const first = {
    eval: jest.fn(),
    close: jest.fn(),
    cleanup: jest.fn(),
    catch: jest.fn(),
    finally: jest.fn(),
  }
  const second = {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    eval: jest.fn((nr: number) => throwsError()),
    close: jest.fn(),
    cleanup: jest.fn(),
    catch: jest.fn(),
    finally: jest.fn(),
  }
  const third = {
    eval: jest.fn(() => throwsError()),
    close: jest.fn(),
    cleanup: jest.fn(),
    catch: jest.fn(),
    finally: jest.fn(),
  }
  let cleanupOrder = 0
  let catchOrder = 0
  let finallyOrder = 0
  let evalOrder = 0
  const sc = await scopey(async (scope) => {
    const test = await scope
      .eval(() => {
        first.eval(evalOrder++);
        return {
          db: {
            close: first.close,
            update: () => {}
          },
        };
      })
      .cleanup((ctx) => {
        ctx.db.close();
        first.cleanup(cleanupOrder++)
      })
      .catch(() => first.catch(catchOrder++))
      .finally(() => first.finally(finallyOrder++))
      .do() as unknown as { db: { close: () => void, update: (o: string) => void } };
    expect(test.db.close).toEqual(first.close);
    await scope
      .eval(() => {
        second.eval(evalOrder++);
        return {
          db: {
            close: second.close,
          },
        };
      })
      .cleanup(() => {
        second.cleanup(cleanupOrder++)
        test.db.update(`update error table set error = 'error'`)
      })
      .catch(() => {
        second.catch(catchOrder++)
      })
      .finally(() => second.finally(finallyOrder++))
      .do();

      await scope
      .eval(() => {
        third.eval();
        return {
          db: {
            close: third.close,
          },
        };
      })
      .cleanup(third.cleanup)
      .catch(third.catch)
      .finally(third.finally)
      .do();
      return { wurst: 4 }
  });
  expect(sc).toBeInstanceOf(Error);
  expect(sc).toEqual({ wurst: 4 });

  expect(first.eval).toHaveBeenCalled();
  expect(first.eval.mock.calls[0][0]).toEqual(0);
  expect(second.eval).toHaveBeenCalled();
  expect(second.eval.mock.calls[0][0]).toEqual(1);
  expect(third.eval).toHaveBeenCalledTimes(0);

  expect(first.close).toHaveBeenCalledTimes(0);
  expect(first.close.mock.calls[0][0]).toEqual(1);
  expect(second.close).toHaveBeenCalledTimes(0);
  expect(second.close.mock.calls[0][0]).toEqual(0);
  expect(third.close).toHaveBeenCalledTimes(0);

  expect(first.cleanup).toHaveBeenCalledTimes(1);
  expect(first.cleanup.mock.calls[0][0]).toEqual(1);
  expect(second.cleanup).toHaveBeenCalledTimes(1);
  expect(second.cleanup.mock.calls[0][0]).toEqual(0);
  expect(third.cleanup).toHaveBeenCalledTimes(0);

  expect(first.catch).toHaveBeenCalledTimes(0);
  expect(first.catch.mock.calls[0][0]).toEqual(1);
  expect(second.catch).toHaveBeenCalledTimes(1);
  expect(second.catch.mock.calls[0][0]).toEqual(0);
  expect(third.catch).toHaveBeenCalledTimes(0);

  expect(first.finally).toHaveBeenCalledTimes(1);
  expect(first.finally.mock.calls[0][0]).toEqual(1);
  expect(second.finally).toHaveBeenCalledTimes(1);
  expect(second.finally.mock.calls[0][0]).toEqual(0);
  expect(third.finally).toHaveBeenCalledTimes(0);
});


it("a scopey happy path", async () => {
  const first = {
    eval: jest.fn(),
    close: jest.fn(),
    cleanup: jest.fn(),
    catch: jest.fn(),
    finally: jest.fn(),
  }
  const second = {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    eval: jest.fn((nr: number) => throwsError()),
    close: jest.fn(),
    cleanup: jest.fn(),
    catch: jest.fn(),
    finally: jest.fn(),
  }
  const third = {
    eval: jest.fn(() => throwsError()),
    close: jest.fn(),
    cleanup: jest.fn(),
    catch: jest.fn(),
    finally: jest.fn(),
  }
  let cleanupOrder = 0
  let catchOrder = 0
  let finallyOrder = 0
  let evalOrder = 0
  const sc = await scopey(async (scope) => {
    const test = await scope
      .eval(() => {
        first.eval(evalOrder++);
        return {
          db: {
            close: first.close,
            update: () => {}
          },
        };
      })
      .cleanup((ctx) => {
        ctx.db.close();
        first.cleanup(cleanupOrder++)
      })
      .catch(() => first.catch(catchOrder++))
      .finally(() => first.finally(finallyOrder++))
      .do() as unknown as { db: { close: () => void, update: (o: string) => void } };
    expect(test.db.close).toEqual(first.close);
    await scope
      .eval(() => {
        second.eval(evalOrder++);
        return {
          db: {
            close: second.close,
          },
        };
      })
      .cleanup(() => {
        second.cleanup(cleanupOrder++)
        test.db.update(`update error table set error = 'error'`)
      })
      .catch(() => {
        second.catch(catchOrder++)
      })
      .finally(() => second.finally(finallyOrder++))
      .do();

      await scope
      .eval(() => {
        third.eval();
        return {
          db: {
            close: third.close,
          },
        };
      })
      .cleanup(third.cleanup)
      .catch(third.catch)
      .finally(third.finally)
      .do();
      return { wurst: 4 }
  });
  expect(sc).toBeInstanceOf(Error);
  expect(sc).toEqual({ wurst: 4 });

  expect(first.eval).toHaveBeenCalled();
  expect(first.eval.mock.calls[0][0]).toEqual(0);
  expect(second.eval).toHaveBeenCalled();
  expect(second.eval.mock.calls[0][0]).toEqual(1);
  expect(third.eval).toHaveBeenCalledTimes(1);

  expect(first.close).toHaveBeenCalledTimes(0);
  expect(first.close.mock.calls[0][0]).toEqual(1);
  expect(second.close).toHaveBeenCalledTimes(0);
  expect(second.close.mock.calls[0][0]).toEqual(0);
  expect(third.close).toHaveBeenCalledTimes(1);

  expect(first.cleanup).toHaveBeenCalledTimes(1);
  expect(first.cleanup.mock.calls[0][0]).toEqual(1);
  expect(second.cleanup).toHaveBeenCalledTimes(1);
  expect(second.cleanup.mock.calls[0][0]).toEqual(0);
  expect(third.cleanup).toHaveBeenCalledTimes(0);

  expect(first.catch).toHaveBeenCalledTimes(0);
  expect(first.catch.mock.calls[0][0]).toEqual(0);
  expect(second.catch).toHaveBeenCalledTimes(0);
  expect(second.catch.mock.calls[0][0]).toEqual(0);
  expect(third.catch).toHaveBeenCalledTimes(0);

  expect(first.finally).toHaveBeenCalledTimes(1);
  expect(first.finally.mock.calls[0][0]).toEqual(1);
  expect(second.finally).toHaveBeenCalledTimes(1);
  expect(second.finally.mock.calls[0][0]).toEqual(0);
  expect(third.finally).toHaveBeenCalledTimes(0);
});