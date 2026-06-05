# Install PromptX on Windows
param(
    [string]$InstallDir = "$env:LOCALAPPDATA\promptx"
)

$Repo = "luoowei/promptx"
$Version = if ($env:PROMPTX_VERSION) { $env:PROMPTX_VERSION } else { "latest" }

Write-Host "Downloading PromptX..." -ForegroundColor Blue

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$Binary = "px_windows_${Arch}.exe"

if ($Version -eq "latest") {
    $Url = "https://github.com/${Repo}/releases/latest/download/${Binary}"
} else {
    $Url = "https://github.com/${Repo}/releases/download/${Version}/${Binary}"
}

# Create install directory
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

# Download
$TempFile = Join-Path $env:TEMP "px_new.exe"
try {
    Invoke-WebRequest -Uri $Url -OutFile $TempFile -ErrorAction Stop
} catch {
    Write-Host "Download failed. Check your internet connection." -ForegroundColor Red
    exit 1
}

# Install
$TargetFile = Join-Path $InstallDir "px.exe"
if (Test-Path $TargetFile) {
    Remove-Item $TargetFile -Force
}

Move-Item $TempFile $TargetFile -Force

Write-Host "PromptX installed to: $TargetFile" -ForegroundColor Green

# Add to PATH if not already there
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
    Write-Host "Added to PATH. Restart your terminal for changes to take effect." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "  Next steps:" -ForegroundColor Cyan
Write-Host "  1. Set API key: `$env:OPENAI_API_KEY = 'sk-...'" -ForegroundColor White
Write-Host "  2. Try it out: px ask 'hello world'" -ForegroundColor White
Write-Host "  3. Interactive: px" -ForegroundColor White
Write-Host ""
Write-Host "  Star the repo: https://github.com/${Repo}" -ForegroundColor Yellow
