echo "💽 Building app for linux..."
GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux .
echo "💽 Building app for darwin..."
GOOS=darwin GOARCH=amd64 go build -o bin/app-arm64-darwin .
