
# 📝 SHON Specification (v0.2)

**SHON (Simple Human-Oriented Notation)** is a human-friendly data format inspired by JSON, designed for readable configuration and structured data modeling.

---

## ✅ New in v0.2

1. **Quoted keys** are allowed (e.g., `"name": "Sean"`).
2. **Trailing commas** are allowed in objects and arrays.
3. **Multi-line strings** are supported using triple single quotes (`'''`).

---

## 🧱 Syntax Basics

### Data Types

- Strings: `"text"` or `'''multiline'''`
- Numbers: `123`, `4.56`
- Booleans: `true`, `false`
- Null: `null`
- Arrays: `[1, 2, 3,]`  ← trailing comma allowed
- Objects:
```shon
{
    name: "Sean",
    "type": "admin", // quoted key
}
```

### Comments

```shon
// This is a comment
```

---

## 📦 Namespaces

Namespaces begin with `@<name>`:

```shon
@person {
    sean: {
        name: "Sean",
        bio: '''
            Developer, comic, tech junkie.
            Loves Go, Palm Springs, and automation.
        ''',
    },
}
```

---

## 🔗 References

- Simple reference: `&namespace.key`
- Filtered reference (always returns array): `&namespace.arrayField[key=value]`

---

## 📝 Schema Support

Schemas are written in SHON and describe types, required fields, patterns, and more.

