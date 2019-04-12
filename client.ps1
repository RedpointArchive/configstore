param()

Push-Location $PSScriptRoot\client
try {
    Write-Output "Generating client..."
    $env:CONFIGSTORE_GOOGLE_CLOUD_PROJECT_ID = "configstore-test-001"
    $env:CONFIGSTORE_GRPC_PORT = "13389"
    $env:CONFIGSTORE_HTTP_PORT = "13390"
    $env:CONFIGSTORE_SCHEMA_PATH = "..\server\schema.json"
    ..\server\server.exe -generate | Out-File -Encoding UTF8 -Force -FilePath .\client.go
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }

    Write-Output "Running & testing client..."
    go.exe test -v
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }
}
finally {
    Pop-Location
}