export interface NamedWritableStream {
  readonly name: string;
  readonly stream: WritableStream<Uint8Array>;
}

export interface FileSystem {
  readonly baseDir: string;
  create(fname: string): Promise<NamedWritableStream>;
  readFileString(fname: string): Promise<string>;

  abs(fname: string): string;

  join(...paths: string[]): string;

  relative(from: string, to?: string): string;

  dirname(fname: string): string;

  nodeImport(fname: string): string;
}
