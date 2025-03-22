
# 📐 SHON Schema Specification (v0.3)

This specification defines the structure and validation rules for SHON documents using SHON-based schemas. Schemas ensure data correctness, documentation, and tooling support.

---

## ✅ What's New in v0.3

1. **Support for field-level comments via `_comment.<field>`**
2. **Support for namespaced keys (e.g., `app.settings.theme`)**
3. **Schemas can validate included files (`@include`) when evaluated together**

---

## 🎯 Schema Field Reference

| Field          | Type     | Description |
|----------------|----------|-------------|
| `type`         | string   | `"string"`, `"number"`, `"boolean"`, `"null"`, `"object"`, `"array"`, `"ref"` |
| `required`     | boolean  | Whether the field must be present |
| `default`      | any      | Default value to apply if field is missing |
| `enum`         | array    | List of allowed values |
| `pattern`      | string   | Regex pattern (only for strings) |
| `keys`         | object   | Field definitions (only for `object` type) |
| `items`        | object   | Schema describing array items |
| `namespace`    | string   | Expected namespace for `ref` type |
| `comment`      | string   | A description of the field (used for tooling/docs) |
| `namespacedKey`| boolean  | If true, key may use dot notation for nested keys |

---

## 🧾 Example Schema

```shon
@schema.person {
    type: "object",
    keys: {
        name: {
            type: "string",
            required: true,
            comment: "Full name of the person"
        },
        phone: {
            type: "string",
            pattern: "^\d{3}-\d{3}-\d{4}$",
            comment: "US phone number"
        },
        address: {
            type: "ref",
            namespace: "address"
        },
        tags: {
            type: "array",
            items: { type: "string" }
        }
    }
}

@schema.config {
    type: "object",
    keys: {
        "app.settings.theme": {
            type: "string",
            enum: ["light", "dark"],
            namespacedKey: true
        },
        "app.settings.mode": {
            type: "string",
            enum: ["dev", "prod"],
            namespacedKey: true
        }
    }
}
```

---

## 🧪 Schema Use Example (matching the schema above)

```shon
@include "./address.shon"

@person {
    sean: {
        name: "Sean",
        phone: "123-456-7890",
        _comment.phone: "Mobile number",
        tags: ["go", "home-lab", "automation"]
    }
}

@config {
    app.settings.theme: "dark",
    app.settings.mode: "dev"
}
```

---

## 🛠 Notes

- Use `comment` for schema tooling; use `_comment.key` in the SHON document itself
- `namespacedKey: true` supports dot notation for shallowly nested data
- A future version may allow defining inheritance or schema composition

