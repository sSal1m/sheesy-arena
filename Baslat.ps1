$ErrorActionPreference = 'Stop'
Set-Location -LiteralPath $PSScriptRoot
$exe = Join-Path $PSScriptRoot 'cheesy-arena.exe'
if (-not (Test-Path -LiteralPath $exe)) {
    throw 'cheesy-arena.exe bulunamadi. Once go build komutunu calistirin.'
}
$running = Get-Process -Name 'cheesy-arena' -ErrorAction SilentlyContinue |
    Where-Object { $_.Path -eq $exe }
if (-not $running) {
    if (Get-NetTCPConnection -LocalPort 8080 -State Listen -ErrorAction SilentlyContinue) {
        throw '8080 portunu baska bir uygulama kullaniyor.'
    }
    $running = Start-Process -FilePath $exe -ArgumentList '-dev' -WorkingDirectory $PSScriptRoot -WindowStyle Hidden -RedirectStandardOutput (Join-Path $PSScriptRoot 'arena-stdout.log') -RedirectStandardError (Join-Path $PSScriptRoot 'arena-stderr.log') -PassThru
}
for ($attempt = 0; $attempt -lt 30; $attempt++) {
    try {
        $response = Invoke-WebRequest -Uri 'http://localhost:8080/' -UseBasicParsing -TimeoutSec 2
        if ($response.StatusCode -eq 200) {
            Start-Process 'http://localhost:8080/'
            exit 0
        }
    } catch {}
    if (-not (Get-Process -Id $running.Id -ErrorAction SilentlyContinue)) {
        throw 'Uygulama kapandi. arena-stderr.log dosyasini kontrol edin.'
    }
    Start-Sleep -Seconds 1
}
throw 'Uygulama zamaninda yanit vermedi. arena-stderr.log dosyasini kontrol edin.'
