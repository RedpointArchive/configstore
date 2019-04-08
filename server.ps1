param(
    [switch][bool]$generate,
    [switch][bool]$SkipDeps,
    [switch][bool]$BuildOnly
)

$Args = ""
if ($generate) {
    $Args = "-generate"
}

Push-Location $PSScriptRoot\server
try {
    Write-Output "Building protocol buffers..."
    $env:DOCKER_BUILDKIT = 1
    docker.exe build .. -f ../Dockerfile --target=protocol_build --tag=configstore_protocol_build
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }

    $DockerContainerId = $null
    try {
        $DockerContainerId = $(docker.exe create configstore_protocol_build)
        if (Test-Path $PSScriptRoot/workdir_ts/) {
            Remove-Item $PSScriptRoot/workdir_ts/ -Force -Recurse
        }
        if (Test-Path $PSScriptRoot/workdir_go/) {
            Remove-Item $PSScriptRoot/workdir_go/ -Force -Recurse
        }
        docker.exe cp ${DockerContainerId}:/workdir_ts/ $PSScriptRoot/workdir_ts/
        docker.exe cp ${DockerContainerId}:/workdir_go/ $PSScriptRoot/workdir_go/
        Start-Sleep -Seconds 1
        Copy-Item $PSScriptRoot/workdir_go/meta.pb.go $PSScriptRoot/server/meta.pb.go -Force
        if (Test-Path $PSScriptRoot/server-ui/src/api) {
            Remove-Item $PSScriptRoot/server-ui/src/api -Recurse -Force
            Start-Sleep -Seconds 1
        }
        Copy-Item -Recurse $PSScriptRoot/workdir_ts/api $PSScriptRoot/server-ui/src/api -Force
        if (Test-Path $PSScriptRoot/workdir_ts/) {
            Remove-Item $PSScriptRoot/workdir_ts/ -Force -Recurse
        }
        if (Test-Path $PSScriptRoot/workdir_go/) {
            Remove-Item $PSScriptRoot/workdir_go/ -Force -Recurse
        }
    }
    finally {
        docker.exe rm $DockerContainerId
    }

    Write-Output "Building server..."
    go.exe build -mod vendor -o server.exe .
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }

    if (!$BuildOnly) {
        Write-Output "Running & testing server..."
        $env:CONFIGSTORE_GOOGLE_CLOUD_PROJECT_ID = "configstore-test-001"
        $env:CONFIGSTORE_GRPC_PORT = "13389"
        $env:CONFIGSTORE_HTTP_PORT = "13390"
        $env:CONFIGSTORE_SCHEMA_PATH = "schema.json"
        .\server.exe $Args
        if ($LastExitCode -ne 0) {
            exit $LastExitCode
        }
    }
}
finally {
    Pop-Location
}