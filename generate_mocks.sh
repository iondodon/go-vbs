#!/bin/bash

# Loop through all Go files in all subdirectories
for file in $(find . -name "*.go" | grep -v "mocks_test"); do
    # Derive the package name and mock destination directory
    pkg=$(dirname $file)
    mockDir="${pkg}/mocks_test"

    # Create the mock directory if it doesn't exist
    mkdir -p $mockDir

    # Generate the mock
    mockgen -source=$file -destination=$mockDir/mock_$(basename $file)
done