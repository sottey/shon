# SHON Schema Specification v0.6

This document defines the SHON schema format for validating SHON 0.6 documents. A schema defines the expected structure, data types, and constraints of a SHON namespace.

---

## 🧱 Basic Structure

A schema file is itself written in SHON:

```shon
$schema: "0.6"

@namespace {
    field1: { type: "string" },
    field2: { type: "decimal" }
}
```

---

## 📦 Supported Types

| Type       | Description                             |
|------------|-----------------------------------------|
| `string`   | UTF-8 encoded text                      |
| `integer`  | Whole numbers                           |
| `number`   | Any numeric value                       |
| `decimal`  | Fixed-point, high-precision number      |
| `boolean`  | `true` or `false`                       |
| `timestamp`| ISO 8601 datetime string                |
| `array`    | List of values of a specific type       |
| `tuple`    | Fixed-length list with typed positions  |
| `struct`   | Named fields with defined types         |
| `map`      | Arbitrary key/value pairs               |
| `ref`      | Reference to another SHON path          |

---

## 🧱 Fields

### Required fields:
- `type`: must be one of the supported types

### Optional fields:
- `items`: for arrays and tuples
- `properties`: for structs
- `format`: for timestamps
- `required`: list of required field names
- `enum`: list of accepted values

---

## 🔹 Struct Example

```shon
@user {
    type: "struct",
    properties: {
        name: { type: "string" },
        age: { type: "integer" },
        isActive: { type: "boolean" }
    },
    required: ["name", "age"]
}
```

---

## 🔹 Tuple Example

```shon
@point {
    type: "tuple",
    items: ["float", "float"]
}
```

---

## 🔹 Named Tuple Example

```shon
@Vec3 {
    type: "tuple",
    items: [
        { name: "x", type: "float" },
        { name: "y", type: "float" },
        { name: "z", type: "float" }
    ]
}
```

---

## 🔹 Array Example

```shon
@tags {
    type: "array",
    items: { type: "string" }
}
```

---

## 🔹 Decimal and Timestamp

```shon
@transaction {
    amount: { type: "decimal" },
    date: { type: "timestamp", format: "iso8601" }
}
```

---

## 🔹 Map Example

```shon
@translations {
    type: "map",
    values: { type: "string" }
}
```

---

## 🔗 Reference Field

```shon
@team {
    lead: { type: "ref" }
}
```

---

## 🔍 Schema Meta Field

A SHON data file can declare its schema like this:

```shon
$schema: "./user.shos"
```

---

## 🧪 Validation Notes

- Fields not defined in the schema are ignored unless `additionalProperties: false` is used.
- Use `required` to enforce presence.
- Struct field order is preserved for readability, but not enforced.

