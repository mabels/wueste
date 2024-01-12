import { SimpleTypeFactory } from "./generated/simpletype";
import { PayloadFactory } from "wueste/payload";
import { WuestenRetVal } from "wueste/wueste";

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

const payload = PayloadFactory.Builder().Coerce({
  Type: "SimpleType",
  Data: builder.Get().unwrap() as unknown as Record<string, unknown>,
});

const payload2 = SimpleTypeFactory.FromPayload(payload.unwrap());

if (payload2.Ok().float64 != 1.1) {
  throw new Error("float64 mismatch");
}

console.log(`Ready for production: ${payload.Ok().Type}=>${payload2.Ok().float64}`);
