import * as Wueste from "wueste";

import { SimpleTypeFactory } from "../src-generated/go/simple_type";

const test = Wueste.Result.Ok("Test");
if (test.is_err()) {
  console.log(test.unwrap_err());
}
const builder = SimpleTypeFactory.Builder();
builder.bool(true);

console.log("Ready for production");
