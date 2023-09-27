import { fromEnv } from "./helper";
import { helperTestFactory } from "./generated/wasm/helper_test";

describe("helper", () => {
  it("from Environment", () => {
    const builder = helperTestFactory.Builder();
    const env = {
      HELPERTEST_TEST: "HELLO",
      HELPERTEST_SUB_BOOL: "true",
      HELPERTEST_SUB_NUM: "1.1",
      HELPERTEST_SUB_INT: "42",
      HELPERTEST_SUB_STR: "HI",
      HELPERTEST_OPT_SUB_BOOL: "false",
      HELPERTEST_OPT_SUB_NUM: "3.3",
      HELPERTEST_OPT_SUB_INT: "12",
      HELPERTEST_OPT_SUB_STR: "BYE",
    };
    fromEnv(builder, env);
    const result = builder.Get();
    expect(result.unwrap_err()).toBeUndefined();
    expect(result.is_ok()).toBeTruthy();
    expect(helperTestFactory.ToObject(result.unwrap())).toEqual({
      test: "HELLO",
      sub: {
        bool: true,
        num: 1.1,
        int: 42,
        str: "HI",
        optBool: false,
        optNum: 3.3,
        optInt: 12,
        optStr: "BYE",
      },
    });
  });
});
