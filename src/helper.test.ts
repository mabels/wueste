import { fromEnv, toHash } from "./helper";
import { helperTest, helperTestFactory, helperTestGetter } from "./generated/wasm/helpertest";
import { WuestenRetVal } from "./wueste";
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
  it("toHash Exclude String", () => {
    const hash = toHash(helperTestGetter(ref), ["helperTest.sub.helperTestSub.bool"]);
    // echo -n 'testtest1test2a14.200000000000000e+1k12023-03-30T00:00:00.000Z1.100000000000000e+04.200000000000000e+1himurks' | openssl sha1 -hmac ""
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
    expect(fn.mock.calls.map((i) => i[1])).toEqual([
      "test",
      "test1",
      "test2",
      "a1",
      42,
      "k1",
      new Date("2023-03-30").toISOString(),
      true,
      1.1,
      42,
      "hi",
      "murks",
    ]);
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
