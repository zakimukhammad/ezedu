# EzEdu Production Build & Verification Script
$ErrorActionPreference = "Stop"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host " EzEdu -- Production Build Verification" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

$baseDir = Get-Location

# 1. Build Go Backend
Write-Host ""
Write-Host "[1/3] Building Go Backend Server Binary..." -ForegroundColor Yellow
Set-Location "$baseDir\backend"
& "C:\Program Files\Go\bin\go.exe" build -ldflags="-s -w" -o bin/server.exe ./cmd/server
if ($LASTEXITCODE -eq 0) {
    $item = Get-Item "bin/server.exe"
    $binSize = [math]::Round($item.Length / 1MB, 2)
    Write-Host "  OK: Go Backend Binary compiled successfully! ($binSize MB)" -ForegroundColor Green
} else {
    Write-Host "  FAIL: Go compilation failed!" -ForegroundColor Red
    exit 1
}

# 2. Build Astro Static Frontend
Write-Host ""
Write-Host "[2/3] Building Astro Frontend Static Pages (SSG)..." -ForegroundColor Yellow
Set-Location "$baseDir\frontend"
cmd /c "npm run build"
if ($LASTEXITCODE -eq 0) {
    $htmlCount = (Get-ChildItem -Path "dist" -Recurse -Filter "*.html").Count
    Write-Host "  OK: Astro SSG Build successful! ($htmlCount static pages generated in frontend/dist/)" -ForegroundColor Green
} else {
    Write-Host "  FAIL: Astro build failed!" -ForegroundColor Red
    exit 1
}

# 3. Summary
Set-Location $baseDir
Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host " PRODUCTION BUILD COMPLETE AND VERIFIED!" -ForegroundColor Green
Write-Host " Go Server Binary : backend/bin/server.exe" -ForegroundColor Gray
Write-Host " Static Bundle    : frontend/dist/" -ForegroundColor Gray
Write-Host " Deployment Files : deployment/Caddyfile, deployment/ezedu.service" -ForegroundColor Gray
Write-Host "=========================================" -ForegroundColor Cyan
