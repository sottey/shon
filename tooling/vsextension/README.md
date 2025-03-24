# SHON Language Support for VSCode

This extension adds syntax highlighting and basic support for [SHON](https://github.com/sottey/shon) (`.shon`) and SHON Schema (`.shos`) files in Visual Studio Code.

SHON (Structured Human-Optimized Notation) is a human-readable data format with built-in support for:
- Namespaces
- Tuples (including named)
- References
- Decimal and timestamp types
- Structs and maps
- Comments (single-line and multi-line)

## Features

✅ Syntax highlighting for:
- `@namespace` declarations
- `$typed()` values like `$decimal()` and `$timestamp()`
- Named tuples (e.g., `Vec3(...)`)
- SHON references: `&namespace.key`
- Comments: `// single-line` and `/* multi-line */`

✅ Bracket/brace matching  
✅ Auto-closing for quotes and braces  
✅ Support for both `.shon` and `.shos` files

## Installation

### Manual

1. Clone or download this repository.
2. Open the folder in VSCode.
3. Press `F5` to launch the extension development host.
4. Create or open `.shon` or `.shos` files.

### VSIX (if built)

```bash
code --install-extension shon-0.6.0.vsix
