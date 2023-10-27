import { Result } from "wueste/result";

import { SimpleTypeFactory } from "../src/generated/go/simpletype";
import { WuestenRetVal } from "../src/wueste";

const test = Result.Ok(42);
if (test.is_err()) {
  console.log(test.unwrap_err());
}
const builder = SimpleTypeFactory.Builder();
builder.bool(true);
builder.string("test");
builder.createdAt(new Date());
builder.float64(1.1);
builder.int64(42);
builder.sub((sub) => {
  if (!sub) {
    throw new Error("Not implemented");
  }
  sub.Test("test");
  sub.Open(() => {
    return WuestenRetVal({});
  });
});
const payload = SimpleTypeFactory.ToPayload(builder.Get());

console.log(`Ready for production: ${payload.Ok().Type}`);
