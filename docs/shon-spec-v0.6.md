# SHON Specification v0.6

SHON (Structured Human-Optimized Notation) is a data serialization format designed for readability, schema support, and practical use in modern systems. Version 0.6 introduces advanced types and syntax improvements.

---

## 🔧 Syntax Features

### ✅ Comments
- **Single-line**: `// comment here`
- **Multi-line**: `/* this is a
   multi-line comment */`

### ✅ Namespaces
```shon
@users {
  ...
}
```

### ✅ Field Assignment
```shon
key: value
```

---

## 📦 Supported Types

| SHON Syntax              | Description                      |
|--------------------------|----------------------------------|
| `"string"`               | String                           |
| `42`, `true`, `false`    | Number and boolean               |
| `$decimal("12.34")`      | Decimal with precision           |
| `$timestamp("2024-01-01T00:00:00Z")` | ISO 8601 timestamp      |
| `$tuple(1, "a", true)`   | Anonymous tuple                  |
| `Vec3(1.0, 2.0, 3.0)`    | Named tuple                      |
| `[1, 2, 3]`              | Array                            |
| `{ key: value }`         | Map or Struct (based on schema)  |
| `&ref.to.path`           | Reference                        |

---

## 🧱 Data Structures

### 🔹 Arrays
```shon
numbers: [1, 2, 3]
```

### 🔹 Maps
```shon
translations: {
    en: "Hello",
    es: "Hola"
}
```

### 🔹 Structs
Structs look like maps but are validated against a schema with fixed fields.
```shon
user: {
    name: "Sean",
    active: true
}
```

### 🔹 Tuples
```shon
$tuple(1, "a", true)
Vec3(1.0, 2.0, 3.0) // Named tuple
```

---

## 🔗 References
```shon
manager: &people.sean
```

---

## 🕓 Timestamps
```shon
created: $timestamp("2025-03-22T14:30:00Z")
```

---

## 💵 Decimals
```shon
price: $decimal("19.95")
```

---

## 🔄 Example
```shon
@invoice {
    id: "INV001",
    total: $decimal("1042.75"),
    created: $timestamp("2025-03-22T10:00:00Z"),
    items: [
        $tuple("Widget", 3, $decimal("9.99"))
    ],
    paid: false
}
```

