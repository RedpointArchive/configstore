param()

Push-Location $PSScriptRoot
try {
    Write-Output "Generating client..."
    .\protoc-3.6.1-win32\bin\protoc.exe --go_out=plugins=grpc:testclient .\testclient\testclient.proto
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }
    Move-Item -Force testclient\testclient\testclient.pb.go testclient\configstoreExample\testclient.pb.go 
    Remove-Item -Recurse -Force testclient\testclient

    Write-Output "Running & testing client..."
    Push-Location testclient
    try {
        go test -v
        if ($LastExitCode -ne 0) {
            exit $LastExitCode
        }
    } finally {
        Pop-Location
    }
} finally {
    Pop-Location
}