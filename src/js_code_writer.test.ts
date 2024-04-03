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
