import { NestedTypeFactory } from "../../src-generated/go/nested_type";
import { SimpleTypeFactory } from "../../src-generated/go/simple_type";
it("SimpleType-Error", () => {
  const builder = SimpleTypeFactory.Builder();
  builder.sub({ Test: { toString: 5 } as unknown as string });
  builder.float64("WTF" as unknown as number);
  expect(builder.Get().unwrap_err().message).toEqual(
    [
      "Attribute[SimpleType.bool] is required",
      "Attribute[SimpleType.createdAt] is required",
      "Attribute[SimpleType.float64] is required",
      "Attribute[SimpleType.int64] is required",
      "Attribute[SimpleType.string] is required",
      "Attribute[SimpleType.sub] is required",
    ].join("\n"),
  );
});

it("SimpleType-BuilderSet", () => {
  const builder = SimpleTypeFactory.Builder();
  const now = new Date();
  builder.bool(true);
  builder.createdAt(now);
  builder.float64(42.42);
  builder.int64(42);
  builder.string("String42");
  builder.sub({ Test: "Test42" });
  expect(builder.Get().unwrap()).toEqual({
    bool: true,
    createdAt: now,
    default_bool: true,
    default_createdAt: new Date("2023-12-31T23:59:59.000Z"),
    default_float64: 5000,
    default_int64: 64,
    default_string: "hallo",
    float64: 42.42,
    int64: 42,
    opt_sub: undefined,
    optional_bool: undefined,
    optional_createdAt: undefined,
    optional_default_bool: true,
    optional_default_createdAt: new Date("2023-12-31T23:59:59.000Z"),
    optional_default_float32: 50,
    optional_default_int32: 32,
    optional_default_string: "hallo",
    optional_float32: undefined,
    optional_int32: undefined,
    optional_string: undefined,
    string: "String42",
    sub: {
      Test: "Test42",
    },
  });
});

it(`SimpleType-Builder Object-JSON-Object`, () => {
  const builder = SimpleTypeFactory.Builder();
  const now = new Date();
  expect(
    builder
      .Coerce({
        bool: true,
        createdAt: now,
        float64: "42.42",
        int64: "42",
        string: "String42",
        sub: {
          Test: "Test42",
        },
        opt_sub: {
          Test: "Test32",
        },
        optional_bool: true,
        optional_createdAt: now,
        optional_float32: 32.32,
        optional_int32: 32,
        optional_string: "String32",
      })
      .is_ok(),
  ).toBeTruthy();
  builder.float64("42.43");
  const json = JSON.stringify(SimpleTypeFactory.ToObject(builder.Get().unwrap()));
  const fromJson = SimpleTypeFactory.Builder();
  fromJson.Coerce(JSON.parse(json));
  expect(fromJson.Get().unwrap()).toEqual(builder.Get().unwrap());
  expect(fromJson.Get().unwrap().float64).toEqual(42.43);
});

it(`SimpleType-Builder Object-Clone`, () => {
  const builder = SimpleTypeFactory.Builder();
  const now = new Date();
  builder.Coerce({
    bool: true,
    createdAt: now,
    float64: 42.42,
    int64: 42,
    string: "String42",
    sub: {
      Test: "Test42",
    },
    opt_sub: {
      Test: "Test32",
    },
    optional_bool: true,
    optional_createdAt: now,
    optional_float32: 32.32,
    optional_int32: 32,
    optional_string: "String32",
  });
  expect(SimpleTypeFactory.Clone(builder.Get().unwrap()).unwrap()).toEqual(builder.Get().unwrap());
});

it("SimpleType-BuilderCoerce", () => {
  const builder = SimpleTypeFactory.Builder();
  const now = new Date();
  builder.Coerce({
    bool: true,
    createdAt: now,
    float64: 42.42,
    int64: 42,
    string: "String42",
    sub: {
      Test: "Test42",
    },
    opt_sub: {
      Test: "Test32",
    },
    optional_bool: true,
    optional_createdAt: now,
    optional_float32: 32.32,
    optional_int32: 32,
    optional_string: "String32",
  });
  expect(builder.Get().unwrap()).toEqual({
    bool: true,
    createdAt: now,
    default_bool: true,
    default_createdAt: new Date("2023-12-31T23:59:59.000Z"),
    default_float64: 5000,
    default_int64: 64,
    default_string: "hallo",
    float64: 42.42,
    int64: 42,
    opt_sub: {
      Test: "Test32",
    },
    optional_bool: true,
    optional_createdAt: now,
    optional_default_bool: true,
    optional_default_createdAt: new Date("2023-12-31T23:59:59.000Z"),
    optional_default_float32: 50,
    optional_default_int32: 32,
    optional_default_string: "hallo",
    optional_float32: 32.32,
    optional_int32: 32,
    optional_string: "String32",
    string: "String42",
    sub: {
      Test: "Test42",
    },
  });
});

it("SimpleType-BuilderCoerce-Default", () => {
  const builder = SimpleTypeFactory.Builder();
  const now = new Date();
  builder.Coerce({
    bool: true,
    createdAt: now,
    float64: 42.42,
    int64: 42,
    string: "String42",
    sub: {
      Test: "Test42",
    },
    default_int64: 56,
    opt_sub: {
      Test: "Test32",
    },
    optional_bool: true,
    optional_createdAt: now,
    optional_float32: 32.32,
    optional_int32: 32,
    optional_string: "String32",
  });
  expect(builder.Get().unwrap().default_int64).toEqual(56);
  builder.default_int64("42.9" as unknown as number);
  expect(builder.Get().unwrap().default_int64).toEqual(42);
});

it(`NestedType-Builder Object-JSON-Object`, () => {
  const builder = NestedTypeFactory.Builder();
  const now = new Date();
  expect(
    builder
      .Coerce({
        arrayBool: [true, "false"],
        arrayInteger: [42.44, "43.44"],
        arrayNumber: [42.42, "43.43"],
        arrayString: ["42", 43],
        arraySubType: [{ Test: "Test42" }, { Test: 43 }],
        arrayarrayBool: [[[[true]]], [[["false"]]]],
        arrayarrayFlatSchema: [[[[{ Test: "Test42" }], [{ Test: 43 }]]]],
        bool: true,
        float64: 48.9,
        createdAt: now,
        int64: 49,
        string: "xxx",
        sub: {
          Test: "Test42",
        },
      })
      .unwrap(),
  ).toEqual(builder.Get().unwrap());
  const json = JSON.stringify(NestedTypeFactory.ToObject(builder.Get().unwrap()));
  const fromJson = NestedTypeFactory.Builder();
  fromJson.Coerce(JSON.parse(json));
  expect(fromJson.Get().unwrap()).toEqual(builder.Get().unwrap());
  expect(fromJson.Get().unwrap().arrayBool).toEqual([true, false]);
  expect(fromJson.Get().unwrap().arrayInteger).toEqual([42, 43]);
  expect(fromJson.Get().unwrap().arrayNumber).toEqual([42.42, 43.43]);
  expect(fromJson.Get().unwrap().arrayString).toEqual(["42", "43"]);
  expect(fromJson.Get().unwrap().arraySubType).toEqual([{ Test: "Test42" }, { Test: "43" }]);
  expect(fromJson.Get().unwrap().arrayarrayBool).toEqual([[[[true]]], [[[false]]]]);
  expect(fromJson.Get().unwrap().arrayarrayFlatSchema).toEqual([
    [[[{ Test: "Test42", opt_Test: undefined }], [{ Test: "43", opt_Test: undefined }]]],
  ]);
});
