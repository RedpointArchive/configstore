param(
    [switch][bool]$generate
)

$Args = ""
if ($generate) {
    $Args = "-generate"
}

Push-Location $PSScriptRoot\server
try {
    Write-Output "Building server..."
    go build -mod vendor -o server.exe .
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }

    Write-Output "Running & testing server..."
    $env:CONFIGSTORE_GOOGLE_CLOUD_PROJECT_ID="configstore-test-001"
    $env:CONFIGSTORE_GRPC_PORT="13389"
    $env:CONFIGSTORE_HTTP_PORT="13390"
    $env:CONFIGSTORE_SCHEMA_PATH="schema.json"
    .\server.exe $Args
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }
} finally {
    Pop-Location
}