import { WuestenAttr, WuestenReflection } from "./wueste";

// type Builder<T, P, O> = WuestenAttr<T, Partial<T> | Partial<P> | Partial<O>>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export function fromEnv<T, P>(builder: WuestenAttr<T, P>, env: Record<string, string>): WuestenAttr<T, P> {
  // const ref = builder.Reflection();
  // if (ref.type !== "object") {
  //   throw new Error("reflection top type must be object");
  // }
  // walk(ref, [ref.title || ref.id || ""], (path, ref) => {
  //   if (ref.type === "object" || ref.type === "array") {
  //     return;
  //   }
  //   // ref.coerceFromString(env[path.map((p) => p.toUpperCase().replace(/[^A-Za-z0-9]/, "_")).join("_")]);
  // });
  return builder;
}

export function fromHash(ref: WuestenReflection, exclude: string[][] = []): string {
  // const ref = builder.Reflection();
  // if (ref.type !== "object") {
  //   throw new Error("reflection top type must be object");
  // }

  // crypto.subtle.digest("SHA-256", new TextEncoder().encode("test")).then((hashBuffer) => {
  // async function digestMessage(message) {
  //   const msgUint8 = new TextEncoder().encode(message); // encode as (utf-8) Uint8Array
  //   const hashBuffer = await crypto.subtle.digest("SHA-256", msgUint8); // hash the message
  //   const hashArray = Array.from(new Uint8Array(hashBuffer)); // convert buffer to byte array
  //   const hashHex = hashArray
  //     .map((b) => b.toString(16).padStart(2, "0"))
  //     .join(""); // convert bytes to hex string
  //   return hashHex;
  // }
  const excl = new Set(exclude.map((p) => p.join("|")));
  walk(ref, [], (path, ref) => {
    const key = path.join("|");
    if (ref.type === "object" || ref.type === "array") {
      return;
    }
    if (excl.has(key)) {
      return;
    }
    // ref.getAsString();
  });
  return "XXXX";
}

function walk(reflection: WuestenReflection, path: string[], actionFN: (path: string[], ref: WuestenReflection) => void) {
  switch (reflection.type) {
    case "object": {
      if (!reflection.properties) {
        throw new Error("no properties");
      }
      actionFN(path, reflection);
      for (const item of reflection.properties) {
        walk(item.property, path.concat(item.name), actionFN);
      }
      break;
    }
    case "array":
      actionFN(path, reflection);
      walk(reflection.items, path.concat(`[${path.length}]`), actionFN);
      break;
    case "string":
    case "number":
    case "boolean":
    case "integer":
      actionFN(path, reflection);
      break;
    default:
      throw new Error("unknown type");
  }
}
