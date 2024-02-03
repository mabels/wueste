import * as path from "node:path";
import * as fs from "node:fs";
// the wasm env is limited in size
const whiteListEnv: Record<string, boolean> = {
  HOME: true,
  USER: true,
  SHELL: true,
  PATH: true,
};

process.env = Object.keys(process.env).reduce((a, i) => {
  if (!whiteListEnv[i]) {
    delete a[i];
  }
  return a;
}, process.env);

let whereAreMyDeps = __dirname;
if (fs.existsSync(path.join(whereAreMyDeps, "../dist"))) {
  whereAreMyDeps = path.join(whereAreMyDeps, "../dist/");
}

const execFile = path.join(whereAreMyDeps, "./wasm_exec_node.js");
process.argv = ["node", execFile, path.join(whereAreMyDeps, "./generator.wasm"), ...process.argv];
require(execFile);
//.catch((e) => console.error(e));
