$targets = @(
    @{ GOOS="linux"; GOARCH="amd64" }
)

# Function: Copy-FilesByExtension
# This function finds and copies all files with a given extension from one folder to another.
# It avoids copying files that are already in the destination folder.
function Copy-FilesByExtension {
    param (
        [string]$sourceFolder = "./internal",    # The source folder (default: current directory)
        [string]$destinationFolder = "./templates",  # The destination folder (default: ./templates)
        [string]$fileExtension = "*.css" # The file extension to search for (default: .css)
    )

    # Convert relative paths to absolute paths
    $sourceFolder = (Resolve-Path $sourceFolder).Path
    $destinationFolder = (Resolve-Path $destinationFolder).Path

    # Find all files with the given extension and copy them to the destination folder
    
    Get-ChildItem -Path $sourceFolder -Directory
        | Where-Object { $_.FullName -ne $destinationFolder }
        | Get-ChildItem -File -Filter $fileExtension -Recurse  | ForEach-Object {

        $filePath = (Resolve-Path $_.FullName).Path
        
        # Check if the file is NOT already in the destination folder
        if ($filePath -notlike "$destinationFolder\*") {
            Copy-Item -Path $filePath -Destination $destinationFolder -Force
        }
    }
}


foreach ($target in $targets) {
    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    $output = "HeartBeatMonitor-$env:GOOS-$env:GOARCH"
    go build -o ./Release/$env:GOOS/$output
    Write-Host "Built: $output"
    Copy-Item -Recurse -Force ./templates ./Release/$env:GOOS/
    Copy-Item -Force .\appsettings.json ./Release/$env:GOOS/
    
    Copy-FilesByExtension -sourceFolder "./internal" -destinationFolder ./Release/$env:GOOS/templates -fileExtension "*.css"

    Compress-Archive -Force -Path ./Release/$env:GOOS\* -CompressionLevel Fastest -DestinationPath .\Release\$output.zip
    Remove-Item -Force -Recurse ./Release/$env:GOOS/
}