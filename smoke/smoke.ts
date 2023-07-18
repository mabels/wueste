import { result } from "wueste";

import { SimpleTypeFactory } from "../src-generated/go/simple_type";

const test = result.Result.Ok(42);
if (test.is_err()) {
  console.log(test.unwrap_err());
}
const builder = SimpleTypeFactory.Builder();
builder.bool(true);

console.log("Ready for production");
