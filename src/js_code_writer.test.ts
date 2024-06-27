import { JSCodeWriter } from "./js_code_writer";
import { MockFileService } from "@adviser/cement/node/mock_file_service";

function toCode(s: string): string[] {
  return s
    .split("\n")
    .map((i) => i.trim())
    .filter((i) => i.length);
}

it("wr.block", async () => {
  const sfs = new MockFileService();
  const rs = await sfs.create("test.js");
  const writer = new JSCodeWriter({
    resultStream: rs,
    fileEngine: sfs,
  });
  writer.block("if (true)", (wr) => {
    wr.writeLn("console.log('true')");
  });
  await writer.close();
  expect(toCode(sfs.files["test.js"].content)).toEqual(["if (true) {", "console.log('true')", "}"]);
});

it("async else async wr.block", async () => {
  const sfs = new MockFileService();
  const rs = await sfs.create("test.js");
  const writer = new JSCodeWriter({
    resultStream: rs,
    fileEngine: sfs,
  });
  await writer.block(
    "if (true)",
    async (wr) => {
      await new Promise<void>((resolve) => {
        setTimeout(() => {
          wr.writeLn("console.log('true')");
          resolve();
        }, 100);
      });
      wr.writeLn("console.log('true2')");
    },
    {
      elseFn: async (wr) =>
        new Promise<void>((resolve) => {
          setTimeout(() => {
            wr.writeLn("console.log('false2')");
            resolve();
          }, 100);
          wr.writeLn("console.log('false')");
        }),
    },
  );
  await writer.close();
  expect(toCode(sfs.files["test.js"].content)).toEqual([
    "if (true) {",
    "console.log('true')",
    "console.log('true2')",
    "} else {",
    "console.log('false')",
    "console.log('false2')",
    "}",
  ]);
});

it("async else no wr.block", async () => {
  const sfs = new MockFileService();
  const rs = await sfs.create("test.js");
  const writer = new JSCodeWriter({
    resultStream: rs,
    fileEngine: sfs,
  });
  await writer.block(
    "if (true)",
    async (wr) => {
      await new Promise<void>((resolve) => {
        setTimeout(() => {
          wr.writeLn("console.log('true')");
          resolve();
        }, 100);
      });
      wr.writeLn("console.log('true2')");
    },
    {
      elseFn: (wr) => {
        wr.writeLn("console.log('false')");
        wr.writeLn("console.log('false2')");
      },
    },
  );
  await writer.close();
  expect(toCode(sfs.files["test.js"].content)).toEqual([
    "if (true) {",
    "console.log('true')",
    "console.log('true2')",
    "} else {",
    "console.log('false')",
    "console.log('false2')",
    "}",
  ]);
});

it("no else async wr.block", async () => {
  const sfs = new MockFileService();
  const rs = await sfs.create("test.js");
  const writer = new JSCodeWriter({
    resultStream: rs,
    fileEngine: sfs,
  });
  await writer.block(
    "if (true)",
    async (wr) => {
      wr.writeLn("console.log('true')");
      wr.writeLn("console.log('true2')");
    },
    {
      elseFn: async (wr) => {
        await new Promise<void>((resolve) => {
          setTimeout(() => {
            wr.writeLn("console.log('false')");
            resolve();
          }, 100);
        });
        wr.writeLn("console.log('false2')");
      },
    },
  );
  await writer.close();
  expect(toCode(sfs.files["test.js"].content)).toEqual([
    "if (true) {",
    "console.log('true')",
    "console.log('true2')",
    "} else {",
    "console.log('false')",
    "console.log('false2')",
    "}",
  ]);
});

it("else wr.block", async () => {
  const sfs = new MockFileService();
  const rs = await sfs.create("test.js");
  const writer = new JSCodeWriter({
    resultStream: rs,
    fileEngine: sfs,
  });
  writer.block(
    "if (true)",
    (wr) => {
      wr.writeLn("console.log('true')");
    },
    {
      else: "dann",
      elseFn: (wr) => {
        wr.writeLn("console.log('false')");
      },
    },
  );
  await writer.close();
  expect(toCode(sfs.files["test.js"].content)).toEqual([
    "if (true) {",
    "console.log('true')",
    "}dann{",
    "console.log('false')",
    "}",
  ]);
});
