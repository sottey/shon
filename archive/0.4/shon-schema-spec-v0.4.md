
# 📐 SHON Schema Specification (v0.4)

This version adds support for type annotations, constants, tags, and namespace aliases.

---

## ✅ What's New in v0.4

| Feature        | Syntax         | Purpose                        |
|----------------|----------------|--------------------------------|
| Type annotation | `$type`        | Identifies object type         |
| Tags            | `$tags`        | Attaches metadata for grouping |
| Constants       | `@const`       | Reusable values                |
| Namespace alias | `@alias`       | Shorten long namespace paths   |

---

## 🧩 Schema Field Reference

| Field          | Type     | Description |
|----------------|----------|-------------|
| `type`         | string   | "string", "number", "boolean", "null", "object", "array", "ref" |
| `required`     | boolean  | Whether the field must be present |
| `default`      | any      | Default value to apply if field is missing |
| `enum`         | array    | List of allowed values |
| `pattern`      | string   | Regex pattern (only for strings) |
| `keys`         | object   | Field definitions (only for `object` type) |
| `items`        | object   | Schema describing array items |
| `namespace`    | string   | Expected namespace for `ref` type |
| `comment`      | string   | Description of the field (for tooling/docs) |
| `namespacedKey`| boolean  | If true, key may use dot notation |
| `$type`        | string   | Optional field to match declared object type |
| `$tags`        | array    | List of string tags for that object |

---

## 📄 Example Schema

```shon
@schema.person {
    type: "object",
    keys: {
        name: { type: "string", required: true },
        phone: {
            type: "string",
            pattern: &const.US_PHONE_PATTERN,
            comment: "US phone number"
        },
        address: {
            type: "ref",
            namespace: "address"
        },
        $type: {
            type: "string",
            enum: ["user", "admin"]
        },
        $tags: {
            type: "array",
            items: { type: "string" }
        }
    }
}
```

---

## 💡 Notes

- `$type` and `$tags` are reserved metadata fields
- Constants can be used in schemas as values (e.g., regex patterns)
- Namespace aliases affect reference resolution only

