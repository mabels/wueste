import { fromEnv } from "./helper";
import { helperTestFactory } from "./generated/wasm/helpertest";

describe("helper", () => {
  it("from Environment", () => {
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
