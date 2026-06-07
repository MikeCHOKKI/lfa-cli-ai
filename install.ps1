#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Installs LFA CLI on Windows.
.DESCRIPTION
    Downloads the latest LFA CLI binary and adds it to PATH.
.PARAMETER Version
    Specific version to install (default: latest).
.EXAMPLE
    .\install.ps1
    .\install.ps1 -Version 0.1.0
#>

param(
    [string]$Version = "latest"
)

$ErrorActionPreference = "Stop"
$Repo = "MikeCHOKKI/lfa-cli-ai"
$Binary = "lfa"

# Detect architecture
$Arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { "amd64" }
    "ARM64" { "arm64" }
    default { "amd64" }
}

# Get latest version
if ($Version -eq "latest") {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $release.tag_name.TrimStart('v')
}

$binaryName = "lfa-windows-$Arch.exe"
$downloadUrl = "https://github.com/$Repo/releases/download/v$Version/$binaryName"
$installDir = Join-Path $env:USERPROFILE ".local\bin"

if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

Write-Host "Downloading LFA CLI v$Version for windows/$Arch..."
Invoke-WebRequest -Uri $downloadUrl -OutFile (Join-Path $installDir "$Binary.exe")

# Add to PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$installDir*") {
    $newPath = "$currentPath;$installDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Added $installDir to PATH (user-level)"
}

Write-Host ""
Write-Host "LFA CLI v$Version installed successfully!"
Write-Host "Run 'lfa setup' to configure your environment."
