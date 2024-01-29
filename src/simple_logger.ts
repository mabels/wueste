import { Level, Logger, WithLogger } from "./logger";

export class SimpleLogger implements Logger {
  out = {} as Record<string, unknown>;
  With(): WithLogger {
    throw new Error("Method not implemented.");
  }
  Msg(...args: string[]): void {
    this.out["msg"] = args.join(" ");
    console.log(JSON.stringify(this.out));
    this.out = {};
  }
  Flush(): Promise<void> {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Module(key: string): Logger {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  SetDebug(...modules: (string | string[])[]): Logger {
    throw new Error("Method not implemented.");
  }
  Str(key: string, value: string): Logger {
    this.out[key] = value;
    return this;
  }
  Error(): Logger {
    this.out["level"] = "error";
    return this;
  }
  Warn(): Logger {
    this.out["level"] = "warn";
    return this;
  }
  Debug(): Logger {
    throw new Error("Method not implemented.");
  }
  Log(): Logger {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  WithLevel(level: Level): Logger {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Err(err: unknown): Logger {
    throw new Error("Method not implemented.");
  }
  Info(): Logger {
    throw new Error("Method not implemented.");
  }
  Timestamp(): Logger {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Any(key: string, value: unknown): Logger {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Dur(key: string, nsec: number): Logger {
    throw new Error("Method not implemented.");
  }
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  Uint64(key: string, value: number): Logger {
    throw new Error("Method not implemented.");
  }
}
