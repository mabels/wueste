// import { Payload, PayloadFactory } from "../../src/generated/go/payload";

import { NestedTypeFactory, NestedTypeGetter } from "../../src/generated/go/nestedtype";
import { NestedType$Payload, NestedType$PayloadFactory } from "../../src/generated/go/nestedtype$payload";
import { SimpleTypeFactory, SimpleTypeFactoryImpl, SimpleTypeObject, SimpleTypeParam } from "../../src/generated/go/simpletype";
import { WuesteJsonBytesDecoder, WuesteJsonBytesEncoder, WuestenFactoryInferT, WuestenReflectionObject } from "../../src/wueste";

const simpleTypeParam: SimpleTypeParam = {
  bool: true,
  createdAt: new Date(),
  float64: "42.42",
  int64: "42",
  string: "String42",
  sub: {
    Test: "Test42",
    Open: {
      X: {
        Y: {
          Z: 47,
        },
      },
    },
  },
  opt_sub: {
    Test: "Test32",
    Open: {
      X: {
        Y: {
          Z: 49,
        },
      },
    },
  },
  optional_bool: true,
  optional_createdAt: new Date(),
  optional_float32: 32.32,
  optional_int32: 32,
  optional_string: "String32",
};

it("SimpleType-Error", () => {
  const builder = SimpleTypeFactory.Builder();
  builder.sub({
    Test: { toString: 5 } as unknown as string,
    Open: {
      X: {
        Y: {
          Z: 42,
        },
      },
    },
  });
  builder.float64("WTF" as unknown as number);
  expect(builder.Get().unwrap_err().message).toEqual(
    [
      "Attribute[SimpleType.string] is required",
      "Attribute[SimpleType.createdAt] is required",
      "Attribute[SimpleType.float64] is required",
      "Attribute[SimpleType.int64] is required",
      "Attribute[SimpleType.bool] is required",
      "Attribute[SimpleType.sub.Test] is required",
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
  builder.sub({
    Test: "Test42",
    Open: {
      X: {
        Y: {
          Z: 42,
        },
      },
    },
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
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
  });
});

it(`SimpleType-Builder Object-JSON-Object`, () => {
  const builder = SimpleTypeFactory.Builder();
  expect(
    builder
      .Coerce({
        bool: true,
        float64: 42.42,
      })
      .unwrap_err().message,
  ).toEqual(
    [
      "Attribute[SimpleType] is Attribute[SimpleType.string] not found:string",
      "Attribute[SimpleType.createdAt] not found:createdAt",
      "Attribute[SimpleType.int64] not found:int64",
      "Attribute[SimpleType.sub] not found:sub",
    ].join("\n"),
  );
  expect(builder._attr._bool.Get().unwrap()).toEqual(true);
  expect(builder._attr._int64.Get().is_err()).toBeTruthy();
  expect(builder._attr._float64.Get().unwrap()).toEqual(42.42);
});

it(`SimpleType-Builder Object-JSON-Object`, () => {
  const builder = SimpleTypeFactory.Builder();
  expect(builder.Coerce(simpleTypeParam).is_ok()).toBeTruthy();
  builder.float64("42.43");
  const payload = SimpleTypeFactory.ToPayload(builder.Get(), WuesteJsonBytesEncoder).unwrap();
  expect(payload.Type).toBe("https://SimpleType");
  const fromJson = SimpleTypeFactory.Builder();
  fromJson.Coerce(JSON.parse(new TextDecoder().decode(payload.Data as unknown as Uint8Array)));
  expect(fromJson.Get().unwrap()).toEqual(builder.Get().unwrap());
  expect(fromJson.Get().unwrap().float64).toEqual(42.43);
});

it(`SimpleType-Builder Payload-JSON-Payload`, () => {
  const builder = SimpleTypeFactory.Builder();
  expect(builder.Coerce(simpleTypeParam).is_ok()).toBeTruthy();
  const payload = SimpleTypeFactory.ToPayload(builder.Get(), WuesteJsonBytesEncoder).unwrap();
  const fromPayload = SimpleTypeFactory.FromPayload(payload, WuesteJsonBytesDecoder).unwrap();
  expect(fromPayload).toEqual(builder.Get().unwrap());
});

it(`SimpleType-Builder Payload-JSON-Payload`, () => {
  const builder = SimpleTypeFactory.Builder();
  expect(builder.Coerce(simpleTypeParam).is_ok()).toBeTruthy();
  const payload = SimpleTypeFactory.ToPayload(builder.Get().unwrap(), WuesteJsonBytesEncoder).unwrap();
  // const fromPayload = SimpleTypeFactory.Builder();
  (payload as { Type: string }).Type = "Kaput";
  expect(SimpleTypeFactory.FromPayload(payload).unwrap_err().message).toEqual(
    "WuestePayload Type mismatch:[SimpleType,https://SimpleType,SimpleType] != Kaput",
  );
});

it(`Factory Contains Schema`, () => {
  expect((SimpleTypeFactory.Schema() as WuestenReflectionObject).id).toEqual("https://SimpleType");
});

it(`Factory Contains Schema`, () => {
  const builder = SimpleTypeFactory.Builder();
  expect(builder.Coerce(simpleTypeParam).is_ok()).toBeTruthy();
  const fn = jest.fn();
  SimpleTypeFactory.Getter(builder.Get().unwrap()).Apply(fn);
  expect(fn).toBeCalledTimes(30);
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
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
    opt_sub: {
      Test: "Test32",
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
    optional_bool: true,
    optional_createdAt: now,
    optional_float32: 32.32,
    optional_int32: 32,
    optional_string: "String32",
  });
  expect(SimpleTypeFactory.Clone(builder.Get().unwrap()).unwrap()).toEqual(builder.Get().unwrap());
});

it(`SimpleType-Builder ToObject`, () => {
  const builder = SimpleTypeFactory.Builder();
  const now = new Date();
  const dict: SimpleTypeObject = {
    bool: true,
    createdAt: now,
    float64: 42.42,
    int64: 42,
    "default-bool": true,
    "default-createdAt": now,
    "default-float64": 5000,
    "default-int64": 64,
    "default-string": "hallo",
    string: "String42",
    sub: {
      Test: "Test42",
      "opt-Test": "Test32",
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
    "opt-sub": {
      Test: "Test32",
      "opt-Test": "Test32",
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
    "optional-default-bool": true,
    "optional-default-createdAt": now,
    "optional-default-float32": 50,
    "optional-default-int32": 32,
    "optional-default-string": "hallo",
    "optional-bool": true,
    "optional-createdAt": now,
    "optional-float32": 32.32,
    "optional-int32": 32,
    "optional-string": "String32",
  };
  const ref = SimpleTypeFactory.ToObject(builder.Coerce(dict).unwrap());
  expect(ref).toEqual(dict);
});

it("SimpleType-BuilderCoerce", () => {
  const builder = SimpleTypeFactory.Builder();
  const now = new Date("2023-01-27");
  builder.Coerce({
    bool: true,
    createdAt: now,
    float64: 42.42,
    int64: 42,
    string: "String42",
    sub: {
      Test: "Test42",
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
    opt_sub: {
      Test: "Test32",
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
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
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
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
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
  });
});

it("Nested-Getter", () => {
  const nested = createNested();
  const fn = jest.fn();
  NestedTypeGetter(nested).Apply(fn);
  expect(fn.mock.calls.map((i) => i[1])).toEqual([
    true,
    false,
    "42",
    "43",
    42.42,
    43.43,
    42,
    43,
    true,
    false,
    "Test37",
    "X",
    "Y",
    "Z",
    37,
    "39",
    "X",
    "Y",
    "Z",
    39,
    "Test49",
    "X",
    "Y",
    "Z",
    49,
    "47",
    "X",
    "Y",
    "Z",
    47,
    "Test63",
    "X",
    "Y",
    "Z",
    63,

    "xxx",
    "hallo",
    "hallo",
    new Date("2023-01-27T00:00:00.000Z"),
    new Date("2023-12-31T23:59:59.000Z"),
    new Date("2023-12-31T23:59:59.000Z"),
    48.9,
    5000,
    50,
    49,
    64,
    32,
    true,
    true,
    true,
    "Test57",
    "X",
    "Y",
    "Z",
    57,
  ]);
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
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
    default_int64: 56,
    opt_sub: {
      Test: "Test32",
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
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

function createNested() {
  const builder = NestedTypeFactory.Builder();
  const now = new Date("2023-01-27");
  return builder
    .Coerce({
      arrayBool: [true, "false"],
      arrayInteger: [42.44, "43.44"],
      arrayNumber: [42.42, "43.43"],
      arrayString: ["42", 43],
      arraySubType: [
        {
          Test: "Test49",
          Open: {
            X: {
              Y: {
                Z: 49,
              },
            },
          },
        },
        {
          Test: 47,
          Open: {
            X: {
              Y: {
                Z: 47,
              },
            },
          },
        },
      ],
      arrayarrayBool: [[[[true]]], [[["false"]]]],
      arrayarrayFlatSchema: [
        [
          [
            [
              {
                Test: "Test37",
                Open: {
                  X: {
                    Y: {
                      Z: 37,
                    },
                  },
                },
              },
            ],
            [
              {
                Test: 39,
                Open: {
                  X: {
                    Y: {
                      Z: 39,
                    },
                  },
                },
              },
            ],
          ],
        ],
      ],
      bool: true,
      float64: 48.9,
      createdAt: now,
      int64: 49,
      string: "xxx",
      sub: {
        Test: "Test57",
        Open: {
          X: {
            Y: {
              Z: 57,
            },
          },
        },
      },
      sub_flat: {
        Test: "Test63",
        Open: {
          X: {
            Y: {
              Z: 63,
            },
          },
        },
      },
    })
    .unwrap();
}

it(`NestedType-Builder Object-JSON-Object`, () => {
  const builder = NestedTypeFactory.Builder();
  // const now = new Date();
  const nested = createNested();
  builder.Coerce(NestedTypeFactory.ToObject(nested));
  expect(nested).toEqual(builder.Get().unwrap());
  const json = JSON.stringify(NestedTypeFactory.ToObject(builder.Get().unwrap()));
  const fromJson = NestedTypeFactory.Builder();
  fromJson.Coerce(JSON.parse(json));
  const dict = NestedTypeFactory.ToObject(fromJson.Get().unwrap());
  expect(fromJson.Get().unwrap()).toEqual(builder.Get().unwrap());
  expect(dict.arrayBool).toEqual([true, false]);
  expect(dict.arrayInteger).toEqual([42, 43]);
  expect(dict.arrayNumber).toEqual([42.42, 43.43]);
  expect(dict.arrayString).toEqual(["42", "43"]);
  expect(dict.arraySubType).toEqual([
    {
      Test: "Test49",
      Open: {
        X: {
          Y: {
            Z: 49,
          },
        },
      },
    },
    {
      Test: "47",
      Open: {
        X: {
          Y: {
            Z: 47,
          },
        },
      },
    },
  ]);
  expect(dict.arrayarrayBool).toEqual([[[[true]]], [[[false]]]]);
  expect(dict.arrayarrayFlatSchema).toEqual([
    [
      [
        [
          {
            Test: "Test37",
            Open: {
              X: {
                Y: {
                  Z: 37,
                },
              },
            },
          },
        ],
        [
          {
            Test: "39",
            Open: {
              X: {
                Y: {
                  Z: 39,
                },
              },
            },
          },
        ],
      ],
    ],
  ]);
});

it(`Payload OpenObject`, () => {
  const json: NestedType$Payload = {
    Test: "x",
    Open: {
      X: {
        Z: 42,
      },
    },
  };
  const obj = NestedType$PayloadFactory.Builder().Coerce(json).unwrap();
  const ref = NestedType$PayloadFactory.ToObject(obj);

  expect(ref).toEqual(json);
});

it("Test type T I O", () => {
  const t: SimpleTypeFactoryImpl["T"] = {
    bool: true,
    createdAt: new Date(),
    float64: 42.42,
    int64: 42,
    string: "String42",
    default_bool: true,
    default_createdAt: new Date("2023-12-31T23:59:59.000Z"),
    default_float64: 5000,
    default_int64: 64,
    default_string: "hallo",
    sub: {
      Test: "Test42",
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
  };
  expect(t.bool).toBeTruthy();
  expect(t.int64).toEqual(42);
});

it("WuestenFactoryInferT", () => {
  const t: WuestenFactoryInferT<SimpleTypeFactoryImpl> = {
    bool: true,
    createdAt: new Date(),
    float64: 42.42,
    int64: 42,
    string: "String42",
    default_bool: true,
    default_createdAt: new Date("2023-12-31T23:59:59.000Z"),
    default_float64: 5000,
    default_int64: 64,
    default_string: "hallo",
    sub: {
      Test: "Test42",
      Open: {
        X: {
          Y: {
            Z: 42,
          },
        },
      },
    },
  };
  expect(t.bool).toBeTruthy();
  expect(t.int64).toEqual(42);
});
