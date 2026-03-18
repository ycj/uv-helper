$end=(Get-Date).AddMinutes(10)
while($true) {
  $j = & 'C:\Program Files\GitHub CLI\gh.exe' run view 23227262627 --json number,status,conclusion,url,workflowName | ConvertFrom-Json
  Write-Host "[$(Get-Date -Format HH:mm:ss)] $($j.status) ($($j.conclusion))"
  if ($j.status -ne 'in_progress') { break }
  if ((Get-Date) -ge $end) { Write-Host 'Timed out'; exit 2 }
  Start-Sleep -Seconds 5
}
& 'C:\Program Files\GitHub CLI\gh.exe' run view 23227262627 --log
