<#
.SYNOPSIS
    Installs BalancaCLP as a Windows service using NSSM.

.DESCRIPTION
    - Requires nssm.exe in windows\nssm.exe (download from https://nssm.cc/download).
    - Requires the .exe to already exist (run build.ps1 first).
    - Reads .env in project root, passes vars to the service environment.
    - Configures auto-start, restart on failure, log rotation.

.PARAMETER ServiceName
    Service name (default: BalancaCLP).

.PARAMETER User
    Optional user account to run the service as. Default: LocalSystem.

.NOTES
    Run as Administrator (right-click PowerShell -> Run as administrator).
#>
[CmdletBinding()]
param(
    [string]$ServiceName = "BalancaCLP",
    [string]$User = ""
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
$nssm = Join-Path $scriptDir "nssm.exe"
$exePath = Join-Path $projectRoot "dist\BalancaCLP\BalancaCLP.exe"
$workDir = Join-Path $projectRoot "dist\BalancaCLP"
$logDir = Join-Path $projectRoot "logs"
$envFile = Join-Path $projectRoot ".env"
$scalesJson = Join-Path $projectRoot "scales.json"

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg) { Write-Host "    $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "    $msg" -ForegroundColor Yellow }

# ---------------------------------------------------------------------------
# Preflight checks
# ---------------------------------------------------------------------------
Write-Step "Preflight checks..."

$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    throw "This script must be run as Administrator."
}

if (-not (Test-Path $nssm)) {
    throw @"
nssm.exe not found at $nssm

Download steps:
  1. Open https://nssm.cc/download
  2. Download nssm-2.24 (or newer)
  3. Extract win64\nssm.exe
  4. Copy to: $nssm
"@
}

if (-not (Test-Path $exePath)) {
    throw "Executable not found at $exePath. Run windows\build.ps1 first."
}

if (-not (Test-Path $scalesJson)) {
    throw "scales.json not found at $scalesJson. Create it (see scales.json.example)."
}

if (-not (Test-Path $envFile)) {
    Write-Warn ".env not found at $envFile — service will only have OS-level env vars."
}

New-Item -ItemType Directory -Force -Path $logDir | Out-Null
Write-Ok "All checks passed"

# ---------------------------------------------------------------------------
# Stop & remove any existing service with the same name
# ---------------------------------------------------------------------------
Write-Step "Checking for existing service..."
$svc = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
if ($svc) {
    Write-Warn "Service $ServiceName exists — stopping and removing."
    & $nssm stop $ServiceName 2>&1 | Out-Null
    & $nssm remove $ServiceName confirm | Out-Null
    Start-Sleep -Seconds 1
}

# ---------------------------------------------------------------------------
# Install service
# ---------------------------------------------------------------------------
Write-Step "Installing service $ServiceName..."
& $nssm install $ServiceName $exePath
if ($LASTEXITCODE -ne 0) { throw "nssm install failed" }

& $nssm set $ServiceName AppDirectory $workDir | Out-Null
& $nssm set $ServiceName DisplayName "Balanca CLP Integration" | Out-Null
& $nssm set $ServiceName Description "Integra balancas via CLP Siemens (Snap7) com SQL Server. UI em http://localhost:8080" | Out-Null
& $nssm set $ServiceName Start SERVICE_AUTO_START | Out-Null
& $nssm set $ServiceName DependOnService Tcpip | Out-Null

# Logs
$stdoutLog = Join-Path $logDir "balanca-stdout.log"
$stderrLog = Join-Path $logDir "balanca-stderr.log"
& $nssm set $ServiceName AppStdout $stdoutLog | Out-Null
& $nssm set $ServiceName AppStderr $stderrLog | Out-Null
& $nssm set $ServiceName AppRotateFiles 1 | Out-Null
& $nssm set $ServiceName AppRotateOnline 1 | Out-Null
& $nssm set $ServiceName AppRotateSeconds 86400 | Out-Null   # rotate daily
& $nssm set $ServiceName AppRotateBytes 10485760 | Out-Null  # or after 10 MB

# Restart on failure
& $nssm set $ServiceName AppExit Default Restart | Out-Null
& $nssm set $ServiceName AppRestartDelay 5000 | Out-Null     # 5 s
& $nssm set $ServiceName AppThrottle 10000 | Out-Null

# Optional service account
if ($User) {
    Write-Step "Configuring to run as $User (you'll be prompted for password by NSSM)..."
    & $nssm set $ServiceName ObjectName $User
}

# ---------------------------------------------------------------------------
# Environment variables
# ---------------------------------------------------------------------------
Write-Step "Loading environment variables..."
$envLines = @()
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        $line = $_.Trim()
        if ($line -and -not $line.StartsWith('#') -and $line.Contains('=')) {
            $envLines += $line
        }
    }
}
# Force SCALES_CONFIG_PATH to the project root scales.json (not the one inside the bundle)
$envLines = $envLines | Where-Object { -not $_.StartsWith("SCALES_CONFIG_PATH=") }
$envLines += "SCALES_CONFIG_PATH=$scalesJson"

if ($envLines.Count -gt 0) {
    # NSSM accepts multiple KEY=VALUE entries as positional args
    & $nssm set $ServiceName AppEnvironmentExtra $envLines | Out-Null
    Write-Ok "$($envLines.Count) environment variables loaded"
}

# ---------------------------------------------------------------------------
# Start
# ---------------------------------------------------------------------------
Write-Step "Starting service..."
& $nssm start $ServiceName
Start-Sleep -Seconds 2

$svc = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
if ($svc -and $svc.Status -eq "Running") {
    Write-Ok "Service $ServiceName is RUNNING"
    Write-Ok "  Logs:    $logDir"
    Write-Ok "  Web UI:  http://localhost:8080  (and http://<server-ip>:8080)"
    Write-Ok "  Manage:  services.msc  (or 'nssm edit $ServiceName' as Admin)"
} else {
    Write-Warn "Service may have failed to start. Check logs:"
    Write-Warn "  $stderrLog"
    if ($svc) { Write-Warn "  Status: $($svc.Status)" }
}
