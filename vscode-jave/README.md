# Jave Language Support for VS Code

Syntax highlighting extension for the Jave programming language.

## Features

- Syntax highlighting for `.jave` and `.jv` files
- Keyword recognition for current v0.1 syntax (`outy`, `seq`, `given`, `maybe`, etc.)
- Comment highlighting (line and block)
- Docstring highlighting (`doc< ... >`)
- String literal highlighting with escape sequences
- Builtin and carryon symbol highlighting (`Pront`, `Prontulate`, `Girth`, `Slotify`, `Strangs.Combobulate`)
- Bracket matching and auto-closing pairs

## Installation (Local Development)

### Automated Installation (Cross-Platform)

The easiest way to install is using the provided Go tool, which works on Windows, Linux, and macOS:

**Using Mage (recommended):**
```bash
cd /path/to/jave/repo
mage installExtension
```

**Using Go directly:**
```bash
cd /path/to/jave/repo
go run ./tools/install-extension
```

This tool will:
- Detect your operating system automatically
- Find your VS Code extensions directory
- Create a symlink to `vscode-jave` (or copy if symlink fails)
- Provide helpful next-step instructions

### Manual Installation

If you prefer manual installation:

**Windows (PowerShell):**
```powershell
.\install-vscode-extension.ps1
```

**Linux/macOS (Bash):**
```bash
mkdir -p ~/.vscode/extensions
ln -s "$(pwd)/vscode-jave" ~/.vscode/extensions/jave-language-0.1.0
```

**Copy Method (all platforms):**
Copy `vscode-jave/` to:
- Windows: `%USERPROFILE%\.vscode\extensions\jave-language-0.1.0`
- Linux/macOS: `~/.vscode/extensions/jave-language-0.1.0`

### Verification

1. Open a `.jave` file
2. Check the language mode in the bottom-right status bar - it should show "Jave"
3. Keywords like `outy`, `seq`, `maybe`, `give`, `up` should be highlighted

### Testing the Installer

The installer tool includes tests you can run:
```bash
go test ./tools/install-extension -v
```

## Workspace Settings (Automatic for Jave Developers)

If you're working with Jave code, the extension will be automatically recognized once installed. The Jave repository includes pre-configured VS Code tasks for common workflows:

**For Writing Jave Code:**
- `Compile Current Jave File` - Compile the active `.jave` file to `.jbin`
- `Run Current Jave File` - Compile and execute the active file with `baggage run`
- `Compile and Run Current Jave File` - Full workflow with sponsor messages

**For Jave Toolchain Development:**
- `[Dev] Build Jave Toolchain` - Build javec, baggage, javevm from source
- `[Dev] Install Toolchain to PATH` - Install binaries with `go install`
- `[Dev] Run Jave Tests` - Run Go test suite
- `[Dev] Mage Build/Test` - Use Mage build system

### Prerequisites

To use the Jave code tasks, you need the Jave toolchain installed:

**Option 1: Install from source (if you have Go 1.26+):**
```bash
go install github.com/asciifaceman/jave/cmd/javec@latest
go install github.com/asciifaceman/jave/cmd/baggage@latest
go install github.com/asciifaceman/jave/cmd/javevm@latest
```

**Option 2: Download pre-built binaries (future):**
```bash
# Download from GitHub releases (when available)
# Extract and add to PATH
```

**Verify installation:**
```bash
javec --version
baggage --version
javevm --version
```

You can also add these settings to your VS Code workspace for additional configuration:

```json
{
  "files.associations": {
    "*.jave": "jave",
    "*.jv": "jave"
  },
  "[jave]": {
    "editor.tabSize": 4,
    "editor.insertSpaces": false,
    "editor.detectIndentation": false
  }
}
```

## Syntax Reference

### Keywords

**Visibility:**
- `outy` - Exported symbol
- `inny` - Internal symbol

**Declarations:**
- `allow` - Variable declaration
- `seq`, `sequence` - Function/sequence declaration

**Control Flow:**
- `maybe` - Conditional (if)
- `furthermore` - Else-if
- `otherwise` - Else
- `given` - Loop construct
- `again` - Loop continuation marker (while-style)
- `within` - Iteration marker (for-each style)
- `give ... up` - Return statement

**Entry Points:**
- `Foreward` - Standard entry point sequence
- `Foremost` - Priority entry point sequence

**Imports:**
- `install ... from ...` - Import statement
- `highschool/...` - Standard library path prefix

### Types

**Primitive Types:**
- `exact` - Integer type
- `vag` - Floating-point type
- `truther` - Boolean type
- `strang` - String type
- `nada` - Void/no-value type
- `naw` - Null/nil type

**Collection Types:**
- `table<T>` - Ordered indexed collection
- `enumeration<T>` - Dynamic list
- `lexis<K,V>` - Map/dictionary

### Builtin Functions and Modules

**Functions:**
- `Pront(...)` - Output function
- `Prontulate<...>` - Builtin formatted output
- `Girth(...)` - Size/length function
- `Slotify(...)` - Replace first formatting directive in template text

**Standard Library Modules:**
- `Strangs` - String utilities module (canonical)
- `Srangs` - Legacy alias for Strangs (warns)
- `Embellishments` - Pretty display helpers for structured output

**Module Methods:**
- `Strangs.Combobulate<template, ...>` - String assembly with interpolation

**Comments and Docstrings:**
- `>>|` - Line comment
- `=[ ... ]=` - Block comment
- `doc< ... >` - Structured docstring block (multiline)

### Literals

**Strings:** `"hello world"` with escape sequences `\n`, `\t`, `\"`, `\\`

**String Interpolation Directives:**
- `%exact` - Integer placeholder
- `%vag` - Float placeholder
- `%truther`, `%tru` - Boolean placeholder
- `%strang` - String placeholder
- `%v` - Generic placeholder (legacy, warns)

**Numbers:**
- Integers: `42`, `100`
- Floats: `3.14`, `0.5`

**Booleans:**
- `yee` - True
- `nee` - False

**Null:** `naw`

**Collections:**
- Table: `[1, 2, 3]`
- Multi-dimensional: `[[1, 2], [3, 4]]`
- Enumeration: `<"Ada", "Linus", "Grace">`
- Lexis: `{"Ada": 36, "Linus": 55}`

### Operators

**Comparison (Word Operators):**
- `samewise` - Equal (==)
- `notsamewise` - Not equal (!=)
- `bigly` - Greater than (>)
- `lessly` - Less than (<)
- `biglysame` - Greater or equal (>=)
- `lesslysame` - Less or equal (<=)

**Logical (Word Operators):**
- `plusalso` - Logical AND
- `orelse` - Logical OR
- `notno` - Logical NOT

**Arithmetic (Symbol Operators):**
- `+` - Addition
- `-` - Subtraction
- `*` - Multiplication
- `/` - Division
- `%` - Modulo

**Assignment:**
- `2b=2` - Assignment operator

### Punctuation

- `;;` - Statement terminator
- `-->` - Return type arrow
- `->` - Control flow arrow (for conditionals and loops)
- `<<` `>>` - Type wrapper angles
- `<` `>` - Generic/parameter angles
- `.` - Member access
- `:` - Key-value separator (in lexis literals)

### Control Flow Examples

**Conditionals:**
```jave
maybe (<X bigly 5>) -> {
    Pront("large");;
} furthermore (<X lessly 5>) -> {
    Pront("small");;
} otherwise -> {
    Pront("middle");;
}
```

**While-style Loop:**
```jave
given (<X lesslysame 10>) again -> {
    Pront(X);;
    X 2b=2 X + 1;;
}
```

**For-style Loop:**
```jave
given (<allow exact I 2b=2 0;; I lessly 10;; I 2b=2 I + 1;;>) -> {
    Pront(I);;
}
```

**For-each Iteration:**
```jave
given (<Name within Names>) -> {
    Pront(Name);;
}
```

## Examples

```jave
>>| Hello World
outy seq Foremost<> --> <<nada>> {
    Pront("hello, jave");;
    give up;;
}
```

```jave
>>| With variables and conditionals
outy seq Foremost<> --> <<nada>> {
    allow vag Score 2b=2 0.85;;
    
    maybe (<Score biglysame 0.6>) -> {
        Pront("Pass!");;
    } otherwise -> {
        Pront("Fail!");;
    }
    
    give up;;
}
```

```jave
>>| With imports and string assembly
install Strangs from highschool/English;;

outy seq Foremost<> --> <<nada>> {
    allow exact Count 2b=2 42;;
    allow strang Message 2b=2 Strangs.Combobulate<"Answer: %exact", Count>;;
    Prontulate<"Answer: %exact", Count>;;
    give up;;
}
```

```jave
// With collections and loops
install Strangs from highschool/English;;

outy seq Foremost<> --> <<nada>> {
    allow table<exact> Numbers 2b=2 [1, 2, 3, 4, 5];;
    allow exact Sum 2b=2 0;;
    
    given (<Value within Numbers>) -> {
        Sum 2b=2 Sum + Value;;
    }
    
    Pront(Strangs.Combobulate<"Sum: %exact", Sum>);;
    give up;;
}
```

## Development

To modify the syntax highlighting:

1. Edit `syntaxes/jave.tmLanguage.json` - TextMate grammar definitions
2. Edit `language-configuration.json` - Bracket matching and comment configuration
3. Reload VS Code to test changes

## IntelliSense (LSP Preview)

The extension now includes a lightweight LSP client and can launch `javels` for:

- hover docs from source docstrings/manifests
- signature help for sequence/builtin calls

By default, the extension prefers bundled `javels` binaries included in release extension packages, so end users do not need Go installed.

Default launch strategy:

- `jave.languageServer.useBundled = true`
- bundled binary path: `vscode-jave/bin/javels-<os>-amd64[.exe]`
- fallback command: `javels`
- development fallback: `go run ./cmd/javels` when working in a Jave source workspace
- fallback args: `[]`
- cwd: first workspace folder

This works best when you open the Jave repository (or another workspace that has `cmd/javels`).

You can override launch settings in VS Code:

```json
{
    "jave.languageServer.useBundled": true,
    "jave.languageServer.command": "go",
    "jave.languageServer.args": ["run", "./cmd/javels"],
    "jave.languageServer.cwd": ""
}
```

Developer note: set `jave.languageServer.useBundled` to `false` when testing source changes to `cmd/javels` with `go run`.

Install command for local CLI usage:

```bash
go install github.com/asciifaceman/jave/cmd/javels@latest
```

If installed globally, you can switch to:

```json
{
    "jave.languageServer.command": "javels",
    "jave.languageServer.args": [],
    "jave.languageServer.cwd": ""
}
```

## IntelliSense Roadmap

Current preview covers hover/signature help. For richer completion, diagnostics, and semantic navigation, continue expanding the language server while keeping this syntax extension as the lexical layer.

Recommended next steps:

1. Add `textDocument/completion` from scope-aware identifiers + known stdlib carryon exports.
2. Add `textDocument/publishDiagnostics` by reusing parser/sema diagnostics incrementally.
3. Add `definition/references` for sequences and imports.
4. Add semantic tokens for declaration/reference role coloring beyond TextMate regex highlighting.
5. Add quick docs indexing for generated `site/reference` pages to support "open docs" commands from symbol hover.

Because docstrings are now first-class in the AST, they can directly feed hover/signature help without duplicate annotation formats.

## Troubleshooting

**Extension not loading:**
- Verify the symlink or copy is in the correct location
- Check that the directory name matches `jave-language-0.1.0`
- Reload VS Code window (`Ctrl+Shift+P` → "Developer: Reload Window")

**No syntax highlighting:**
- Check file extension is `.jave` or `.jv`
- Verify language mode in status bar (bottom-right)
- Try manually setting language: `Ctrl+Shift+P` → "Change Language Mode" → "Jave"

**Changes not applying:**
- After editing grammar files, reload VS Code
- You may need to close and reopen `.jave` files
- For persistent issues, restart VS Code completely

**LSP timeout on initialize:**
- This usually means no bundled binary and no `javels` available on PATH.
- In Jave source workspaces, the extension will now automatically attempt `go run ./cmd/javels`.
- Outside source workspaces, either install `javels` or use release extension bundles that include `vscode-jave/bin/javels-*` binaries.

## License

Same as the Jave language project.
