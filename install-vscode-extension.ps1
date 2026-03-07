# Jave VS Code Extension Installer
# Run this script from the repository root to install the Jave language extension

Write-Host "Installing Jave language extension for VS Code..." -ForegroundColor Cyan

# Get the extensions directory
$extensionsDir = "$env:USERPROFILE\.vscode\extensions"
$extensionTarget = Join-Path $extensionsDir "jave-language-0.1.0"
$extensionSource = Join-Path $PSScriptRoot "vscode-jave"

# Create extensions directory if it doesn't exist
if (-not (Test-Path $extensionsDir)) {
    New-Item -ItemType Directory -Force -Path $extensionsDir | Out-Null
    Write-Host "Created extensions directory: $extensionsDir" -ForegroundColor Green
}

# Remove existing symlink or directory if it exists
if (Test-Path $extensionTarget) {
    Write-Host "Removing existing extension at $extensionTarget" -ForegroundColor Yellow
    Remove-Item -Path $extensionTarget -Force -Recurse
}

# Create symlink
try {
    New-Item -ItemType SymbolicLink -Path $extensionTarget -Target $extensionSource -Force | Out-Null
    Write-Host "✓ Extension installed successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Cyan
    Write-Host "1. Reload VS Code: Ctrl+Shift+P → 'Developer: Reload Window'" -ForegroundColor White
    Write-Host "2. Open a .jave file to see syntax highlighting" -ForegroundColor White
    Write-Host ""
    Write-Host "To compile and run Jave programs, install the toolchain:" -ForegroundColor Yellow
    Write-Host "  go install github.com/asciifaceman/jave/cmd/javec@latest" -ForegroundColor Gray
    Write-Host "  go install github.com/asciifaceman/jave/cmd/baggage@latest" -ForegroundColor Gray
    Write-Host "  go install github.com/asciifaceman/jave/cmd/javevm@latest" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Extension location: $extensionTarget" -ForegroundColor Gray
}
catch {
    Write-Host "✗ Failed to create symlink. You may need to run as Administrator." -ForegroundColor Red
    Write-Host ""
    Write-Host "Alternative: Copy the extension manually:" -ForegroundColor Yellow
    Write-Host "  Copy-Item -Recurse -Force ""$extensionSource"" ""$extensionTarget""" -ForegroundColor White
    exit 1
}
