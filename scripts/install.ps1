# LegionTret Windows Install Script
# By Death Legion Team
#
# Run in PowerShell:
#   irm https://raw.githubusercontent.com/deathlegionteam/legiontret/main/scripts/install.ps1 | iex

$ErrorActionPreference = "Stop"

$GITHUB_REPO = "deathlegionteam/legiontret"
$INSTALL_DIR = "$env:USERPROFILE\AppData\Local\LegionTret"
$VERSION = if ($args.Count -gt 0) { $args[0] } else { "latest" }

Write-Host ""
Write-Host "  =======================================================" -ForegroundColor Cyan
Write-Host "  |       LegionTret by Death Legion Team               |" -ForegroundColor Cyan
Write-Host "  |       Run LLMs locally. Simple. Fast. Free.         |" -ForegroundColor Cyan
Write-Host "  =======================================================" -ForegroundColor Cyan
Write-Host ""

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "unknown" }

if ($arch -eq "unknown") {
    Write-Host "Error: Unsupported architecture." -ForegroundColor Red
    exit 1
}

Write-Host "  Detected: windows/$arch" -ForegroundColor Blue

# Determine download URL
$downloadUrl = if ($VERSION -eq "latest") {
    "https://github.com/$GITHUB_REPO/releases/latest/download/legiontret-windows-$arch.exe"
} else {
    "https://github.com/$GITHUB_REPO/releases/download/$VERSION/legiontret-windows-$arch.exe"
}

Write-Host "  Downloading LegionTret $VERSION..." -ForegroundColor Blue

# Create install directory
New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null

# Download
$binaryPath = Join-Path $INSTALL_DIR "legiontret.exe"
try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $binaryPath -UseBasicParsing
} catch {
    Write-Host "Error: Download failed. Please check your internet connection." -ForegroundColor Red
    exit 1
}

# Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$INSTALL_DIR*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$INSTALL_DIR", "User")
    $env:Path = "$env:Path;$INSTALL_DIR"
}

# Verify
Write-Host ""
Write-Host "  =======================================================" -ForegroundColor Green
Write-Host "  LegionTret installed successfully!" -ForegroundColor Green
Write-Host "  =======================================================" -ForegroundColor Green
Write-Host ""
Write-Host "  Get started:" -ForegroundColor White
Write-Host "    legiontret run gemma3       # Run Gemma 3" -ForegroundColor Cyan
Write-Host "    legiontret run llama3       # Run Llama 3" -ForegroundColor Cyan
Write-Host "    legiontret pull mistral     # Download Mistral" -ForegroundColor Cyan
Write-Host "    legiontret list             # List models" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Note: You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
