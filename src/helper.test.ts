import { fromEnv, toHash } from "./helper";
import { helperTest, helperTestFactory, helperTestGetter } from "./generated/wasm/helpertest";

const ref: helperTest = {
  test: "test",
  sub: {
    array: [
      {
        test: "test1",
      },
      {
        test: "test2",
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
  it("toHash", async () => {
    const hash = await toHash(helperTestGetter(ref), new Set(["helperTest.sub.bool"]));
    // echo -n 'testtest1test21.100000000000000e+04.200000000000000e+1himurks' | openssl sha1 -hmac ""
    expect(Buffer.from(hash).toString("hex")).toEqual("efb5f75a1899d00208b5a2f0bd7894c532b205f0");
  });

  it("hashit", () => {
    const fn = jest.fn();
    helperTestGetter(ref).Apply(fn);
    expect(fn.mock.calls.map((i) => i[1])).toEqual(["test", "test1", "test2", true, 1.1, 42, "hi", "murks"]);
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
