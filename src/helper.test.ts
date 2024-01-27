import {
  fromEnv,
  walk,
  toHash,
  toPathValue,
  groups,
  walkSchema,
  walkSchemaFilter,
  xFilter,
  getValueByAttrName,
  WalkSchemaObjectCollector,
} from "./helper";
import { helperTest, helperTestFactory, helperTestGetter } from "./generated/wasm/helpertest";
import { WuestenFactory, WuestenReflection, WuestenReflectionObject, WuestenReflectionObjectItem, WuestenRetVal } from "./wueste";
import { helperTest$helperTestSubBuilder, helperTest$helperTestSub$arrayBuilder } from "./generated/wasm/helpertest$helpertestsub";

const ref: helperTest = {
  test: "test",
  sub: {
    array: [
      {
        test: "test1",
      },
      {
        test: "test2",
        open: {
          k1: new Date("2023-03-30"),
          a1: 42,
        },
      },
    ],
    bool: true,
    num: 1.1,
    int: 42,
    str: "hi",
    opt_str: "murks",
  },
};

describe("helper", () => {
  it("sanitize json dict", () => {
    const strategy = (a: unknown) => {
      if (typeof a === "number") {
        // floats
        if (a % 1) {
          return "" + a;
        }
        if (a > 0x7fffffff) {
          return "" + a;
        }
      } else if (typeof a === "bigint") {
        return a.toString();
      }
      return a;
    };
    const param = [
      1,
      0x7fffffff,
      0x80000000,
      BigInt("11111111111111111111111111111111111111111"),
      0.78,
      {
        a: 1,
        b: BigInt("11111111111111111111111111111111111111111"),
        c: 0.78,
        d: { a: 1, b: BigInt("11111111111111111111111111111111111111111"), c: 0.78 },
        e: [
          1,
          BigInt("11111111111111111111111111111111111111111"),
          0.78,
          {
            a: 1,
            b: BigInt("11111111111111111111111111111111111111111"),
            c: 0.78,
            d: { a: 1, b: BigInt("11111111111111111111111111111111111111111"), c: 0.78 },
          },
        ],
      },
    ];
    const result = walk(param, strategy) as ArrayLike<unknown>;
    // expect(param[4]).toBe(BigInt("11111111111111111111111111111111111111111"));
    param.push(12);
    expect(result.length).toEqual(param.length - 1);
    expect(result).toEqual([
      1,
      2147483647,
      "2147483648",
      "11111111111111111111111111111111111111111",
      "0.78",
      {
        a: 1,
        b: "11111111111111111111111111111111111111111",
        c: "0.78",
        d: { a: 1, b: "11111111111111111111111111111111111111111", c: "0.78" },
        e: [
          1,
          "11111111111111111111111111111111111111111",
          "0.78",
          {
            a: 1,
            b: "11111111111111111111111111111111111111111",
            c: "0.78",
            d: { a: 1, b: "11111111111111111111111111111111111111111", c: "0.78" },
          },
        ],
      },
    ]);
    // expect(result).toEqual(param);
  });
  it("toHash Exclude String", () => {
    const hash = toHash(helperTestGetter(ref), ["helperTest.sub.helperTestSub.bool"]);
    // echo -n 'testtest1test2a14.200000000000000e+1k12023-03-30T00:00:00.000Z1.100000000000000e+04.200000000000000e+1himurks' | openssl sha1 -hmac ""
    //
    expect(Buffer.from(hash).toString("hex")).toEqual("c9bcb79097342ddec7af9cba01e55a545c6da696");
  });

  it("toHash Exclude Regex", async () => {
    const hash = toHash(helperTestGetter(ref), [/.*\.bool$/]);
    // echo -n 'testtest1test2a14.200000000000000e+1k12023-03-30T00:00:00.000Z1.100000000000000e+04.200000000000000e+1himurks' | openssl sha1 -hmac ""
    expect(Buffer.from(hash).toString("hex")).toEqual("c9bcb79097342ddec7af9cba01e55a545c6da696");
  });

  it("hashit", () => {
    const fn = jest.fn();
    helperTestGetter(ref).Apply(fn);
    expect(fn.mock.calls.map((i) => toPathValue(i[0]))).toEqual([
      "test",
      "test1",
      "test2",
      { a1: 42, k1: new Date("2023-03-30") },
      42,
      new Date("2023-03-30"),
      true,
      1.1,
      42,
      "hi",
      "murks",
    ]);
  });

  it("grouping helper", () => {
    expect(groups(helperTestGetter(ref))).toEqual({
      bool: [
        {
          path: "helperTest.sub.helperTestSub.bool",
          schema: {
            type: "boolean",
            "x-groups": ["group1", "bool"],
          },
          ref: true,
        },
      ],
      group1: [
        {
          path: "helperTest.test",
          ref: "test",
          schema: {
            type: "string",
            "x-groups": ["group1", "test"],
          },
        },
        {
          path: "helperTest.sub.helperTestSub.bool",
          schema: {
            type: "boolean",
            "x-groups": ["group1", "bool"],
          },
          ref: true,
        },
        {
          path: "helperTest.sub.helperTestSub.num",
          schema: {
            type: "number",
            "x-groups": ["group1", "num"],
          },
          ref: 1.1,
        },
        {
          path: "helperTest.sub.helperTestSub.int",
          schema: {
            type: "integer",
            "x-groups": ["group1", "int"],
          },
          ref: 42,
        },
        {
          path: "helperTest.sub.helperTestSub.str",
          schema: {
            type: "string",
            "x-groups": ["group1", "str"],
          },
          ref: "hi",
        },
      ],
      int: [
        {
          path: "helperTest.sub.helperTestSub.int",
          schema: {
            type: "integer",
            "x-groups": ["group1", "int"],
          },
          ref: 42,
        },
      ],
      num: [
        {
          path: "helperTest.sub.helperTestSub.num",
          schema: {
            type: "number",
            "x-groups": ["group1", "num"],
          },
          ref: 1.1,
        },
      ],
      str: [
        {
          path: "helperTest.sub.helperTestSub.str",
          schema: {
            type: "string",
            "x-groups": ["group1", "str"],
          },
          ref: "hi",
        },
      ],
      test: [
        {
          path: "helperTest.test",
          ref: "test",
          schema: {
            type: "string",
            "x-groups": ["group1", "test"],
          },
        },
      ],
    });
  });

  it("coerce is function", () => {
    const builder = helperTestFactory.Builder();
    builder.Coerce(ref);
    // string|number|boolean|Date
    builder.test((test: string) => {
      return WuestenRetVal(test.toUpperCase());
    });
    builder.sub((sub?: helperTest$helperTestSubBuilder) => {
      sub!.array((array?: helperTest$helperTestSub$arrayBuilder) => {
        return WuestenRetVal(array!.Get().unwrap().concat({ test: "test3" }));
      });
    });
    const obj = builder.Get().unwrap();
    expect(obj.test).toEqual("TEST");
    expect(obj.sub.array).toEqual([
      {
        open: undefined,
        test: "test1",
      },
      {
        open: {
          a1: 42,
          k1: new Date("2023-03-30T00:00:00.000Z"),
        },
        test: "test2",
      },
      { test: "test3" },
    ]);
  });

  it.skip("from Environment", () => {
    const builder = helperTestFactory.Builder();
    const env = {
      HELPERTEST_TEST: "HELLO",
      HELPERTEST_SUB_BOOL: "true",
      HELPERTEST_SUB_NUM: "1.1",
      HELPERTEST_SUB_INT: "42",
      HELPERTEST_SUB_STR: "HI",
      HELPERTEST_SUB_OPT_BOOL: "false",
      HELPERTEST_SUB_OPT_NUM: "3.3",
      HELPERTEST_SUB_OPT_INT: "12",
      HELPERTEST_SUB_OPT_STR: "BYE",
    };
    const result = fromEnv(builder, env).Get();
    expect(result.is_ok()).toBeTruthy();
    expect(helperTestFactory.ToObject(result.unwrap())).toEqual({
      test: "HELLO",
      sub: {
        bool: true,
        num: 1.1,
        int: 42,
        str: "HI",
        "opt-bool": false,
        "opt-num": 3.3,
        "opt-int": 12,
        "opt-str": "BYE",
      },
    });
  });
});

it("walk simple number", () => {
  expect(
    walk(4, (x) => {
      return (x as number) + 1;
    }),
  ).toEqual(5);
});

it("walk simple string", () => {
  expect(
    walk("str", (x) => {
      return "test" + x;
    }),
  ).toEqual("teststr");
});

it("walk simple bool", () => {
  expect(
    walk(false, (x) => {
      return !x;
    }),
  ).toEqual(true);
});

it("walk simple array", () => {
  expect(
    walk([4, "str", false], (x) => {
      if (typeof x === "number") {
        return x + 1;
      }
      if (typeof x === "string") {
        return "test" + x;
      }
      if (typeof x === "boolean") {
        return !x;
      }
      return x;
    }),
  ).toEqual([5, "teststr", true]);
});

it("walk simple object", () => {
  expect(
    walk({ x: 4, y: "str", z: false }, (x) => {
      if (typeof x === "number") {
        return x + 1;
      }
      if (typeof x === "string") {
        return "test" + x;
      }
      if (typeof x === "boolean") {
        return !x;
      }
      return x;
    }),
  ).toEqual({ x: 5, y: "teststr", z: true });
});

it("walk literal -> obj", () => {
  expect(
    walk(4, (x) => {
      if (x === 4) {
        return { x: 5, y: "teststr", z: true };
      }
      if (typeof x === "number") {
        return x + 1;
      }
      if (typeof x === "string") {
        return "test" + x;
      }
      if (typeof x === "boolean") {
        return !x;
      }
      return x;
    }),
  ).toEqual({
    x: 6,
    y: "testteststr",
    z: false,
  });
});

it("walk literal -> array", () => {
  expect(
    walk(4, (x) => {
      if (x === 4) {
        return [5, "teststr", true];
      }
      if (typeof x === "number") {
        return x + 1;
      }
      if (typeof x === "string") {
        return "test" + x;
      }
      if (typeof x === "boolean") {
        return !x;
      }
      return x;
    }),
  ).toEqual([6, "testteststr", false]);
});

it("walk array", () => {
  expect(
    walk([0, 1], (x) => {
      if (typeof x === "number") {
        return { x: "" + x + 4 };
      }
      return x;
    }),
  ).toEqual([{ x: "04" }, { x: "14" }]);
});

it("walk object", () => {
  expect(
    walk({ x: 0, y: 1 }, (x) => {
      if (typeof x === "number") {
        return { x: "" + x + 4 };
      }
      return x;
    }),
  ).toEqual({ x: { x: "04" }, y: { x: "14" } });
});

it("walk object-null", () => {
  expect(
    walk(null, (x) => {
      return x;
    }),
  ).toEqual(null);
});

it("walk object-null", () => {
  expect(
    walk({ x: { y: null } }, (x) => {
      return x;
    }),
  ).toEqual({ x: { y: null } });
});

it("walk replace with null", () => {
  expect(
    walk({ x: { y: null } }, (x) => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      if (typeof x === "object" && x !== null && typeof (x as any).y === "object") {
        return null;
      }
      return x;
    }),
  ).toEqual({ x: null });
});

it("walk object replace array", () => {
  expect(
    walk({ x: { y: 7 } }, (x) => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      if (typeof x === "object" && x !== null && (x as any).y) {
        return [7, 8];
      }
      if (7 === x) {
        return 8;
      }
      return x;
    }),
  ).toEqual({ x: [8, 8] });
});

const filterResultHelperTest = [
  [
    {
      id: "https://github.com/mabels/wueste#readme",
      name: undefined,
      title: "helperTest",
      type: "object",
    },
    {
      id: undefined,
      name: "test",
      title: undefined,
      type: "objectitem",
    },
  ],
];

const filterResultHelperTestSub = [
  [
    {
      id: "https://github.com/mabels/wueste#readme",
      name: undefined,
      title: "helperTest",
      type: "object",
    },
    {
      id: undefined,
      name: "sub",
      title: undefined,
      type: "objectitem",
    },
    {
      id: "https://github.com/mabels/wueste#sub",
      name: undefined,
      title: "helperTestSub",
      type: "object",
    },
    {
      id: undefined,
      name: "bool",
      title: undefined,
      type: "objectitem",
    },
  ],
  [
    {
      id: "https://github.com/mabels/wueste#readme",
      name: undefined,
      title: "helperTest",
      type: "object",
    },
    {
      id: undefined,
      name: "sub",
      title: undefined,
      type: "objectitem",
    },
    {
      id: "https://github.com/mabels/wueste#sub",
      name: undefined,
      title: "helperTestSub",
      type: "object",
    },
    {
      id: undefined,
      name: "num",
      title: undefined,
      type: "objectitem",
    },
  ],
  [
    {
      id: "https://github.com/mabels/wueste#readme",
      name: undefined,
      title: "helperTest",
      type: "object",
    },
    {
      id: undefined,
      name: "sub",
      title: undefined,
      type: "objectitem",
    },
    {
      id: "https://github.com/mabels/wueste#sub",
      name: undefined,
      title: "helperTestSub",
      type: "object",
    },
    {
      id: undefined,
      name: "int",
      title: undefined,
      type: "objectitem",
    },
  ],
  [
    {
      id: "https://github.com/mabels/wueste#readme",
      name: undefined,
      title: "helperTest",
      type: "object",
    },
    {
      id: undefined,
      name: "sub",
      title: undefined,
      type: "objectitem",
    },
    {
      id: "https://github.com/mabels/wueste#sub",
      name: undefined,
      title: "helperTestSub",
      type: "object",
    },
    {
      id: undefined,
      name: "str",
      title: undefined,
      type: "objectitem",
    },
  ],
];

interface smallSchema {
  id?: string;
  type: string;
  name?: string;
  title?: string;
}
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function filterPath(b: any[][]): smallSchema[][] {
  return b.map((i) =>
    i.map((j: smallSchema) => ({
      id: j.id,
      type: j.type,
      name: j.name,
      title: j.title,
    })),
  );
}

it("groupsAttributes", () => {
  const walk = jest.fn();
  walkSchema(helperTestFactory.Schema(), walkSchemaFilter(xFilter("x-groups", "group1"), walk));
  expect(filterPath(walk.mock.calls.map((i) => i[0]))).toEqual([...filterResultHelperTest, ...filterResultHelperTestSub]);
});

function testRender(paths: WuestenReflection[][]): string {
  const obj = paths[0][paths[0].length - 2] as WuestenReflectionObject;
  return `class ${obj.title} {${paths
    .map(
      (i) =>
        `${(i[i.length - 1] as WuestenReflectionObjectItem).name}: ${(i[i.length - 1] as WuestenReflectionObjectItem).property.type}`,
    )
    .join(";\n")}
}`;
}

it("Key-Type-Generator", async () => {
  const helper = await import("./generated/wasm/helpertest");

  const factory = getValueByAttrName(helper, (key, val) => {
    if (key.endsWith("Factory")) {
      return val;
    }
  }) as WuestenFactory<unknown, unknown, unknown>;

  const oc = new WalkSchemaObjectCollector();
  walkSchema(factory.Schema(), walkSchemaFilter(xFilter("x-groups", "group1"), oc.add));
  expect(
    Array.from(oc.objects.entries()).map((a) => {
      return [a[0], filterPath(a[1])];
    }),
  ).toEqual([
    ["helperTest", filterResultHelperTest],
    ["helperTestSub", filterResultHelperTestSub],
  ]);
  expect(testRender(oc.objects.get("helperTestSub")!)).toBe(
    `class helperTestSub {bool: boolean;\nnum: number;\nint: integer;\nstr: string\n}`,
  );
});
