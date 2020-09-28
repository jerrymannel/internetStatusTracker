env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o internetStatusTracker-pi .
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o internetStatusTracker-macos .