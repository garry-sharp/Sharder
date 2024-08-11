#!/bin/bash

root_dir=""
current_dir=$(pwd)

while [[ ! -f "$current_dir/README.md" ]]; do
    current_dir=$(dirname "$current_dir")
done

root_dir="$current_dir"

build_dir="$root_dir/build"
src_dir="$root_dir/go"

# Compile Go code for every architecture
architectures=("amd64" "arm64" "arm")

mkdir -p "$build_dir"
app_name="sharder"
version=$(git describe --tags --always)

$(
    cd $src_dir
    #GOOS=linux GOARCH=amd64 go build -o $build_dir/$app_name-$version-linux-amd64 main.go

    # Windows (x86_64)
    #GOOS=windows GOARCH=amd64 go build -o $build_dir/$app_name-$version-windows-amd64.exe main.go

    # macOS (x86_64)
    GOOS=darwin GOARCH=amd64 go build -o $build_dir/$app_name-$version-darwin-amd64 main.go

    # Linux (ARM)
    #GOOS=linux GOARCH=arm GOARM=7 go build -o $build_dir/$app_name-$version-linux-arm main.go

    # Linux (ARM64)
    #GOOS=linux GOARCH=arm64 go build -o $build_dir/$app_name-$version-linux-arm64 main.go

    # Windows (ARM64)
    #GOOS=windows GOARCH=arm64 go build -o $build_dir/$app_name-$version-windows-arm64.exe main.go

    # macOS (ARM64)
    #GOOS=darwin GOARCH=arm64 go build -o $build_dir/$app_name-$version-darwin-arm64 main.go
)
