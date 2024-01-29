import { fromSystem } from "./cli_parser";
import { GenerateGroupConfigFactory } from "./generated/generategroupconfig";
import { SimpleLogger } from "./simple_logger";

describe("cli-test", () => {
  it("cli-test full env single array", () => {
    const out = fromSystem(GenerateGroupConfigFactory, {
      log: new SimpleLogger(),
      env: {
        INPUT_FILES_0: "x",
        INPUT_FILES_1: "z",
        DEBUG: "xxx",
        FILTER_X_KEY: "emno",
      },
      args: [],
    });
    expect(out.Ok()).toEqual({
      debug: "xxx",
      input_files: ["x", "z"],
      output_dir: "./",
      include_path: "./",
      filter: {
        x_key: "emno",
        x_value: "primary-key",
      },
    });
  });

  it("cli-test full env string comma array", () => {
    const out = fromSystem(GenerateGroupConfigFactory, {
      log: new SimpleLogger(),
      env: {
        INPUT_FILES: "x,z",
        DEBUG: "xxx",
        FILTER_X_KEY: "emno",
      },
      args: [],
    });
    expect(out.Ok()).toEqual({
      debug: "xxx",
      input_files: ["x", "z"],
      output_dir: "./",
      include_path: "./",
      filter: {
        x_key: "emno",
        x_value: "primary-key",
      },
    });
  });

  it("cli-test full env single array", () => {
    const out = fromSystem(GenerateGroupConfigFactory, {
      log: new SimpleLogger(),
      env: {
        INPUT_FILES_0: "x",
        INPUT_FILES_1: "z",
        DEBUG: "xxx",
        FILTER_X_KEY: "emno",
      },
      args: [],
    });
    expect(out.Ok()).toEqual({
      debug: "xxx",
      input_files: ["x", "z"],
      output_dir: "./",
      include_path: "./",
      filter: {
        x_key: "emno",
        x_value: "primary-key",
      },
    });
  });

  it("cli-test full args no env", () => {
    const out = fromSystem(GenerateGroupConfigFactory, {
      log: new SimpleLogger(),
      env: {},
      args: ["--input-files", "x", "--input-files", "z", "--debug", "xxx", "--filter-x-key", "emno"],
    });
    expect(out.Ok()).toEqual({
      debug: "xxx",
      input_files: ["x", "z"],
      output_dir: "./",
      include_path: "./",
      filter: {
        x_key: "emno",
        x_value: "primary-key",
      },
    });
  });

  it("cli-test full args override env", () => {
    const out = fromSystem(GenerateGroupConfigFactory, {
      log: new SimpleLogger(),
      env: {
        INPUT_FILES: "a,b",
        FILTER_X_VALUE: "the-key",
      },
      args: ["--input-files", "x", "--input-files", "z", "--debug", "xxx", "--filter-x-key", "emno"],
    });
    expect(out.Ok()).toEqual({
      debug: "xxx",
      input_files: ["x", "z"],
      output_dir: "./",
      include_path: "./",
      filter: {
        x_key: "emno",
        x_value: "the-key",
      },
    });
  });
});
