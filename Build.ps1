# Build for Windows (64-bit)
Write-Host "Building HolyQuestGo.exe (Windows)"
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o HolyQuestGo.exe

# Build for Linux (64-bit)
Write-Host "Building HolyQuestGo.bin (Linux)"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o HolyQuestGo.bin

Write-Host "Builds complete: HolyQuestGo.exe (Windows), HolyQuestGo.bin (Linux)"