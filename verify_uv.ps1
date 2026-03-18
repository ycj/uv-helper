# Load the Windows PowerShell profile (the one we just installed)
$profilePath = Join-Path $env:USERPROFILE 'Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1'
if (Test-Path $profilePath) { . $profilePath }
# Reset last_dir to force re-evaluation
$global:__uv_helper_last_dir = $null
Set-Location 'C:\Users\10854\Documents\testuv'
# Trigger the prompt function which runs uv-helper
prompt | Out-Null
Write-Output "=== RESULTS ==="
Write-Output "VIRTUAL_ENV=$env:VIRTUAL_ENV"
Write-Output "PATH_starts_with=$($env:PATH.Substring(0, [Math]::Min(80, $env:PATH.Length)))"
