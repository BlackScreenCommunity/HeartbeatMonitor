$targets = @(
    @{ GOOS="linux"; GOARCH="amd64" }
)

foreach ($target in $targets) {
    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    $output = "HeartBeatMonitor-$env:GOOS-$env:GOARCH"
    go build -o $output
    Write-Host "Built: $output"
}