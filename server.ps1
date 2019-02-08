param()

Push-Location $PSScriptRoot
try {
    Write-Output "Building server..."
    go build -mod vendor -o testclient\testserver.exe .
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }

    Write-Output "Running & testing server..."
    $env:CONFIGSTORE_GOOGLE_CLOUD_PROJECT_ID="configstore-test-001"
    $env:CONFIGSTORE_GRPC_PORT="13389"
    $env:CONFIGSTORE_HTTP_PORT="13390"
    .\testclient\testserver.exe
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }
} finally {
    Pop-Location
}