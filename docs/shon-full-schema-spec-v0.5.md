
# 📐 SHON Schema Specification (v0.5)

SHON schemas are SHON documents that describe the expected structure, type, and constraints of other SHON files.

---

## ✅ Supported Schema Fields

| Field          | Type     | Description |
|----------------|----------|-------------|
| `type`         | string   | Type of the value: `"string"`, `"number"`, `"boolean"`, `"null"`, `"object"`, `"array"`, `"ref"` |
| `required`     | boolean  | Whether the field must be present |
| `default`      | any      | A fallback value if not supplied |
| `enum`         | array    | List of valid values |
| `pattern`      | string   | Regex pattern (only for strings) |
| `keys`         | object   | Describes fields of an object (type = `object`) |
| `items`        | schema   | Describes items of an array (type = `array`) |
| `namespace`    | string   | Required namespace (for type = `ref`) |
| `comment`      | string   | Description of field (for docs/tooling) |
| `namespacedKey`| boolean  | Enables support for dot-notated keys |
| `$type`        | string   | Expected logical type (if using type tagging) |
| `$tags`        | array    | List of tags associated with the object |

---

## 🔧 Example Schema

```shon
@schema.person {
    type: "object",
    keys: {
        name: { type: "string", required: true },
        phone: {
            type: "string",
            pattern: &const.US_PHONE_PATTERN,
            comment: "US phone number format"
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

@schema.address {
    type: "object",
    keys: {
        street: { type: "string" },
        city: { type: "string" },
        zip: { type: "string" }
    }
}
```

---

## 🧠 Schema Application

- Can be manually specified using `$schema` in the SHON file
- Tools can load and apply schema rules to validate:
  - Type correctness
  - Presence of required fields
  - Reference integrity
  - Pattern and enum conformity

---

## 🔐 Security Notes

- When using remote schemas, validate access and trust
- If users can supply `$schema`, consider restricting to known paths

