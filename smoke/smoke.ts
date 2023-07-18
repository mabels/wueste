import * as Wueste from "wueste";

import { SimpleTypeFactory  } from "../src-generated/go/simple_type";

const builder = SimpleTypeFactory.Builder();
builder.bool(true);

console.log("Ready for production");
