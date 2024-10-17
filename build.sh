#!/bin/bash

set -e

rm -rf build

# Define the target operating systems and architectures and extension
TARGETS=(
    "linux/386/"
    "linux/amd64/"
    "linux/arm64/"
    "darwin/amd64/"
    "darwin/arm64/"
    "windows/386/.exe"
    "windows/amd64/.exe"
)

product_name="md"
version="1.0.0"
start_dir="${PWD}"

for target in "${TARGETS[@]}"; do
    # Split the target operating system and architecture into separate variables
    IFS='/' read -r os arch extension <<< "${target}"

    # Define the output directory
    output_dir="build/${os}_${arch}_${version}"
    output_name="${product_name}${extension}"

    # Create the output directory if it doesn't exist
    mkdir -p "${output_dir}"

    # Build the program for the target operating system and architecture
    GOOS="${os}" GOARCH="${arch}" go build -ldflags="-s -w" -trimpath -o "${output_dir}/${output_name}" .

    # Zip the executable file
    cd "${output_dir}"
    zip "build.zip" "${output_name}"
    mv "build.zip" "../${product_name}_${os}_${arch}_${version}.zip"

    cd "${start_dir}"

    # Print a message indicating that the build was successful
    echo "Build for ${os}/${arch} successful."
done
