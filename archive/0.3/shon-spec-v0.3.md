
# 📝 SHON Specification (v0.3)

**SHON (Simple Human-Oriented Notation)** is a structured, human-friendly data format designed for readable configuration, extensibility, and tooling.

---

## ✅ New in v0.3

1. **Include directive** for modular files:
```shon
@include "./some-file.shon"
```

2. **Inline comments on keys** (e.g. after values), or optional `_comment.<field>` keys:
```shon
name: "Sean", // full name
_comment.phone: "Must be US format"
```

3. **Namespaced keys** (optional support):
```shon
app.settings.theme: "dark"
```

---

## 🔁 Include Directive

```shon
@include "./address.shon"
@include "./user.shon"
```

- Imports and merges SHON data from another file
- Includes must appear at the top level
- Intended for modularity, reuse, and separation of concerns

---

## 💬 Inline and Field-Level Comments

You can use inline `//` comments:

```shon
phone: "123-456-7890", // Must match US format
```

Or explicitly define comments using `_comment.<field>`:

```shon
phone: "123-456-7890",
_comment.phone: "Must be a valid US number"
```

---

## 🗂 Namespaced Keys (Optional)

```shon
@config {
    app.settings.theme: "dark",
    app.settings.version: "1.2.3"
}
```

- Treats nested namespaces as flat keys using dot notation
- Parsers may choose to expand these to nested objects

---

## 💡 Example SHON v0.3

```shon
@include "./address.shon"

@person {
    sean: {
        name: "Sean", // Full name
        phone: "123-456-7890", // Mobile number
        _comment.phone: "US format only",
        tags: ["go", "home-lab", "automation",],
    }
}

@config {
    app.settings.theme: "dark",
    app.settings.mode: "dev",
}
```
