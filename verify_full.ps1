# Load the Windows PowerShell profile
$profilePath = Join-Path $env:USERPROFILE 'Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1'
if (Test-Path $profilePath) { . $profilePath }

# Step 1: Enter testuv → should activate
$global:__uv_helper_last_dir = $null
Set-Location 'C:\Users\10854\Documents\testuv'
prompt | Out-Null
Write-Output "=== AFTER cd testuv ==="
Write-Output "VIRTUAL_ENV=$env:VIRTUAL_ENV"

# Step 2: Leave to a non-uv directory → should deactivate
Set-Location 'C:\Users\10854\Documents'
prompt | Out-Null
Write-Output "=== AFTER cd Documents ==="
Write-Output "VIRTUAL_ENV=$env:VIRTUAL_ENV"

# Step 3: Re-enter testuv → should activate again
Set-Location 'C:\Users\10854\Documents\testuv'
prompt | Out-Null
Write-Output "=== AFTER cd testuv again ==="
Write-Output "VIRTUAL_ENV=$env:VIRTUAL_ENV"
