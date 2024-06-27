import { WuestenReflection, WuestenReflectionObject, WuestenReflectionObjectItem } from "./wueste";
import { FileService, NamedWritableStream } from "@adviser/cement";

export function sanitize(str: string) {
  return str.replace(/[^a-zA-Z0-9_]/g, "_");
}
export function quote(str: string) {
  const san = sanitize(str);
  if (san !== str) {
    return JSON.stringify(str);
  }
  return str;
}

export function pascalize(str: string) {
  return str.replace(/(?:^\w|[A-Z]|\b\w|\s+)/g, (match) => {
    if (+match === 0) return ""; // or if (/\s+/.test(match)) for white spaces
    return match.toUpperCase();
  });
}

export function typePath(typ: WuestenReflection[]): WuestenReflectionObjectItem[] {
  const out = [] as WuestenReflectionObjectItem[];
  for (const o of typ) {
    // if (o.type === 'object') {
    //     // out.push(`.${typeName(o as WuestenReflectionObject)}`);
    // }
    if (o.type === "objectitem") {
      out.push(o);
    }
  }
  return out;
}

export function typeName(obj: WuestenReflectionObject, suffix?: string): string {
  const name = obj.title || obj.id;
  if (!name) {
    throw new Error("no name");
  }
  return pascalize(sanitize(name + (suffix ? suffix : "")));
}

export function isOptional(obj: WuestenReflectionObject, oi: WuestenReflectionObjectItem): boolean {
  return obj.required?.indexOf(oi.name) === -1;
}

interface BlockOps {
  readonly open: string;
  readonly close: string;
  readonly indent: string;
  readonly elseFn?: BlockFn;
  readonly else: string;
  readonly indentStep: string;
  readonly noCloseingNewLine: boolean;
  readonly importCollector: ImportCollector;
  readonly resultStream: NamedWritableStream;
  readonly codeWriter: WritableStreamDefaultWriter<Uint8Array>;
  readonly codeBuffer: Uint8Array[];
  readonly fileEngine: FileService;
}

type BlockFn = (writer: JSCodeWriter) => void | Promise<void>;

function toArray<T>(v: T | T[]): T[] {
  if (Array.isArray(v)) {
    return v;
  } else {
    return [v];
  }
}

const encoder = new TextEncoder();

export class JSCodeWriterContext {}

function isPromise<T>(v: unknown): v is Promise<T> {
  return !!v && (v as { then: () => void }).then !== undefined;
}

function resolve(v: void | Promise<void>, next?: () => void | Promise<void>): void | Promise<void> {
  if (isPromise(v)) {
    return v.then(() => {
      if (next) {
        return resolve(next());
      }
    });
  }
  if (next) {
    return resolve(next());
  }
}
export class JSCodeWriter {
  readonly opts: BlockOps;
  constructor(
    opts: Partial<BlockOps> & {
      resultStream: NamedWritableStream;
      fileEngine: FileService;
    },
  ) {
    let codeWriter: WritableStreamDefaultWriter<Uint8Array> | undefined = undefined;
    let codeBuffer: Uint8Array[] | undefined = undefined;
    if (opts && !opts.codeWriter && !opts.codeBuffer) {
      codeBuffer = [];
      codeWriter = new WritableStream({
        write(chunk) {
          codeBuffer!.push(chunk);
        },
      }).getWriter();
    }
    this.opts = {
      codeWriter: opts?.codeWriter ?? codeWriter!,
      codeBuffer: opts?.codeBuffer ?? codeBuffer!,
      indent: "",
      indentStep: "  ",
      noCloseingNewLine: false,
      importCollector: opts?.importCollector ?? new ImportCollector(this),
      ...opts,
      else: opts?.else ?? " else ",
      open: opts?.open ?? "{",
      close: opts?.close ?? "}",
    };
  }
  blockClose(opts: BlockOps) {
    opts.noCloseingNewLine ? this.write(opts.close) : this.writeLn(opts.close);
  }
  blockElse(opts: BlockOps): void | Promise<void> {
    if (opts.elseFn) {
      this.writeLn(opts.close + opts.else + opts.open);
      return resolve(
        opts.elseFn!(
          new JSCodeWriter({
            ...this.opts,
            ...opts,
            indent: this.opts.indent + this.opts.indentStep,
          }),
        ),
        () => {
          this.blockClose(opts);
        },
      );
    } else {
      this.blockClose(opts);
    }
  }

  block(iStmts: string | string[], fn: BlockFn, iOpts?: Partial<BlockOps>): void | Promise<void> {
    const opts = {
      ...this.opts,
      ...iOpts,
    };
    this.writeLn([...toArray(iStmts), opts.open].join(" "));
    return resolve(
      fn(
        new JSCodeWriter({
          ...this.opts,
          ...opts,
          indent: this.opts.indent + this.opts.indentStep,
        }),
      ),
      () => {
        return resolve(this.blockElse(opts));
      },
    );
  }

  write(...iStmts: string[]) {
    this.opts.codeWriter.write(encoder.encode([this.opts.indent, ...toArray(iStmts).join(" ")].join("")));
  }
  writeLn(...iStmts: string[]) {
    this.write(...iStmts);
    this.opts.codeWriter.write(encoder.encode("\n"));
  }

  readonly(...parts: string[]): string {
    return ["readonly", ...parts].join(" ");
  }
  static(...parts: string[]): string {
    return ["static", ...parts].join(" ");
  }
  const(lhs: string, rhs: string): string {
    return ["const", this.assign(lhs, rhs)].join(" ");
  }
  let(lhs: string, rhs: string): string {
    return ["let", this.assign(lhs, rhs)].join(" ");
  }
  generic(type: string, ...args: string[]): string {
    return `${type}<${args.join(", ")}>`;
  }
  declareType(type: string, name: string, optional = false): string {
    return `${name}${optional ? "?" : ""}: ${type}`;
  }
  assign(lhs: string, rhs: string): string {
    return `${lhs} = ${rhs}`;
  }
  tuple(lhs: string, rhs: string, comma = ","): string {
    return `${lhs}: ${rhs}${comma}`;
  }
  call(name: string, ...args: string[]): string {
    return `${name}(${args.join(", ")})`;
  }
  deref(...parts: string[]): string {
    return parts.join(".");
  }
  line(...parts: string[]): string {
    return [...parts].join(" ") + ";";
  }
  import(name: string, file: string, as?: string) {
    this.opts.importCollector.add(name, file, as);
  }

  async close(): Promise<void> {
    await this.opts.codeWriter.close();
    const stream = this.opts.resultStream.stream;
    const writer = stream.getWriter();
    this.opts.importCollector.render(
      new JSCodeWriter({
        resultStream: this.opts.resultStream,
        fileEngine: this.opts.fileEngine,
        codeWriter: writer,
      }),
    );
    for (const chunk of this.opts.codeBuffer!) {
      await writer.write(chunk);
    }
    await writer.close();
  }
}

class ImportItem {
  readonly name: string;
  readonly as?: string;
  constructor(name: string, as?: string) {
    this.name = name;
    this.as = as;
  }
  render(): string {
    return this.as ? `${this.name} as ${this.as}` : `${this.name}`;
  }
}

class ImportCollector {
  readonly imports = new Map<string, ImportItem[]>();
  readonly codeWriter: JSCodeWriter;
  constructor(jcw: JSCodeWriter) {
    this.codeWriter = jcw;
  }
  add(name: string, file: string, as?: string) {
    file = this.codeWriter.opts.fileEngine.abs(file);
    // console.log('add', file);
    let found = this.imports.get(file);
    if (!found) {
      found = [];
      this.imports.set(file, found);
    }
    if (found.find((i) => i.name === name)) {
      return;
    }
    found.push(new ImportItem(name, as));
  }
  render(w: JSCodeWriter): void {
    const imports = Array.from(this.imports.entries());
    imports.sort((a, b) => a[0].localeCompare(b[0]));
    for (const [file, items] of imports) {
      if (items.length === 0) {
        continue;
      }
      //   console.log("result:", this.codeWriter.opts.resultStream.name, "import:", file);
      const ifile = this.codeWriter.opts.fileEngine.nodeImport(
        this.codeWriter.opts.fileEngine.relative(
          this.codeWriter.opts.fileEngine.dirname(this.codeWriter.opts.resultStream.name),
          file,
        ),
      );

      if (items.length > 5) {
        w.block(
          ["import"],
          (w) => {
            for (const item of items) {
              w.writeLn(item.render() + ",");
            }
          },
          {
            close: `} from "${ifile}";`,
          },
        );
      } else {
        w.writeLn(`import { ${items.map((i) => i.render()).join(", ")} } from "${ifile}";`);
      }
    }
    w.writeLn();
  }
}
