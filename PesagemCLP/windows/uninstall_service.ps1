<#
.SYNOPSIS
    Stops and removes the BalancaCLP Windows service.

.PARAMETER ServiceName
    Service name to remove (default: BalancaCLP).
#>
[CmdletBinding()]
param(
    [string]$ServiceName = "BalancaCLP"
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$nssm = Join-Path $scriptDir "nssm.exe"

$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    throw "This script must be run as Administrator."
}

if (-not (Test-Path $nssm)) {
    throw "nssm.exe not found at $nssm"
}

$svc = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
if (-not $svc) {
    Write-Host "Service $ServiceName does not exist."
    return
}

Write-Host "==> Stopping $ServiceName..." -ForegroundColor Cyan
& $nssm stop $ServiceName 2>&1 | Out-Null

Write-Host "==> Removing $ServiceName..." -ForegroundColor Cyan
& $nssm remove $ServiceName confirm | Out-Null

Write-Host "Service $ServiceName removed." -ForegroundColor Green
