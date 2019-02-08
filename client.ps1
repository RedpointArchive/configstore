param()

Push-Location $PSScriptRoot\client
try {
    Write-Output "Fetch deps..."
    Push-Location $env:SYSTEMDRIVE
    go get -u github.com/golang/protobuf/proto
    Pop-Location

    Write-Output "Generating client..."
    ..\server\server.exe -generate | Out-File -Encoding UTF8 -Force -FilePath .\client.go
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }

    Write-Output "Running & testing client..."
    go test -v
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }
} finally {
    Pop-Location
}