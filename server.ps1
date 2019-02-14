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
    if (!$SkipDeps) {
        Write-Output "Fetch deps..."
        go get -u github.com/golang/protobuf/protoc-gen-go
        if ($LastExitCode -ne 0) {
            exit $LastExitCode
        }
    }

    Write-Output "Generate meta.go..."
    .\protoc.exe --go_out=plugins=grpc:meta .\meta.proto
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }
    $Content = Get-Content -Raw -Path .\meta\meta.pb.go
    $Content = $Content.Replace("package meta", "package main")
    Set-Content -Path .\meta\meta.pb.go -Value $Content
    Move-Item -Path .\meta\meta.pb.go -Destination .\meta.pb.go -Force

    Write-Output "Generate TypeScript gRPC client..."
    Push-Location .\typescript
    try {
        yarn
        if ($LastExitCode -ne 0) {
            exit $LastExitCode
        }
        if (!(Test-Path src\api)) {
            New-Item -ItemType Directory -Path src\api | Out-Null
        }
        ..\protoc.exe `
            --plugin="protoc-gen-ts=.\node_modules\.bin\protoc-gen-ts.cmd" `
            --js_out="import_style=commonjs,binary:src/api" `
            --ts_out="service=true:src/api" `
            -I .. `
            meta.proto
        if ($LastExitCode -ne 0) {
            exit $LastExitCode
        }
        $Content = Get-Content -Raw -Path .\src\api\meta_pb.js
        $Content = $Content.Replace("// GENERATED CODE", "/* eslint-disable */`n// GENERATED CODE")
        Set-Content -Path .\src\api\meta_pb.js -Value $Content
    } finally {
        Pop-Location
    }

    Write-Output "Building server..."
    go build -mod vendor -o server.exe .
    if ($LastExitCode -ne 0) {
        exit $LastExitCode
    }

    if (!$BuildOnly) {
        Write-Output "Running & testing server..."
        $env:CONFIGSTORE_GOOGLE_CLOUD_PROJECT_ID="configstore-test-001"
        $env:CONFIGSTORE_GRPC_PORT="13389"
        $env:CONFIGSTORE_HTTP_PORT="13390"
        $env:CONFIGSTORE_SCHEMA_PATH="schema.json"
        .\server.exe $Args
        if ($LastExitCode -ne 0) {
            exit $LastExitCode
        }
    }
} finally {
    Pop-Location
}