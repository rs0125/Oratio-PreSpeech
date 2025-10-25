#!/bin/bash

# Build script for AWS Lambda deployment

echo "Building Lambda function for Linux AMD64..."

# Set environment variables for cross-compilation
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# Build the binary
go build -tags lambda.norpc -o bootstrap cmd/lambda/main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    
    # Create deployment package
    echo "Creating deployment package..."
    zip -j function.zip bootstrap
    
    if [ $? -eq 0 ]; then
        echo "Deployment package created: function.zip"
        echo "File size: $(du -h function.zip | cut -f1)"
    else
        echo "Failed to create zip file"
        exit 1
    fi
    
    # Clean up
    rm bootstrap
    echo "Cleanup complete"
else
    echo "Build failed"
    exit 1
fi

echo ""
echo "To deploy to AWS Lambda:"
echo "1. aws lambda update-function-code --function-name PreSpeechLayer --zip-file fileb://function.zip"
echo "2. Or upload function.zip via AWS Console"
