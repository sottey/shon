
# 📐 SHON Schema Specification (v0.2)

This defines the SHON schema format, used to validate SHON documents.

## Supported Fields

| Field      | Type    | Description |
|------------|---------|-------------|
| `type`     | string  | "string", "number", "boolean", "null", "object", "array", "ref" |
| `required` | boolean | Field must be present |
| `default`  | any     | Default value |
| `enum`     | array   | Allowed values |
| `pattern`  | string  | Regex (for strings) |
| `keys`     | object  | Field definitions (object only) |
| `items`    | schema  | Array item schema |
| `namespace`| string  | Expected namespace (for `ref`) |

---

## Example: Schema for a "person" object

```shon
@schema.person {
    type: "object",
    keys: {
        name: { type: "string", required: true },
        bio: { type: "string" },
        address: { type: "ref", namespace: "address" },
        tags: { type: "array", items: { type: "string" }},
    }
}

@schema.address {
    type: "object",
    keys: {
        street: { type: "string", required: true },
        city: { type: "string" },
        zip: { type: "string" }
    }
}
```
