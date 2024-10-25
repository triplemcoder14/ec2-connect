#!/bin/bash

# Add Go to the PATH
export PATH=$PATH:/usr/local/go/bin

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go could not be found. Please install Go."
    exit 1
fi

# Create a release directory
mkdir -p release

# Define versions and architectures
VERSION="v1.0.0"
PLATFORMS=("darwin/amd64" "darwin/arm64" "linux/386" "linux/amd64" "linux/arm64" "windows/386" "windows/amd64" "windows/arm64")

# Build for each platform
for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%/*}"
    ARCH="${PLATFORM#*/}"
    OUTPUT="release/ec2-ssh_${VERSION}-${OS}-${ARCH}"

    # Set GOOS and GOARCH for cross-compilation
    echo "Building for $OS/$ARCH..."
    GOOS="$OS" GOARCH="$ARCH" go build -o "$OUTPUT" ./initssh.go

    # Check if the binary was created successfully
    if [ ! -f "$OUTPUT" ]; then
        echo "Failed to create binary for $OS/$ARCH"
        exit 1
    fi

    # Package the binary
    tar -czvf "${OUTPUT}.tar.gz" -C release "$OUTPUT"
done

echo "Build completed!"

