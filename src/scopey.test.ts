
type PromiseOrValue<T> = Promise<T> | T

function scopey<T>(fn: (scope: Scope) => PromiseOrValue<T>) {
    const ret = fn(new Scope())
    return ret
}


type WithoutPromise<T> = T extends Promise<infer U> ? U : T

class EvalBuilder<T extends PromiseOrValue<unknown>> {
    readonly scope: Scope
    readonly evalFn: () => PromiseOrValue<T>
    cleanupFn?: (t: WithoutPromise<T>) => PromiseOrValue<void>
    catchFn?: (err: Error) => PromiseOrValue<void>
    finallyFn?: () => PromiseOrValue<void>

    constructor(scope: Scope, fn: () => PromiseOrValue<T>) {
        this.evalFn = fn
        this.scope = scope
    }

    cleanup(fn: (t: WithoutPromise<T>) => PromiseOrValue<void>): this {
        this.cleanupFn = fn
        return this
    }
    catch(fn: (err: Error) => PromiseOrValue<void>): this {
        this.catchFn = fn
        return this
    }
    finally(fn: () => PromiseOrValue<void>): this {
        this.finallyFn = fn
        return this
    }
    do(): PromiseOrValue<void> {
        let ctx: PromiseOrValue<T> | undefined
        try {
            ctx = this.evalFn()
        } catch (err) {
            this.scope.onCatch(() => {
                this.catchFn?.(err as Error)
            })
            throw err
        } finally {
            if (ctx) {
                this.scope.onCleanup(() => {
                    this.cleanupFn?.(ctx as WithoutPromise<T>)
                })
            }
            this.scope.onFinally(() => {
                this.finallyFn?.()
            })
        }
    }
}

class Scope {
    builders: EvalBuilder<unknown>[] = []
    eval<T>(fn: () => PromiseOrValue<T>): EvalBuilder<T> {
        const bld = new EvalBuilder<T>(this, fn)
        this.builders.push(bld as EvalBuilder<unknown>)
        return bld
    }

    cleanups: (() => PromiseOrValue<void>)[] = []
    onCleanup(fn: () => PromiseOrValue<void>) {
        this.cleanups.push(fn)
    }
    catchFn?: (err: Error) => PromiseOrValue<void>
    onCatch(fn: () => PromiseOrValue<void>) {
        this.catchFn = fn
    }
    finallys: (() => PromiseOrValue<void>)[]= []
    onFinally(fn: () => PromiseOrValue<void>) {
        this.finallys.push(fn)
    }
}

function throwsError() {
    throw new Error("error")
}


it('a scopey is a function', () => {
    scopey(async (scope) => {
        await scope.eval(() => {
            throwsError()
            return {
                db: {
                    close: () => { }
                }
            }
        }).cleanup((ctx) => {
            ctx.db.close()
        }).catch((err) => {
            console.log(err)
        }).finally(() => {
        }).do()
    })
})

