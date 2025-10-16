<#
.SYNOPSIS
    Builds BalancaCLP.exe (PyInstaller onedir bundle) for Windows.

.DESCRIPTION
    - Verifies Python is installed (3.9+ recommended; on Windows Server 2012, use 3.9).
    - Creates a virtualenv under windows\.venv
    - Installs dependencies from windows\requirements-win.txt
    - Checks for snap7.dll at windows\snap7\snap7.dll (download manually if missing).
    - Runs PyInstaller using windows\balanca_clp.spec
    - Output: windows\dist\BalancaCLP\BalancaCLP.exe

.NOTES
    Run from project root or any directory — script resolves paths relative to itself.
    Must be run on Windows (PyInstaller cannot cross-compile from Linux).
#>

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
$venvDir = Join-Path $scriptDir ".venv"
$snap7Dir = Join-Path $scriptDir "snap7"
$snap7Dll = Join-Path $snap7Dir "snap7.dll"
$distDir = Join-Path $projectRoot "dist"
$buildDir = Join-Path $projectRoot "build"
$specFile = Join-Path $scriptDir "balanca_clp.spec"
$reqFile = Join-Path $scriptDir "requirements-win.txt"

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg) { Write-Host "    $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "    $msg" -ForegroundColor Yellow }

# ---------------------------------------------------------------------------
# 1. Check Python
# ---------------------------------------------------------------------------
Write-Step "Checking Python..."
$python = Get-Command python -ErrorAction SilentlyContinue
if (-not $python) {
    throw "Python not found in PATH. Install Python 3.9 from https://www.python.org/downloads/ (on Server 2012, use Python 3.9 — newer versions may refuse to install)."
}
$pyVersion = & python --version 2>&1
Write-Ok "$pyVersion"

# ---------------------------------------------------------------------------
# 2. Check snap7.dll
# ---------------------------------------------------------------------------
Write-Step "Checking snap7.dll..."
if (-not (Test-Path $snap7Dll)) {
    New-Item -ItemType Directory -Force -Path $snap7Dir | Out-Null
    Write-Warn "snap7.dll not found at $snap7Dll"
    Write-Warn ""
    Write-Warn "Download steps:"
    Write-Warn "  1. Open https://sourceforge.net/projects/snap7/files/1.4.2/"
    Write-Warn "  2. Download snap7-full-1.4.2.7z"
    Write-Warn "  3. Extract release\Windows\Win64\snap7.dll (you need 7-Zip)"
    Write-Warn "  4. Copy snap7.dll to: $snap7Dll"
    Write-Warn ""
    throw "snap7.dll is required for the build. See instructions above."
}
Write-Ok "snap7.dll present"

# ---------------------------------------------------------------------------
# 3. Create / refresh venv
# ---------------------------------------------------------------------------
Write-Step "Setting up virtualenv at $venvDir..."
if (-not (Test-Path $venvDir)) {
    & python -m venv $venvDir
    if ($LASTEXITCODE -ne 0) { throw "Failed to create virtualenv" }
    Write-Ok "venv created"
} else {
    Write-Ok "venv already exists"
}

$venvPython = Join-Path $venvDir "Scripts\python.exe"
$venvPyInst = Join-Path $venvDir "Scripts\pyinstaller.exe"

# ---------------------------------------------------------------------------
# 4. Install dependencies
# ---------------------------------------------------------------------------
Write-Step "Installing dependencies..."
& $venvPython -m pip install --upgrade pip --disable-pip-version-check
if ($LASTEXITCODE -ne 0) { throw "pip upgrade failed" }
& $venvPython -m pip install -r $reqFile --disable-pip-version-check
if ($LASTEXITCODE -ne 0) { throw "pip install failed" }
Write-Ok "Dependencies installed"

# ---------------------------------------------------------------------------
# 5. Run PyInstaller
# ---------------------------------------------------------------------------
Write-Step "Running PyInstaller (this takes a minute or two)..."
Push-Location $projectRoot
try {
    if (Test-Path $distDir) { Remove-Item -Recurse -Force $distDir }
    if (Test-Path $buildDir) { Remove-Item -Recurse -Force $buildDir }
    & $venvPyInst --clean --noconfirm $specFile
    if ($LASTEXITCODE -ne 0) { throw "PyInstaller failed" }
}
finally {
    Pop-Location
}

# ---------------------------------------------------------------------------
# 6. Verify output
# ---------------------------------------------------------------------------
$exePath = Join-Path $distDir "BalancaCLP\BalancaCLP.exe"
if (-not (Test-Path $exePath)) {
    throw "Build appears to have succeeded but $exePath does not exist."
}

Write-Step "Build complete"
Write-Ok "Executable: $exePath"
Write-Ok ""
Write-Ok "Next steps:"
Write-Ok "  1. Copy/edit .env and scales.json next to the exe (or in project root)"
Write-Ok "  2. Test manually:   & '$exePath'"
Write-Ok "  3. Install service: .\windows\install_service.ps1  (as Administrator)"
