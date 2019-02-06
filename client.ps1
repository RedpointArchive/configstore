param()

Push-Location $PSScriptRoot
try {
    Write-Output "Generating client..."
    .\protoc-3.6.1-win32\bin\protoc.exe --go_out=plugins=grpc:testclient/configstoreExample .\testclient\testclient.proto
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }

    Write-Output "Running & testing client..."
    Push-Location testclient
    try {
        go run .
        if ($LastExitCode -ne 0) {
            exit $LastExitCode
        }
    } finally {
        Pop-Location
    }
} finally {
    Pop-Location
}