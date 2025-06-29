# PowerShell script to auto-install from GitHub repository
try {
    Write-Host "🌐 Downloading configuration-002 repository from GitHub..."

    $tempZip = Join-Path $env:TEMP "configuration-002.zip"
    $extractDir = Join-Path $env:TEMP ("configuration-002-" + [guid]::NewGuid().ToString())
    $repoUrl = "https://github.com/PeterCullenBurbery/configuration-002/archive/refs/heads/main.zip"

    # Download the ZIP file
    Invoke-WebRequest -Uri $repoUrl -OutFile $tempZip -UseBasicParsing

    Write-Host "📂 Extracting ZIP to: $extractDir"
    Expand-Archive -Path $tempZip -DestinationPath $extractDir -Force

    # Detect the actual extracted folder (GitHub zips include repo-name-branch)
    $unzippedRoot = Get-ChildItem -Path $extractDir | Where-Object { $_.PSIsContainer } | Select-Object -First 1
    if (-not $unzippedRoot) {
        throw "❌ Failed to detect extracted folder inside $extractDir"
    }

    $repoPath = $unzippedRoot.FullName
    $exePath = Join-Path $repoPath "canonical-file-structure\go-programs\call_installer\call_installer.exe"

    if (-not (Test-Path $exePath)) {
        throw "❌ call_installer.exe not found at expected location: $exePath"
    }

    Write-Host "🚀 Running call_installer.exe..."
    & $exePath $repoPath
} catch {
    Write-Error "❌ Script failed: $_"
} finally {
    # Cleanup optional
    Write-Host "🧹 Temporary folder: $extractDir"
    # Remove-Item -Recurse -Force $extractDir
}
