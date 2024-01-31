import { generateGroupType } from "./generate_group_type";
import { GenerateGroupConfig } from "./generated/generategroupconfig";
import { SimpleLogger } from "./simple_logger";
import { FileSystem, NamedWritableStream } from "./file_system";
import path from "node:path";
import fs from "node:fs/promises";

interface FileCollector {
  readonly name: string;
  content: string;
}

class MockFileSystem implements FileSystem {
  readonly baseDir: string;
  constructor(baseDir: string = process.cwd()) {
    this.baseDir = this.abs(baseDir);
  }

  nodeImport(fname: string): string {
    // console.log('nodeImport:'+ fname);
    if (path.isAbsolute(fname)) {
      return fname;
    } else {
      return "./" + path.normalize(fname);
    }
  }

  readFileString(fname: string): Promise<string> {
    return fs.readFile(fname, { encoding: "utf-8" });
  }

  dirname(fname: string): string {
    return path.dirname(fname);
  }

  basename(fname: string): string {
    return path.basename(fname);
  }

  join(...paths: string[]): string {
    return path.join(...paths);
  }

  relative(from: string, to?: string): string {
    if (to === undefined) {
      to = from;
      from = process.cwd();
    }
    const ret = path.relative(from, to);
    // console.log('relative:'+ from + " -> " + to +   "= " + ret);
    return ret;
  }

  abs(fname: string): string {
    if (this.isAbsolute(fname)) {
      return fname;
    } else {
      const cwd = process.cwd();
      return path.resolve(cwd, fname);
    }
  }

  isAbsolute(fname: string): boolean {
    return path.isAbsolute(fname);
  }

  readonly files = {} as Record<string, FileCollector>;

  async create(fname: string): Promise<NamedWritableStream> {
    let oName = fname;
    if (!path.isAbsolute(fname)) {
      oName = this.abs(fname);
    }

    const fc = {
      name: oName,
      content: "",
    };
    this.files[oName] = fc;
    const decoder = new TextDecoder();

    return {
      name: oName,
      stream: new WritableStream<Uint8Array>({
        write(chunk) {
          fc.content = fc.content + decoder.decode(chunk);
        },
        close() {},
        abort() {
          throw new Error("not implemented");
        },
      }),
    };
  }
}

function cleanCode(code: string): string[] {
  return code
    .split("\n")
    .map((i) => i.trim())
    .filter((i) => i);
}

it("test generated key types", async () => {
  const cfg: GenerateGroupConfig = {
    input_files: ["src/generate_group_type.schema.json"],
    output_dir: "src/generated/keys",
    include_path: "src/generated/wasm",
    filter: {
      x_key: "x-groups",
      x_value: "primary-key",
    },
  };
  const fs = new MockFileSystem();

  await generateGroupType(cfg.input_files[0], {
    log: new SimpleLogger(),
    fs,
    filter: cfg.filter,
    includePath: cfg.include_path,
    outDir: cfg.output_dir,
  });
  expect(Object.keys(fs.files).length).toBe(2);
  expect(Object.values(fs.files)[0].name).toBe(fs.abs("src/generated/keys/generategroupconfigkey.ts"));
  expect(Object.values(fs.files)[1].name).toBe(fs.abs("src/generated/keys/generategroupconfig$filterkey.ts"));
  expect(cleanCode(Object.values(fs.files)[0].content)).toEqual([
    'import { GenerateGroupConfig } from "./../wasm/generategroupconfig";',
    "export interface GenerateGroupConfigKeyType {",
    "readonly debug?: string;",
    "readonly output_dir: string;",
    "}",
    "export class GenerateGroupConfigKey {",
    "static Coerce(val: GenerateGroupConfig): GenerateGroupConfigKeyType {",
    "return {",
    "debug: val.debug,",
    "output_dir: val.output_dir,",
    "};",
    "}",
    "}",
  ]);

  expect(cleanCode(Object.values(fs.files)[1].content)).toEqual([
    'import { GenerateGroupConfig$Filter } from "./../wasm/generategroupconfig$filter";',
    "export interface GenerateGroupConfig$FilterKeyType {",
    "readonly x_key: string;",
    "}",
    "export class GenerateGroupConfig$FilterKey {",
    "static Coerce(val: GenerateGroupConfig$Filter): GenerateGroupConfig$FilterKeyType {",
    "return {",
    "x_key: val.x_key,",
    "};",
    "}",
    "}",
  ]);
});
