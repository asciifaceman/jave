# Bundled Jave LSP Binaries

This directory is populated during release packaging with platform-specific `javels` binaries:

- `javels-linux-amd64`
- `javels-darwin-amd64`
- `javels-windows-amd64.exe`

The VS Code extension launches these binaries by default so end users do not need Go installed.

For local developer workflows, if bundled binaries are absent, configure fallback launch settings:

```json
{
  "jave.languageServer.useBundled": false,
  "jave.languageServer.command": "go",
  "jave.languageServer.args": ["run", "./cmd/javels"]
}
```
