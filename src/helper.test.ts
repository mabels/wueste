import { fromEnv } from "./helper";
import { helperTestFactory } from "./generated/wasm/helper_test";
import { WuestenReflection } from "./wueste";
import { helperTestSubParam } from "./generated/wasm/helper_test_sub";

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
    builder.Reflection = (): WuestenReflection => {
      return {
        type: "object",
        coerceFromString: () => {
          throw new Error("not implemented");
        },
        properties: [
          {
            name: "test",
            property: {
              type: "string",
              coerceFromString: builder.test.bind(builder),
            },
          },
          {
            name: "sub",
            property: {
              type: "object",
              coerceFromString: () => {
                throw new Error("not implemented");
              },
              properties: [
                {
                  name: "bool",
                  property: {
                    type: "boolean",
                    coerceFromString: (arg) => {
                      builder.sub({ bool: arg } as unknown as helperTestSubParam);
                    },
                  },
                },
                {
                  name: "num",
                  property: {
                    type: "number",
                    coerceFromString: (arg) => {
                      builder.Coerce({ sub: { num: arg } as unknown as helperTestSubParam });
                    },
                  },
                },
                {
                  name: "int",
                  property: {
                    type: "integer",
                    coerceFromString: (arg) => {
                      builder.Coerce({ sub: { int: arg } as unknown as helperTestSubParam });
                    },
                  },
                },
                {
                  name: "str",
                  property: {
                    type: "string",
                    coerceFromString: (arg) => {
                      builder.Coerce({ sub: { str: arg } as unknown as helperTestSubParam });
                    },
                  },
                },
                {
                  name: "opt_bool",
                  property: {
                    type: "boolean",
                    coerceFromString: (arg) => {
                      builder.Coerce({ sub: { opt_bool: arg } as unknown as helperTestSubParam });
                    },
                  },
                },
                {
                  name: "opt_num",
                  property: {
                    type: "number",
                    coerceFromString: (arg) => {
                      builder.Coerce({ sub: { opt_num: arg } as unknown as helperTestSubParam });
                    },
                  },
                },
                {
                  name: "opt_int",
                  property: {
                    type: "integer",
                    coerceFromString: (arg) => {
                      builder.Coerce({ sub: { opt_int: arg } as unknown as helperTestSubParam });
                    },
                  },
                },
                {
                  name: "opt_str",
                  property: {
                    type: "string",
                    coerceFromString: (arg) => {
                      builder.Coerce({ sub: { opt_str: arg } as unknown as helperTestSubParam });
                    },
                  },
                },
              ],
            },
          },
        ],
      };
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
