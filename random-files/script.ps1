# PowerShell script to auto-install from GitHub repository
try {
    Write-Host "🌐 Downloading configuration-002 repository from GitHub..."

    # Set the temp file path for the ZIP file
    $tempZip = Join-Path $env:TEMP "configuration-002.zip"

    # Generate a unique extraction directory under TEMP
    $extractDir = Join-Path $env:TEMP ("configuration-002-" + [guid]::NewGuid().ToString())

    # GitHub ZIP URL for the main branch
    $repoUrl = "https://github.com/PeterCullenBurbery/configuration-002/archive/refs/heads/main.zip"

    # Optional: Ensure TLS 1.2 is used on PowerShell 5.1
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

    # Use WebClient for a faster download than Invoke-WebRequest
    $webClient = New-Object System.Net.WebClient
    $webClient.DownloadFile($repoUrl, $tempZip)

    Write-Host "📂 Extracting ZIP to: $extractDir"

    # Unzip the downloaded repository into the extraction directory
    Expand-Archive -Path $tempZip -DestinationPath $extractDir -Force

    # Find the root folder of the extracted repository (GitHub zips include repo-name-branch)
    $unzippedRoot = Get-ChildItem -Path $extractDir | Where-Object { $_.PSIsContainer } | Select-Object -First 1

    # If the extracted folder wasn't found, raise an error
    if (-not $unzippedRoot) {
        throw "❌ Failed to detect extracted folder inside $extractDir"
    }

    # Store the full path to the root of the unzipped repository
    $repoPath = $unzippedRoot.FullName

    # Compute the full path to the expected call_installer.exe location
    $exePath = Join-Path $repoPath "canonical-file-structure\go-programs\call_installer\call_installer.exe"

    # Validate that call_installer.exe exists at the expected path
    if (-not (Test-Path $exePath)) {
        throw "❌ call_installer.exe not found at expected location: $exePath"
    }

    Write-Host "🚀 Running call_installer.exe..."

    # Execute call_installer.exe and pass the repository root path as an argument
    & $exePath $repoPath

} catch {
    # Log any errors encountered during the process
    Write-Error "❌ Script failed: $_"
} finally {
    # Output the temporary folder path (optional cleanup can be enabled below)
    Write-Host "🧹 Temporary folder: $extractDir"

    # Optional cleanup: uncomment to remove temp directory
    # Remove-Item -Recurse -Force $extractDir
}
