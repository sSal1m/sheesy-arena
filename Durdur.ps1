$exe = Join-Path $PSScriptRoot 'cheesy-arena.exe'
Get-Process -Name 'cheesy-arena' -ErrorAction SilentlyContinue |
    Where-Object { $_.Path -eq $exe } |
    Stop-Process
