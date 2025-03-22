
# 📝 SHON Specification (v0.5)

**SHON (Simple Human-Oriented Notation)** is a human-friendly structured data format designed for readability, validation, and tooling. It is inspired by JSON, but supports comments, namespaces, references, and other extensions suited for configuration and modeling tasks.

---

## 🎯 Goals

- Be readable and writable by humans
- Support structured data like objects and arrays
- Include metadata and schema linkage
- Support comments and modular structure
- Be easily usable with statically typed languages like Go

---

## 🧱 Syntax Overview

### Data Types

| Type     | Example |
|----------|---------|
| String   | `"hello"` or `'''multi-line'''` |
| Number   | `123`, `3.14` |
| Boolean  | `true`, `false` |
| Null     | `null` |
| Object   | `{ key: value }` |
| Array    | `[1, 2, 3,]` (trailing comma allowed) |

---

## 💬 Comments

```shon
// This is a comment
```

- Only single-line comments are supported
- May appear above, inline, or beside fields

---

## 📦 Namespaces

```shon
@person {
    sean: {
        name: "Sean",
        phone: "123-456-7890"
    }
}
```

- Declared using `@namespaceName { ... }`
- The body is a dictionary (object)
- Keys must be unique within each namespace

---

## 🔁 Includes

```shon
@include "./file.shon"
```

- Includes other SHON files into the current context
- Only top-level includes allowed

---

## 🔗 References

```shon
address: &address.main
```

- Use `&namespace.key` to reference another object
- Use `&namespace.arrayField[code=US]` for filtered lookup — always returns an array

---

## 🔧 Metadata Fields

| Field     | Description |
|-----------|-------------|
| `$schema` | Path/URI to schema file for validation (optional) |
| `$type`   | Logical type of the object (e.g., `"user"`, `"admin"`) |
| `$tags`   | Array of descriptive strings used for filtering, environments, etc. |

---

## 🗃 Constants

```shon
@const {
    US_PHONE_PATTERN: "^\d{3}-\d{3}-\d{4}$"
}
```

- Reusable values for schemas, validation, etc.
- Currently not evaluated in-place (used by tools/schemas)

---

## 📛 Namespace Aliases

```shon
@alias {
    addr: address
}
```

- Allows shortened references like `&addr.main`

---

## 📝 Example

```shon
$schema: "./schemas/person.shon"

@const {
    US_PHONE_PATTERN: "^\d{3}-\d{3}-\d{4}$"
}

@alias {
    addr: address
}

@address {
    main: {
        street: "1234 S Main Street",
        city: "Palm Springs",
        zip: "92264"
    }
}

@person {
    sean: {
        name: "Sean",
        phone: "123-456-7890",
        address: &addr.main,
        $type: "user",
        $tags: ["admin", "beta"]
    }
}
```
