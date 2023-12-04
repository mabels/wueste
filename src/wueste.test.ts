import { Result } from "./result";
import {
  WuestenAttributeParameter,
  WuestePayload,
  WuestenFactory,
  WuestenBuilder,
  WuestenDecoder,
  WuestenEncoder,
  WuestenReflection,
  WuestenRecordGetter,
  WuestenGetterBuilder,
  WuesteToIterator,
  WuestenAttributeBase,
  WuestenNames,
  WuestenAttributeFactory,
  WuestenAttributeFactoryOptional,
  WuestenAttribute,
} from "./wueste";
import { Formatter } from "./formatter";

it("array coerce from array", () => {
  const ri = WuesteToIterator<number>([1, 2, 3]);
  expect(ri.is_ok()).toBeTruthy();
  let idx = 0;
  let i = ri.Ok().next();
  for (; !i.done; i = ri.Ok().next()) {
    expect(i.idx).toBe(idx);
    expect(i.idx).toBe(i.value - 1);
    idx++;
  }
  expect(i.idx).toBe(3);
  expect(idx).toBe(3);
});

it("array coerce from empty array", () => {
  const ri = WuesteToIterator<number>([]);
  expect(ri.is_ok()).toBeTruthy();
  let idx = 0;
  let i = ri.Ok().next();
  for (; !i.done; i = ri.Ok().next()) {
    idx++;
  }
  expect(i.idx).toBe(0);
  expect(idx).toBe(0);
});

it("array coerce from object", () => {
  const ri = WuesteToIterator<number>({ a: 1, b: 2, c: 3 });
  expect(ri.is_ok()).toBeTruthy();
  let idx = 0;
  let i = ri.Ok().next();
  for (; !i.done; i = ri.Ok().next()) {
    expect(i.idx).toBe(idx);
    expect(i.idx).toBe(i.value - 1);
    idx++;
  }
  expect(i.idx).toBe(3);
  expect(idx).toBe(3);
});

it("array coerce from empty object", () => {
  const ri = WuesteToIterator<number>([]);
  expect(ri.is_ok()).toBeTruthy();
  let idx = 0;
  let i = ri.Ok().next();
  for (; !i.done; i = ri.Ok().next()) {
    idx++;
  }
  expect(i.idx).toBe(0);
  expect(idx).toBe(0);
});

it("array coerce from generator", () => {
  const ri = WuesteToIterator<number>(function* () {
    yield 1;
    yield 2;
    yield 3;
  });
  expect(ri.is_ok()).toBeTruthy();
  let idx = 0;
  let i = ri.Ok().next();
  for (; !i.done; i = ri.Ok().next()) {
    expect(i.idx).toBe(idx);
    expect(i.idx).toBe(i.value - 1);
    idx++;
  }
  expect(i.idx).toBe(3);
  expect(idx).toBe(3);
});

it("array coerce from generator array", () => {
  const ri = WuesteToIterator<number>(function* () {});
  expect(ri.is_ok()).toBeTruthy();
  let idx = 0;
  let i = ri.Ok().next();
  for (; !i.done; i = ri.Ok().next()) {
    idx++;
  }
  expect(i.idx).toBe(0);
  expect(idx).toBe(0);
});

function attrParam<T, C>(def?: WuestenAttributeParameter<C>): WuestenAttribute<T, C> {
  return WuestenAttributeFactory({
    jsonname: "x",
    varname: "X",
    base: "base",
    ...def,
  });
}

describe("string coerce", () => {
  it("string attribute", () => {
    const coerce = WuestenAttributeFactory(attrParam());
    expect(coerce.Get().unwrap_err().message).toContain("Attribute[base.x] is required");
    expect(coerce.CoerceAttribute({}).unwrap_err().message).toContain("Attribute[base.x] not found");
    expect(coerce.CoerceAttribute({ x: 4 }).unwrap()).toBe("4");
    expect(coerce.Get().unwrap()).toBe("4");
    expect(coerce.CoerceAttribute({ X: 9 }).unwrap()).toBe("9");
    expect(coerce.Get().unwrap()).toBe("9");
  });
  it("string optional attribute", () => {
    const coerce = WuestenAttributeFactoryOptional(attrParam());
    expect(coerce.Get().is_ok()).toBeTruthy();
    expect(coerce.CoerceAttribute({}).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.CoerceAttribute({ x: 4 }).unwrap()).toBe("4");
    expect(coerce.Get().unwrap()).toBe("4");
    expect(coerce.CoerceAttribute({ X: undefined }).unwrap()).toBeUndefined();
    expect(coerce.Get().unwrap()).toBe("4");
  });

  it("string no default", () => {
    const coerce = WuestenAttributeFactory(attrParam());
    expect(coerce.Get().is_err()).toBeTruthy();
    expect(coerce.Coerce({}).unwrap()).toBe("[object Object]");
    expect(coerce.Get().unwrap()).toBe("[object Object]");
    expect(coerce.Coerce(6.4).unwrap()).toBe("6.4");
    expect(coerce.Get().unwrap()).toBe("6.4");
    expect(coerce.Coerce(false).unwrap()).toBe("false");
    expect(coerce.Get().unwrap()).toBe("false");
    expect(coerce.Coerce({ toString: 5 } as unknown as string).unwrap_err().message).toContain("Attribute[base.x] is not a string");
    expect(coerce.Coerce(null as unknown as string).is_err()).toBeTruthy();
    expect(coerce.Coerce(undefined as unknown as string).is_err()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBe("false");
    expect(coerce.Coerce(NaN).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBe("NaN");
    expect(coerce.CoerceAttribute({}).unwrap_err().message).toBe("Attribute[base.x] not found:x");
  });

  it("string default", () => {
    const coerce = wuesten.AttributeString(attrParam({ default: "x" }));
    expect(coerce.Get().unwrap()).toBe("x");
    expect(coerce.Coerce("y").is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBe("y");
    expect(coerce.CoerceAttribute({}).unwrap()).toEqual("x");
    expect(coerce.CoerceAttribute({ x: "z" }).unwrap()).toBe("z");
  });

  it("stringoptional no default", () => {
    const coerce = wuesten.AttributeStringOptional(attrParam());
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(() => coerce.Coerce(null as unknown as string)).not.toThrowError();
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(() => coerce.Coerce(6.4)).not.toThrowError();
    expect(coerce.Get().unwrap()).toBe("6.4");
    expect(() => coerce.Coerce(false)).not.toThrowError();
    expect(coerce.Get().unwrap()).toBe("false");
    expect(() => coerce.Coerce({ toString: 5 } as unknown as string)).not.toThrowError();
    expect(coerce.Get().unwrap()).toBe("false");

    expect(() => coerce.Coerce(false)).not.toThrowError();
    expect(() => coerce.Coerce(null as unknown as string)).not.toThrowError();
    expect(coerce.Get().unwrap()).toBeUndefined();

    expect(() => coerce.Coerce(false)).not.toThrowError();
    expect(() => coerce.Coerce(undefined)).not.toThrowError();
    expect(coerce.Get().unwrap()).toBeUndefined();

    expect(() => coerce.Coerce(NaN)).not.toThrowError();
    expect(coerce.Get().unwrap()).toBe("NaN");
  });

  it("stringoptional default", () => {
    const coerce = wuesten.AttributeStringOptional(attrParam({ default: "x" }));
    expect(coerce.Get().unwrap()).toBe("x");
    expect(coerce.Coerce("y").is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBe("y");
    expect(coerce.Coerce(null as unknown as string).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBe(undefined);
  });
});

describe("datetime coerce", () => {
  it("datetime no default", () => {
    const coerce = wuesten.AttributeDateTime(attrParam());
    expect(coerce.Get().unwrap_err().message).toContain("Attribute[base.x] is required");
    expect(coerce.Coerce({} as string).unwrap_err().message).toContain("Attribute[base.x] is not a Date");
    expect(coerce.Coerce(6.4 as unknown as string).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(new Date("1970-01-01T00:00:00.006Z"));
    expect(coerce.Coerce("2023-01-01").is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(new Date("2023-01-01"));
    expect(coerce.Coerce(new Date("2023-01-01")).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(new Date("2023-01-01"));
  });

  it("datetime default", () => {
    const coerce = wuesten.AttributeDateTime(attrParam<Formatter.Date.CoerceType>({ default: "2023-01-01" }));
    expect(coerce.Get().unwrap()).toEqual(new Date("2023-01-01"));
  });

  it("datetime default", () => {
    const coerce = wuesten.AttributeDateTime(attrParam<Formatter.Date.CoerceType>({ default: 6 as unknown as string }));
    expect(coerce.Get().unwrap()).toEqual(new Date(6));
  });

  it("datetimeoptional no default", () => {
    const coerce = wuesten.AttributeDateTimeOptional(attrParam());
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.Coerce(null as unknown as Date).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.Coerce(6.4 as unknown as string).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(new Date("1970-01-01T00:00:00.006Z"));
  });

  it("datetimeoptional default", () => {
    const coerce = wuesten.AttributeDateTimeOptional(attrParam<Formatter.Date.CoerceType>({ default: "2023-01-01" }));
    expect(coerce.Get().unwrap()).toEqual(new Date("2023-01-01"));
    coerce.Coerce(undefined);
    expect(coerce.Get().unwrap()).toBeFalsy();
  });
});

describe("integer coerce", () => {
  it("integer no default", () => {
    const coerce = wuesten.AttributeInteger(attrParam());
    expect(coerce.Get().unwrap_err().message).toContain("Attribute[base.x] is required");
    expect(coerce.Coerce({} as number).unwrap_err().message).toContain("Attribute[base.x] is not a number");
    expect(coerce.Coerce(6.4).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(6);
  });

  it("integer default", () => {
    const coerce = wuesten.AttributeInteger(attrParam({ default: 7.2 }));
    expect(coerce.Get().unwrap()).toEqual(7);
  });

  it("integer no default", () => {
    const coerce = wuesten.AttributeIntegerOptional(attrParam());
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.Coerce(null as unknown as number).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.Coerce(6.4).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(6);
  });

  it("integer default", () => {
    const coerce = wuesten.AttributeIntegerOptional(attrParam({ default: 7.2 }));
    expect(coerce.Get().unwrap()).toEqual(7);
    expect(coerce.Coerce(undefined).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
  });
});

describe("float coerce", () => {
  it("float no default", () => {
    const coerce = wuesten.AttributeNumber(attrParam());
    expect(coerce.Get().unwrap_err().message).toContain("Attribute[base.x] is required");
    expect(coerce.Coerce({} as string).unwrap_err().message).toContain("Attribute[base.x] is not a number");
    expect(coerce.Coerce(6.4).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(6.4);
  });

  it("float default", () => {
    const coerce = wuesten.AttributeNumber(attrParam({ default: 7.2 }));
    expect(coerce.Get().unwrap()).toEqual(7.2);
  });

  it("float no default", () => {
    const coerce = wuesten.AttributeNumberOptional(attrParam());
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.Coerce(null as unknown as string).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.Coerce(6.4).is_ok).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(6.4);
  });

  it("float default", () => {
    const coerce = wuesten.AttributeNumberOptional(attrParam({ default: 7.2 }));
    expect(coerce.Get().unwrap()).toEqual(7.2);
    coerce.Coerce(undefined);
    expect(coerce.Get().unwrap()).toBeUndefined();
  });
});

describe("bool coerce", () => {
  it("bool no default", () => {
    const coerce = wuesten.AttributeBoolean(attrParam());
    expect(coerce.Get().unwrap_err().message).toContain("Attribute[base.x] is required");
    expect(coerce.Coerce({} as number).unwrap_err().message).toContain("Attribute[base.x] is not a boolean");
    expect(coerce.Coerce(true).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(true);
    expect(coerce.Coerce("false").is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(false);

    expect(coerce.Coerce("bug").is_err()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(false);

    expect(coerce.Coerce(47).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(true);
  });

  it("bool default", () => {
    const coerce = wuesten.AttributeBoolean(attrParam<BooleanFormatter.CoerceType>({ default: "true" }));
    expect(coerce.Get().unwrap()).toEqual(true);
  });

  it("booloptional no default", () => {
    const coerce = wuesten.AttributeBooleanOptional(attrParam());
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.Coerce(null as unknown as string).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.Coerce("on").is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(true);
    expect(coerce.Coerce("bug").is_err()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(true);
    expect(coerce.Coerce(0).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual(false);
  });

  it("booloptional default", () => {
    const coerce = wuesten.AttributeBooleanOptional(attrParam<BooleanFormatter.CoerceType>({ default: 1 }));
    expect(coerce.Get().unwrap()).toEqual(true);
    expect(coerce.Coerce(undefined).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
  });
});

interface Entity {
  id: string;
  test: number;
}

class Builder implements WuestenBuilder<Entity, Entity> {
  Reflection(): WuestenReflection {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  CoerceAttribute(val: unknown): Result<Entity, Error> {
    throw new Error("Method not implemented.");
  }

  readonly _id = wuesten.AttributeString(attrParam({ jsonname: "id", varname: "Id", base: "base" }));
  readonly _test = wuesten.AttributeInteger(attrParam({ jsonname: "test", varname: "Test", base: "base" }));

  readonly param: WuestenAttributeBase<Entity>;

  constructor(param?: WuestenAttributeParameter<Entity>) {
    this.param = WuestenFactoryAttributeMerge(
      attrParam({
        jsonname: "builder",
        varname: "Builder",
        base: "base",
      }),
      param,
    );
    const base = [this.param.base, this.param.jsonname].join(".");
    this._id = wuesten.AttributeString(attrParam({ jsonname: "id", varname: "Id", base }));
    this._id.CoerceAttribute(param);
    this._test = wuesten.AttributeInteger(attrParam({ jsonname: "test", varname: "Test", base }));
    this._test.CoerceAttribute(param);
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  SetNameSuffix(...idxs: number[]): void {
    throw new Error("Method not implemented.");
  }

  Get(): Result<Entity, Error> {
    throw new Error("Method not implemented.");
  }
  Coerce(value: unknown): Result<Entity, Error> {
    if (typeof value !== "object" || value === null) {
      return Result.Err("not an object");
    }
    const results = {
      id: this._id.CoerceAttribute(value),
      test: this._test.CoerceAttribute(value),
    };
    const errors = Object.values(results)
      .filter((r) => r.is_err())
      .map((r) => r.unwrap_err().message);
    if (errors.length > 0) {
      return Result.Err(errors.join(", "));
    }
    return Result.Ok({
      id: results.id.unwrap(),
      test: results.test.unwrap(),
    });
  }

  id(val: string): Builder {
    this._id.Coerce(val);
    return this;
  }
  test(val: number): Builder {
    this._id.Coerce(val);
    return this;
  }
}

class TestFactory extends WuestenFactory<Entity, Entity, Entity> {
  constructor(param: WuestenAttributeParameter<Entity>) {
    super(param);
  }

  Names(): WuestenNames {
    throw new Error("Method not implemented.");
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Builder(param?: WuestenAttributeParameter<Entity>): WuestenBuilder<Entity, Entity> {
    return new Builder(param);
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  AddFormat(name: string, fn: (recv: unknown) => unknown): TestFactory {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ToObject(typ: Entity): Entity {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  FromPayload(val: WuestePayload, decoder: WuestenDecoder): Result<Entity> {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ToPayload(typ: Entity, encoder?: WuestenEncoder): Result<WuestePayload, Error> {
    throw new Error("Method not implemented.");
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Clone(typ: Entity): Result<Entity, Error> {
    throw new Error("Method not implemented.");
  }
  Schema(): WuestenReflection {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Getter(typ: Entity, base: WuestenReflection[]): WuestenGetterBuilder {
    throw new Error("Method not implemented.");
  }
  // Schema(): WuestenSchema {
  //   throw new Error("Method not implemented.");
  // }
}

describe("object coerce", () => {
  it("object no default", () => {
    const coerce = wuesten.AttributeObject(
      attrParam({ jsonname: "x", varname: "X", base: "super" }),
      new TestFactory(attrParam({ jsonname: "x", varname: "X", base: "super" })),
    );
    expect(coerce.Get().is_err()).toBeTruthy();
    expect(
      coerce
        .Coerce({
          id: "test",
        } as Entity)
        .unwrap_err().message,
    ).toContain("Attribute[super.x.test] not found");
    expect(
      coerce
        .Coerce({
          id: "test",
          Test: 6.4,
        } as unknown as Entity)
        .is_ok(),
    ).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual({
      id: "test",
      test: 6,
    });
    expect(coerce.CoerceAttribute({}).is_err()).toBeTruthy();
    expect(coerce.CoerceAttribute({ o: {} }).unwrap_err().message).toContain("Attribute[super.x] not found");
    expect(coerce.CoerceAttribute({ x: {} }).unwrap_err().message).toContain("Attribute[super.x.id] not found");
    expect(coerce.CoerceAttribute({ x: {} }).unwrap_err().message).toContain("Attribute[super.x.test] not found");
    expect(coerce.CoerceAttribute({ x: { id: "bla" } }).unwrap_err().message).toContain("Attribute[super.x.test] not found");
    expect(coerce.CoerceAttribute({ X: { id: { toString: 4 }, Test: "bla" } }).unwrap_err().message).toContain(
      "Attribute[super.x.id] is not a string",
    );
    expect(coerce.CoerceAttribute({ x: { id: { toString: 4 }, Test: "bla" } }).unwrap_err().message).toContain(
      "Attribute[super.x.test] is not a number",
    );
    expect(coerce.CoerceAttribute({ X: { id: "bla", test: "6.7" } }).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual({
      id: "bla",
      test: 6,
    });
  });

  it("WuestenRecordGetter Nothing", () => {
    const fn = jest.fn();
    WuestenRecordGetter(fn, [], undefined);
    expect(fn.mock.calls.length).toBe(0);
  });

  it("WuestenRecordGetter ObjectEmpty", () => {
    const fn = jest.fn();
    WuestenRecordGetter(fn, [], {});
    expect(fn.mock.calls.length).toBe(0);
  });

  it("WuestenRecordGetter ObjectEmpty", () => {
    const fn = jest.fn();
    WuestenRecordGetter(fn, [], {
      a: 1,
      b: {
        c: 2,
        d: [10, 11],
      },
    });
    expect(fn.mock.calls).toEqual([
      [[{ name: "a", property: undefined, type: "objectitem" }], "a"],
      [[{ name: "a", property: undefined, type: "objectitem" }], 1],
      [[{ name: "b", property: undefined, type: "objectitem" }], "b"],
      [
        [
          { name: "b", property: undefined, type: "objectitem" },
          { name: "c", property: undefined, type: "objectitem" },
        ],
        "c",
      ],
      [
        [
          { name: "b", property: undefined, type: "objectitem" },
          { name: "c", property: undefined, type: "objectitem" },
        ],
        2,
      ],
      [
        [
          { name: "b", property: undefined, type: "objectitem" },
          { name: "d", property: undefined, type: "objectitem" },
        ],
        "d",
      ],
      [
        [
          { name: "b", property: undefined, type: "objectitem" },
          { name: "d", property: undefined, type: "objectitem" },
          { id: "[0]", items: undefined, type: "array" },
        ],
        10,
      ],
      [
        [
          { name: "b", property: undefined, type: "objectitem" },
          { name: "d", property: undefined, type: "objectitem" },
          { id: "[1]", items: undefined, type: "array" },
        ],
        11,
      ],
    ]);
  });

  it("WuestenRecordGetter ArrayEmpty", () => {
    const fn = jest.fn();
    WuestenRecordGetter(fn, [], [4, { a: 1, b: { c: 1 } }]);
    expect(fn.mock.calls).toEqual([
      [[{ id: "[0]", items: undefined, type: "array" }], 4],
      [
        [
          { id: "[1]", items: undefined, type: "array" },
          { name: "a", property: undefined, type: "objectitem" },
        ],
        "a",
      ],
      [
        [
          { id: "[1]", items: undefined, type: "array" },
          { name: "a", property: undefined, type: "objectitem" },
        ],
        1,
      ],
      [
        [
          { id: "[1]", items: undefined, type: "array" },
          { name: "b", property: undefined, type: "objectitem" },
        ],
        "b",
      ],
      [
        [
          { id: "[1]", items: undefined, type: "array" },
          { name: "b", property: undefined, type: "objectitem" },
          { name: "c", property: undefined, type: "objectitem" },
        ],
        "c",
      ],
      [
        [
          { id: "[1]", items: undefined, type: "array" },
          { name: "b", property: undefined, type: "objectitem" },
          { name: "c", property: undefined, type: "objectitem" },
        ],
        1,
      ],
    ]);
    /*
    [

    [[{ id: "[0]", items: undefined, type: "array" }], 4],
    [
      [
        { id: "[1]", items: undefined, type: "array" },
        { name: "a", property: undefined, type: "objectitem" },
      ],
      1,
    ],
    [
      [
        { id: "[1]", items: undefined, type: "array" },
        { name: "b", property: undefined, type: "objectitem" },
        { name: "c", property: undefined, type: "objectitem" },
      ],
      1,
    ],
  ]);
  */
  });

  it("objectoptional no default", () => {
    const coerce = wuesten.AttributeObjectOptional(
      attrParam({ jsonname: "x", varname: "X", base: "super" }),
      new TestFactory(attrParam({ jsonname: "x", varname: "X", base: "super" })),
    );
    expect(coerce.Get().is_ok()).toBeTruthy();
    expect(
      coerce
        .Coerce({
          id: "test",
        } as Entity)
        .is_err(),
    ).toBeTruthy();
    expect(
      coerce
        .Coerce({
          id: "test",
          test: 6.4,
        })
        .is_ok(),
    ).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual({
      id: "test",
      test: 6,
    });
    expect(coerce.Coerce(null as unknown as Entity).is_ok()).toBeTruthy();
    expect(coerce.CoerceAttribute({}).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.CoerceAttribute({ x: {} }).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.CoerceAttribute({ x: { id: "bla" } }).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toBeUndefined();
    expect(coerce.CoerceAttribute({ X: { id: "bla", test: "6.7" } }).is_ok()).toBeTruthy();
    expect(coerce.Get().unwrap()).toEqual({
      id: "bla",
      test: 6,
    });
  });

  // it('object default', () => {
  //     const coerce = wuesten.AttributeBoolean({
  //         default: "true" as unknown as boolean
  //     })
  //     expect(coerce.Get()).toEqual(true)
  // });

  // it('object no default', () => {
  //     const coerce = wuesten.AttributeBooleanOptional()
  //     expect(coerce.Get()).toBeUndefined()
  //     expect(() => coerce.Coerce(null)).not.toThrowError()
  //     expect(coerce.Get()).toBeUndefined()
  //     expect(() => coerce.Coerce("on")).not.toThrowError()
  //     expect(coerce.Get()).toEqual(true)
  // });

  // it('object default', () => {
  //     const coerce = wuesten.AttributeBooleanOptional({
  //         default: 1 as unknown as boolean
  //     })
  //     expect(coerce.Get()).toEqual(true)
  //     coerce.Coerce(undefined)
  //     expect(coerce.Get()).toBeUndefined()
  // });
});

// AttributeObject: <B extends WuestenBuilder<T>, T>(wf: WuestenFactory<B, T>): WuestenAttribute<T> => {
//     return new WuestenAttributeObject(wf);
// },
// AttributeObjectOptional: <B extends WuestenBuilder<T>, T>(wf: WuestenFactory<B, T>): WuestenAttribute<T | undefined> => {
//     return new WuestenAttributeObject(wf);
// },
