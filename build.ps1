$targets = @(
    @{ GOOS="linux"; GOARCH="amd64" }
)

foreach ($target in $targets) {
    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    $output = "HeartBeatMonitor-$env:GOOS-$env:GOARCH"
    go build -o ./Release/$env:GOOS/$output
    Write-Host "Built: $output"
    Copy-Item -Recurse -Force ./templates ./Release/$env:GOOS/
    Compress-Archive -Path ./Release/$env:GOOS\* -CompressionLevel Fastest -DestinationPath .\Release\$output.zip
    Remove-Item -Force -Recurse ./Release/$env:GOOS/
}