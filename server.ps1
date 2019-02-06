param()

Push-Location $PSScriptRoot
try {
    Write-Output "Building server..."
    go build -o testclient\testserver.exe .
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }

    Write-Output "Running & testing server..."
    .\testclient\testserver.exe
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }
} finally {
    Pop-Location
}