# WhatsApp Bot Quick Start Script
# This script sets up and runs the example bot

Write-Host "🚀 WhatsApp Bot Quick Start" -ForegroundColor Green
Write-Host ""

# Check if we're in the right directory
if (-not (Test-Path "example")) {
    Write-Host "❌ Error: Please run this script from the whatsmeow root directory" -ForegroundColor Red
    exit 1
}

# Navigate to example directory
Set-Location example

Write-Host "📦 Installing dependencies..." -ForegroundColor Cyan
go mod download

Write-Host ""
Write-Host "🔨 Building example..." -ForegroundColor Cyan
go build

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "✅ Build successful!" -ForegroundColor Green
Write-Host ""
Write-Host "📱 Starting WhatsApp bot..." -ForegroundColor Cyan
Write-Host "   If this is your first time, scan the QR code with WhatsApp" -ForegroundColor Yellow
Write-Host "   Press Ctrl+C to stop the bot" -ForegroundColor Yellow
Write-Host ""

# Run the example
.\example.exe

Write-Host ""
Write-Host "👋 Bot stopped" -ForegroundColor Yellow
