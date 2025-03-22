
# 📝 SHON Specification (v0.4)

**SHON (Simple Human-Oriented Notation)** is a human-readable, JSON-inspired data format designed for configuration, data modeling, and tooling — especially Go-friendly.

---

## ✅ New in v0.4

1. **Type annotations via `$type`**
2. **Constants using `@const`**
3. **Tag system via `$tags`**
4. **Namespace aliasing via `@alias`**

---

## 🧩 Type Annotations

Add `$type` to any object to describe its schema or classification.

```shon
@person {
    sean: {
        name: "Sean",
        phone: "123-456-7890",
        $type: "user"
    }
}
```

- Can be used for polymorphic handling or schema validation
- Helpful in UI tooling and generators

---

## 🔗 Constants

Declare shared constant values using `@const`.

```shon
@const {
    US_PHONE_PATTERN: "^\d{3}-\d{3}-\d{4}$",
    THEME_DEFAULT: "dark"
}
```

- Constants can be referenced in schema or tooling
- Cannot be used directly as runtime values (yet — reference support is planned)

---

## 🏷 Tags

Attach a `$tags` array to any object:

```shon
@feature_flags {
    login: {
        enabled: true,
        $tags: ["public", "beta"]
    },
    payments: {
        enabled: false,
        $tags: ["internal"]
    }
}
```

- Tags can be used for filtering, grouping, environment control
- Tooling can expose or toggle by tags

---

## 📛 Namespace Aliases

Use `@alias` to shorten or remap namespace references:

```shon
@alias {
    addr: address,
    cfg: config
}

@person {
    sean: {
        address: &addr.hq
    }
}
```

- Useful in large projects
- Alias applies to all `&` references after declared

---

## 🔁 Includes

```shon
@include "./address.shon"
@include "./people.shon"
```

---

## 📝 Example (All Features Combined)

```shon
@const {
    US_PHONE_PATTERN: "^\d{3}-\d{3}-\d{4}$"
}

@alias {
    addr: address
}

@address {
    hq: {
        street: "1234 S Main St",
        city: "Palm Springs"
    }
}

@person {
    sean: {
        name: "Sean",
        phone: "123-456-7890",
        address: &addr.hq,
        $type: "user",
        $tags: ["admin", "beta"]
    }
}
```
