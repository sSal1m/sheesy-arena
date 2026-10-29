param([string]$OutputDirectory = (Split-Path -Parent $PSScriptRoot), [ValidateRange(1,2)][int]$SensorChannels = 1)
$ErrorActionPreference = 'Stop'
$destination = [IO.Path]::GetFullPath($OutputDirectory)
$serverFile = Join-Path $destination 'arduino.local.json'
if (Test-Path -LiteralPath $serverFile) { throw 'Existing Arduino keys will not be overwritten. Use a separate output directory for a new installation.' }
New-Item -ItemType Directory -Path $destination -Force | Out-Null
$devices = @()
$rng = [Security.Cryptography.RandomNumberGenerator]::Create()
try {
    foreach ($alliance in @('red','blue')) {
        $suffix = if ($alliance -eq 'red') {201} else {202}
        $key = New-Object byte[] 32
        $rng.GetBytes($key)
        $keyHex = ($key | ForEach-Object { $_.ToString('x2') }) -join ''
        $id = 'fuel-' + $alliance + '-01'
        $devices += [ordered]@{Id=$id;Alliance=$alliance;IP="10.0.100.$suffix";Key=$keyHex;Source='arduino'}
        $headerDirectory = Join-Path $destination "firmware\local\$alliance"
        New-Item -ItemType Directory -Path $headerDirectory -Force | Out-Null
        $keyBytes = ($key | ForEach-Object { '0x' + $_.ToString('x2') }) -join ','
        @"
#pragma once
// PRIVATE: generated device identity/key. Do not commit or share this file.
const char DEVICE_ID[] = "$id";
const uint8_t DEVICE_KEY[32] = {$keyBytes};
const uint8_t DEVICE_MAC[6] = {0x02,0x54,0x26,0x00,0x00,$suffix};
const uint8_t DEVICE_IP[4] = {10,0,100,$suffix};
const uint8_t SERVER_IP[4] = {10,0,100,5};
const uint8_t GATEWAY_IP[4] = {10,0,100,3};
const uint8_t SUBNET[4] = {255,255,255,0};
const uint16_t SERVER_PORT = 1760;
const uint8_t ETHERNET_CS_PIN = 10;
const uint8_t SD_CS_PIN = 4;
// Two independent beam channels on UNO external interrupts D2 and D3.
const uint8_t SENSOR_PINS[2] = {2,3};
const uint8_t SENSOR_CHANNELS = $SensorChannels;
const bool SENSOR_ACTIVE_LOW = true;
const bool SENSOR_PULLUP = true;
const uint32_t SENSOR_DEBOUNCE_US = 1000;
"@ | Set-Content -LiteralPath (Join-Path $headerDirectory 'DeviceConfig.h') -Encoding ascii
    }
} finally { $rng.Dispose() }
# Disabled until the operator has checked the hardware, IPs and wiring.
[ordered]@{Enabled=$false;Listen='10.0.100.5:1760';Devices=$devices} | ConvertTo-Json -Depth 5 | Set-Content -LiteralPath $serverFile -Encoding ascii
Write-Output "Generated disabled server configuration and private RED/BLUE headers in $destination. Review before enabling."
