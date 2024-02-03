import { Logger, MockLogger } from "@adviser/cement";
import { FromSystemResultParsed, fromSystem } from "./cli_parser";
import { GenerateGroupConfigFactory } from "./generated/generategroupconfig";

describe("cli-test", () => {
  it("cli-test full env single array", () => {
    const out = fromSystem(GenerateGroupConfigFactory, {
      log: MockLogger().logger,
      env: {
        INPUT_FILES_0: "x",
        INPUT_FILES_1: "z",
        DEBUG: "xxx",
        FILTER_X_KEY: "emno",
      },
      args: [],
    }) as FromSystemResultParsed<typeof GenerateGroupConfigFactory>;
    expect(out.parsed.Ok()).toEqual({
      debug: "xxx",
      input_files: ["x", "z"],
      output_dir: "./",
      include_path: "./",
      not_selected: false,
      output_format: "TS",
      filter: {
        x_key: "emno",
        x_value: "primary-key",
      },
      filters: [],
    });
  });

  it("cli-test full env string comma array", () => {
    const out = fromSystem(GenerateGroupConfigFactory, {
      log: MockLogger().logger,
      env: {
        INPUT_FILES: "x,z",
        DEBUG: "xxx",
        FILTER_X_KEY: "emno",
      },
      args: [],
    }) as FromSystemResultParsed<typeof GenerateGroupConfigFactory>;
    expect(out.parsed.Ok()).toEqual({
      debug: "xxx",
      input_files: ["x", "z"],
      output_dir: "./",
      include_path: "./",
      not_selected: false,
      output_format: "TS",
      filter: {
        x_key: "emno",
        x_value: "primary-key",
      },
      filters: [],
    });
  });

  it("cli-test full args no env", () => {
    const out = fromSystem(GenerateGroupConfigFactory, {
      log: MockLogger().logger,
      env: {},
      args: ["--input-files", "x", "--input-files", "z", "--debug", "xxx", "--filter-x-key", "emno"],
    }) as FromSystemResultParsed<typeof GenerateGroupConfigFactory>;
    expect(out.parsed.Ok()).toEqual({
      debug: "xxx",
      input_files: ["x", "z"],
      output_dir: "./",
      include_path: "./",
      not_selected: false,
      output_format: "TS",
      filter: {
        x_key: "emno",
        x_value: "primary-key",
      },
      filters: [],
    });
  });

  it("cli-test full args override env", () => {
    const out = fromSystem(GenerateGroupConfigFactory, {
      log: MockLogger().logger,
      env: {
        INPUT_FILES: "a,b",
        FILTER_X_VALUE: "the-key",
      },
      args: ["--input-files", "x", "--input-files", "z", "--debug", "xxx", "--filter-x-key", "emno"],
    }) as FromSystemResultParsed<typeof GenerateGroupConfigFactory>;
    expect(out.parsed.Ok()).toEqual({
      debug: "xxx",
      input_files: ["x", "z"],
      output_dir: "./",
      include_path: "./",
      not_selected: false,
      output_format: "TS",
      filter: {
        x_key: "emno",
        x_value: "the-key",
      },
      filters: [],
    });
  });
});

function logger2console(log: Logger): typeof console {
  return {
    log: (...msg: unknown[]) => log.Info().Any("out", msg).Msg("cli"),
    error: (...msg: unknown[]) => log.Error().Any("out", msg).Msg("cli"),
    warn: (...msg: unknown[]) => log.Warn().Any("out", msg).Msg("cli"),
  } as unknown as typeof console;
}

it("help", async () => {
  const log = MockLogger();
  expect(
    fromSystem(GenerateGroupConfigFactory, {
      log: log.logger,
      parseOut: logger2console(log.logger),
      env: {},
      args: ["--help"],
    }).isHelp,
  ).toBe(true);
  await log.logger.Flush();
  // const logs = log.logCollector.Logs()
  // expect(logs.length).toBe(1)
  // expect({ ...logs[0], out: [cleanCode(logs[0].out.join(""))] }).toEqual({
  //   "level": "info",
  //   "module": "MockLogger",
  //   "msg": "cli",
  //   "out": [
  //     [
  //       "GenerateGroupConfig",
  //       "Configuration for GenerateGroupConfig",
  //       "Options",
  //       "--help                    Prints this usage guide",
  //       "--debug string            Optional. this is debug [env: DEBUG]",
  //       "--not-selected            use all which is not filtered [env: NOT_SELECTED]",
  //       "--output-format string    Defaults to \"TS\". format TS for Typescript, JSchema for JSON Schema [env:",
  //       "OUTPUT_FORMAT]",
  //       "--output-dir string       Defaults to \"./\". undefined [env: OUTPUT_DIR]",
  //       "--include-path string     Defaults to \"./\". undefined [env: INCLUDE_PATH]",
  //       "--input-files string[]    Defaults to []. undefined [env: INPUT_FILES]",
  //       "--filters                 Defaults to []. Optional. undefined [env: FILTERS]",
  //       "--filter-x-key string     Defaults to \"x-groups\". undefined [env: FILTER_X_KEY]",
  //       "--filter-x-value string   Defaults to \"primary-key\". undefined [env: FILTER_X_VALUE]",
  //     ]
  //   ]
  // })
});
