# VS Code Syntax Highlighting for Jave

This document describes the VS Code extension for Jave syntax highlighting.

## Overview

The `vscode-jave` extension provides TextMate-based syntax highlighting for Jave language files (`.jave` and `.jv` extensions).

## Features

### Language Recognition

The extension recognizes:
- **Keywords**: All Jave v0.1 keywords (control flow, declarations, types)
- **Operators**: Word-based operators (`bigly`, `samewise`, `plusalso`, etc.) and arithmetic symbols
- **Types**: Primitive types (`exact`, `vag`, `truther`, `strang`, `nada`, `naw`) and collection types
- **Builtins**: Standard functions (`pront`, `girth`) and modules (`Strangs`, `Pronts`)
- **Literals**: Strings with escape sequences and interpolation directives, numbers, booleans
- **Punctuation**: Statement terminators (`;;`), arrows (`-->`, `->`), assignment (`2b=2`)
- **Comments**: Line (`//`) and block (`/* */`) comments

### Scopes

The grammar assigns semantic scopes following VS Code conventions:

| Scope | Color Theme Mapping |
|-------|---------------------|
| `keyword.control.*` | Control flow keywords (blue) |
| `keyword.operator.*` | Operators (purple/pink) |
| `storage.type.*` | Type names (cyan/green) |
| `entity.name.function.*` | Entry points and function names (yellow) |
| `support.function.*` | Builtin functions (light blue) |
| `support.class.*` | Modules/namespaces (cyan) |
| `constant.language.*` | Boolean/null literals (orange) |
| `constant.numeric.*` | Numbers (green) |
| `string.quoted.*` | String literals (orange/brown) |
| `comment.*` | Comments (gray/green) |

These mappings vary by theme but follow standard conventions.

## Installation

See [vscode-jave/README.md](../vscode-jave/README.md) for installation instructions.

### Prerequisites

The VS Code extension provides syntax highlighting only. To compile and run Jave programs, install the Jave toolchain:

**From source (requires Go 1.26+):**
```bash
go install github.com/asciifaceman/jave/cmd/javec@latest
go install github.com/asciifaceman/jave/cmd/baggage@latest
go install github.com/asciifaceman/jave/cmd/javevm@latest
```

Verify installation:
```bash
javec --version
baggage --version
javevm --version
```

### Extension Installation

#### Quick Install (PowerShell)

From repository root:

```powershell
.\install-vscode-extension.ps1
```

Or run the VS Code task: `Tasks: Run Task` → `Install Jave Extension`

## Grammar Design

### Jave-Specific Challenges

Jave syntax presents unique highlighting challenges:

1. **Word-based operators**: Comparison and logical operators are words (`bigly`, `samewise`) not symbols
2. **Multi-character assignment**: `2b=2` must be recognized as single operator
3. **Two-keyword return**: `give ... up` spans across expressions
4. **Entry point names**: `Foreward`/`Foremost` are special identifiers, not keywords
5. **Module paths**: `highschool/...` in import paths needs special handling
6. **String interpolation**: `%exact`, `%vag`, etc. are not arbitrary placeholders

### Grammar Structure

The TextMate grammar is organized into these repositories:

- `comments` - Line and block comment patterns
- `keywords` - Control flow, declarations, visibility, imports
- `entrypoints` - Special function names (Foreward, Foremost)
- `types` - Primitive and collection type names
- `operators` - Word and symbol operators, assignment
- `literals` - Strings, numbers, booleans, null, collections
- `builtins` - Functions, modules, methods, import paths
- `punctuation` - Statement terminators, arrows, brackets, dots

### Pattern Precedence

Patterns are applied in this order:
1. Comments (highest precedence)
2. Keywords
3. Entry points
4. Types
5. Operators
6. Literals
7. Builtins
8. Punctuation (lowest precedence)

This ensures that keywords are recognized before generic identifiers.

## Testing Syntax Highlighting

### Visual Testing

1. Open any `.jave` file in the `examples/` directory
2. Verify these elements are highlighted:
   - Keywords: `outy`, `seq`, `allow`, `maybe`, `give`, `up`
   - Types: `exact`, `vag`, `truther`, `strang`, `nada`
   - Operators: `2b=2`, `bigly`, `samewise`, `plusalso`
   - Builtins: `pront`, `girth`, `Strangs`, `Combobulate`
   - Punctuation: `;;`, `-->`, `->`
   - Literals: strings, numbers, `yee`/`nee`

### Test Cases

**Conditionals:**
```jave
maybe (<X bigly 5>) -> {
    pront("large");;
}
```
Expected: `maybe`, `bigly`, `pront` all highlighted differently

**Loops:**
```jave
given (<I lessly 10>) again -> {
    I 2b=2 I + 1;;
}
```
Expected: `given`, `again`, `lessly`, `2b=2` highlighted

**Types and Collections:**
```jave
allow table<exact> Numbers 2b=2 [1, 2, 3];;
```
Expected: `allow`, `table`, `exact`, `2b=2` highlighted

**String Interpolation:**
```jave
Strangs.Combobulate<"Value: %exact", X>
```
Expected: `Strangs`, `Combobulate`, `%exact` highlighted

## Known Limitations

### What Works

- All Jave v0.1 keywords and operators
- Standard library modules and functions
- String escape sequences and interpolation directives
- Collection literals (table, enumeration, lexis)
- Comments (line and block)
- Nested structures

### What Doesn't Work Yet

- **Semantic highlighting**: No type checking, scope analysis, or error detection
- **IntelliSense**: No autocomplete or hover information
- **Go to definition**: No symbol navigation
- **Syntax errors**: Invalid syntax is not flagged visually
- **Custom themes**: Some themes may not map scopes optimally

### Future Enhancements (Post v0.1)

- Language server protocol (LSP) implementation
- Real-time diagnostics
- Code completion based on imports and builtins
- Symbol outline and navigation
- Format-on-save support
- Snippet library

## Customizing Colors

To customize Jave syntax colors in your VS Code theme:

1. Open Settings (JSON): `Ctrl+Shift+P` → "Preferences: Open Settings (JSON)"
2. Add token color customizations:

```json
{
  "editor.tokenColorCustomizations": {
    "textMateRules": [
      {
        "scope": "keyword.operator.comparison.word.jave",
        "settings": {
          "foreground": "#FF6B9D"
        }
      },
      {
        "scope": "storage.type.primitive.jave",
        "settings": {
          "foreground": "#4EC9B0"
        }
      }
    ]
  }
}
```

### Useful Scopes for Customization

- `keyword.control.conditional.jave` - maybe, furthermore, otherwise
- `keyword.operator.comparison.word.jave` - bigly, lessly, samewise
- `keyword.operator.assignment.jave` - 2b=2
- `storage.type.primitive.jave` - exact, vag, truther, strang
- `entity.name.function.entrypoint.jave` - Foreward, Foremost
- `support.function.builtin.jave` - pront, girth
- `support.class.builtin.jave` - Strangs, Pronts

## Debugging the Grammar

### Inspecting Tokens

Use VS Code's built-in token inspector:

1. Open Command Palette: `Ctrl+Shift+P`
2. Run: `Developer: Inspect Editor Tokens and Scopes`
3. Click on any token in your Jave file
4. View the assigned scopes and theme color

### Editing the Grammar

The grammar is defined in `vscode-jave/syntaxes/jave.tmLanguage.json`.

After editing:
1. Reload VS Code: `Ctrl+Shift+P` → `Developer: Reload Window`
2. Close and reopen `.jave` files to see changes

### Common Issues

**Keywords not highlighting:**
- Check that pattern uses word boundaries: `\\b(keyword)\\b`
- Verify pattern is in correct repository section
- Ensure keyword spelling matches implementation

**Operators not working:**
- For word operators, use `\\b` boundaries
- For symbol operators, escape special regex characters
- Check pattern precedence order

**Strings breaking:**
- Verify escape sequence patterns
- Check that begin/end patterns are balanced
- Test with nested quotes and escapes

## Contributing

To improve Jave syntax highlighting:

1. Test with real Jave code from `examples/`
2. Identify missing or incorrect highlights
3. Update `jave.tmLanguage.json` patterns
4. Submit changes with test cases

Refer to [Jave v0.1 spec](../specs/jave-v0.1.md) for canonical syntax reference.

## Resources

- [TextMate Language Grammars](https://macromates.com/manual/en/language_grammars)
- [VS Code Syntax Highlighting Guide](https://code.visualstudio.com/api/language-extensions/syntax-highlight-guide)
- [Scope Naming Conventions](https://www.sublimetext.com/docs/scope_naming.html)
- [Jave v0.1 Specification](../specs/jave-v0.1.md)
