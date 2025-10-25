# Build script for AWS Lambda deployment (PowerShell)

Write-Host "Building Lambda function for Linux AMD64..." -ForegroundColor Green

# Set environment variables for cross-compilation
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

# Build the binary
Write-Host "Compiling Go binary..." -ForegroundColor Yellow
go build -tags lambda.norpc -o bootstrap cmd/lambda/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    
    # Create deployment package
    Write-Host "Creating deployment package..." -ForegroundColor Yellow
    
    if (Test-Path function.zip) {
        Remove-Item function.zip
    }
    
    Compress-Archive -Path bootstrap -DestinationPath function.zip -Force
    
    if (Test-Path function.zip) {
        $fileSize = (Get-Item function.zip).Length / 1MB
        Write-Host "Deployment package created: function.zip ($([math]::Round($fileSize, 2)) MB)" -ForegroundColor Green
    } else {
        Write-Host "Failed to create zip file" -ForegroundColor Red
        exit 1
    }
    
    # Clean up
    Remove-Item bootstrap
    Write-Host "Cleanup complete" -ForegroundColor Green
} else {
    Write-Host "Build failed" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "To deploy to AWS Lambda:" -ForegroundColor Cyan
Write-Host "1. aws lambda update-function-code --function-name PreSpeechLayer --zip-file fileb://function.zip" -ForegroundColor White
Write-Host "2. Or upload function.zip via AWS Console" -ForegroundColor White
