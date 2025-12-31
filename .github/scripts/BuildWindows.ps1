# .github/scripts/BuildWindows.ps1
param (
    [string]$BuildDir = "build",
    [string]$ProjectName = "NaxiServer",
    [string]$Version = "v1.0.0"
)

Write-Host "--- Processing Version ---"
$RegexPattern = "\d+\.\d+\.\d+"
$Match = [regex]::Match($Version, $RegexPattern)
if ($Match.Success) {
    $CleanVersion = $Match.Value
} else {
    $CleanVersion = "1.0.0"
}

if (Test-Path $BuildDir) { Remove-Item -Recurse -Force $BuildDir }
New-Item -ItemType Directory -Path $BuildDir

Write-Host "--- Compiling ---"
go build -ldflags="-s -w" -o $BuildDir/$ProjectName/$ProjectName.exe

Write-Host "--- Preparing Deploy ---"
Copy-Item "COPYING" -Destination "$BuildDir\$ProjectName" -Force
New-Item -ItemType Directory -Path "$BuildDir\db"
Copy-Item "dist\db\**" -Destination "$BuildDir\db" -Recurse -Force

Write-Host "--- Creating a ZIP File ---"
$ZipName = "${ProjectName}-Windows-x86_64.zip"
Compress-Archive -Path "$BuildDir\*" -DestinationPath $ZipName -Force
Write-Host "Build and Packaging completed: $ZipName"


Write-Host "--- Complete build of $ProjectName $CleanVersion"