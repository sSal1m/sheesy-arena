param(
    [string]$ArduinoCli = 'arduino-cli',
    [string]$ConfigurationDirectory = (Split-Path -Parent $PSScriptRoot)
)
$ErrorActionPreference = 'Stop'
$repo = Split-Path -Parent $PSScriptRoot
foreach ($alliance in @('red','blue')) {
    $header = Join-Path $ConfigurationDirectory "firmware\local\$alliance\DeviceConfig.h"
    if (-not (Test-Path -LiteralPath $header)) { throw "Missing $alliance private configuration. Run New-ArduinoConfig.ps1 first." }
    $sketch = Join-Path $ConfigurationDirectory "firmware\local\build\$alliance\FuelEthernet"
    New-Item -ItemType Directory -Path $sketch -Force | Out-Null
    Copy-Item -LiteralPath (Join-Path $repo 'firmware\FuelEthernet\FuelEthernet.ino') -Destination $sketch
    Copy-Item -LiteralPath $header -Destination $sketch
    try {
        # Windows PowerShell treats native stderr warnings as ErrorRecords.
        # Keep compiler diagnostics visible and decide success by its exit code.
        $ErrorActionPreference = 'Continue'
        & $ArduinoCli compile --fqbn arduino:avr:uno --warnings all --output-dir (Join-Path $sketch 'output') $sketch 2>&1 | ForEach-Object { $_.ToString() }
        $compileExit = $LASTEXITCODE
    } finally { $ErrorActionPreference = 'Stop' }
    if ($compileExit -ne 0) { throw "$alliance firmware compilation failed" }
}
