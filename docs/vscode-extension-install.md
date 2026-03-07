# VS Code Extension Installation

This document describes how to install the Jave VS Code extension for syntax highlighting.

## Cross-Platform Installation (Recommended)

The repository includes a Go-based installer tool that works on Windows, Linux, and macOS.

### Using Mage

From the repository root:

```bash
mage installExtension
```

### Using Go Directly

```bash
go run ./tools/install-extension
```

### What the Installer Does

The tool:
1. Detects your operating system automatically
2. Locates your VS Code extensions directory (`~/.vscode/extensions`)
3. Creates a symlink from `vscode-jave` to the extensions directory
4. Falls back to copying if symlink creation fails (e.g., on Windows without admin rights)
5. Provides helpful post-install instructions

### Why This Approach?

**Benefits:**
- **Cross-platform**: Same command works on Windows, Linux, macOS
- **Testable**: Go tests verify correct behavior
- **Safe**: Falls back to copying if symlink fails
- **Smart**: Checks for existing installations and toolchain presence

**For Active Development:**
Symlinks are preferred because changes to `vscode-jave/` apply immediately without reinstalling.

## Platform-Specific Installation

If you don't have Go installed or prefer platform-specific methods:

### Windows (PowerShell)

```powershell
.\install-vscode-extension.ps1
```

Requires PowerShell 5.1+ (standard on Windows 10+). May require running as Administrator to create symlinks.

### Linux/macOS (Bash)

```bash
chmod +x install-vscode-extension.sh
./install-vscode-extension.sh
```

Or manually:
```bash
mkdir -p ~/.vscode/extensions
ln -s "$(pwd)/vscode-jave" ~/.vscode/extensions/jave-language-0.1.0
```

### Manual Copy Method (All Platforms)

If symlinks are not an option:

**Windows:**
```powershell
Copy-Item -Recurse -Force vscode-jave "$env:USERPROFILE\.vscode\extensions\jave-language-0.1.0"
```

**Linux/macOS:**
```bash
cp -r vscode-jave ~/.vscode/extensions/jave-language-0.1.0
```

**Note:** Copying means you must reinstall after any changes to the extension.

## VS Code Tasks

The repository includes a VS Code task for installation:

1. Open Command Palette: `Ctrl+Shift+P`
2. Run: `Tasks: Run Task`
3. Select: `Install Jave Extension`

This runs the Go installer tool internally.

## Verification

After installation:

1. **Reload VS Code**: `Ctrl+Shift+P` → `Developer: Reload Window`
2. **Open a `.jave` file** from `examples/`
3. **Check language mode** in bottom-right status bar (should show "Jave")
4. **Verify highlighting**:
   - Keywords: `outy`, `seq`, `maybe`, `give`, `up`
   - Types: `exact`, `vag`, `truther`, `strang`
   - Operators: `2b=2`, `bigly`, `samewise`
   - Builtins: `pront`, `girth`, `Strangs`, `Pronts`, `Combobulate`, `Prontulate`

`Strangs.Combobulate<...>` and `Pronts.Prontulate<...>` are part of Jave v0.1 and should be recognized/highlighted.

### Troubleshooting

**Extension not loading:**
- Ensure the directory name is exactly `jave-language-0.1.0`
- Check the symlink/copy is in `~/.vscode/extensions/`
- Try closing and reopening VS Code completely
- Check VS Code output: `View` → `Output` → `Extensions`

**No syntax highlighting:**
- Verify file extension is `.jave` or `.jv`
- Check language mode in status bar (bottom-right)
- Manually set language: `Ctrl+Shift+P` → `Change Language Mode` → `Jave`

**Symlink creation failed:**
- On Windows, try running PowerShell as Administrator
- Alternatively, use the copy method (no symlink needed)
- The Go installer automatically falls back to copying

**Changes not applying:**
- If installed via copy, you must reinstall after changes
- If installed via symlink, just reload VS Code
- Try: `Ctrl+Shift+P` → `Developer: Reload Window`

## Testing the Installer

The Go installer includes comprehensive tests:

```bash
# Run all tests
go test ./tools/install-extension -v

# Test specific functionality
go test ./tools/install-extension -v -run TestCopyDir
```

**Test Coverage:**
- VS Code extensions directory detection (Windows/Linux/macOS)
- File copying with permission preservation
- Directory recursive copying
- Symlink creation
- Toolchain presence detection

## Uninstalling

To remove the extension:

```bash
# Manual removal
rm -rf ~/.vscode/extensions/jave-language-0.1.0
```

Or on Windows:
```powershell
Remove-Item -Recurse -Force "$env:USERPROFILE\.vscode\extensions\jave-language-0.1.0"
```

Then reload VS Code.

## For Package Maintainers

If distributing Jave as a package, you can:

1. **Bundle the extension** in `/usr/share/jave/vscode-jave/` or equivalent
2. **Provide installation script** that symlinks to user's extensions directory
3. **Use the Go installer** as a reference implementation

Example systemwide installation:
```bash
# Install extension globally
sudo mkdir -p /usr/share/jave
sudo cp -r vscode-jave /usr/share/jave/

# User can then symlink
ln -s /usr/share/jave/vscode-jave ~/.vscode/extensions/jave-language-0.1.0
```

## See Also

- [vscode-jave/README.md](../vscode-jave/README.md) - Extension usage and features
- [docs/syntax-highlighting.md](syntax-highlighting.md) - Grammar design and customization
- [specs/jave-v0.1.md](../specs/jave-v0.1.md) - Language specification
