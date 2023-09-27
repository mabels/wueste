import { Result } from "wueste/result";

import { SimpleTypeFactory } from "../src/generated/go/simpletype";

const test = Result.Ok(42);
if (test.is_err()) {
  console.log(test.unwrap_err());
}
const builder = SimpleTypeFactory.Builder();
builder.bool(true);

console.log("Ready for production");
