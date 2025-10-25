# PreSpeech Layer - AWS Lambda Deployment Guide

## Prerequisites

- AWS CLI installed and configured
- AWS SAM CLI installed (optional, for SAM deployment)
- Go 1.23+ installed

## Building for Lambda

### Using PowerShell (Windows):
```powershell
.\build-lambda.ps1
```

### Using Bash (WSL/Linux/Mac):
```bash
chmod +x build-lambda.sh
./build-lambda.sh
```

This will create `function.zip` ready for deployment.

## Deployment Options

### Option 1: AWS SAM (Recommended)

1. **Install AWS SAM CLI:**
```powershell
# Windows (Chocolatey)
choco install aws-sam-cli

# Or download from: https://aws.amazon.com/serverless/sam/
```

2. **Build:**
```bash
sam build
```

3. **Deploy (first time):**
```bash
sam deploy --guided
```

Follow the prompts:
- Stack Name: `prespeech-layer-stack`
- AWS Region: `us-east-1` (or your preferred region)
- Parameter OpenAIApiKey: `sk-...`
- Parameter SupabaseUrl: `https://xxx.supabase.co`
- Parameter SupabaseKey: `eyJ...`
- Confirm changes before deploy: `Y`
- Allow SAM CLI IAM role creation: `Y`
- Save arguments to configuration file: `Y`

4. **Deploy (subsequent times):**
```bash
sam build && sam deploy
```

### Option 2: Manual Lambda Deployment

1. **Build the function:**
```powershell
.\build-lambda.ps1
```

2. **Create Lambda function via AWS CLI:**
```bash
aws lambda create-function \
  --function-name PreSpeechLayer \
  --runtime provided.al2023 \
  --handler bootstrap \
  --architectures x86_64 \
  --role arn:aws:iam::YOUR_ACCOUNT_ID:role/lambda-execution-role \
  --zip-file fileb://function.zip \
  --timeout 300 \
  --memory-size 512 \
  --environment Variables="{OPENAI_API_KEY=sk-xxx,SUPABASE_URL=https://xxx.supabase.co,SUPABASE_KEY=eyJxxx,GIN_MODE=release}"
```

3. **Update existing function:**
```bash
aws lambda update-function-code \
  --function-name PreSpeechLayer \
  --zip-file fileb://function.zip
```

4. **Set environment variables:**
```bash
aws lambda update-function-configuration \
  --function-name PreSpeechLayer \
  --environment Variables="{OPENAI_API_KEY=sk-xxx,SUPABASE_URL=https://xxx.supabase.co,SUPABASE_KEY=eyJxxx,GIN_MODE=release}"
```

### Option 3: AWS Console Upload

1. Build the function using `build-lambda.ps1`
2. Go to AWS Lambda Console
3. Create new function:
   - Runtime: Custom runtime on Amazon Linux 2023
   - Architecture: x86_64
4. Upload `function.zip`
5. Set Handler to: `bootstrap`
6. Set Timeout to: 5 minutes (300 seconds)
7. Set Memory to: 512 MB
8. Add environment variables:
   - `OPENAI_API_KEY`
   - `SUPABASE_URL`
   - `SUPABASE_KEY`
   - `GIN_MODE=release`

## API Gateway Setup

### Using SAM
SAM automatically creates API Gateway. The endpoint URL is shown in the stack outputs.

### Manual Setup
1. Create new HTTP API in API Gateway
2. Add Lambda integration pointing to `PreSpeechLayer` function
3. Create routes:
   - `POST /generate` → Lambda
   - `GET /session` → Lambda
4. Deploy API to a stage (e.g., `prod`)

## Testing

Get your API endpoint:
```bash
# If using SAM
aws cloudformation describe-stacks \
  --stack-name prespeech-layer-stack \
  --query 'Stacks[0].Outputs[?OutputKey==`ApiEndpoint`].OutputValue' \
  --output text
```

Test the endpoint:
```bash
# Replace with your actual endpoint
curl -X POST https://YOUR_API_ID.execute-api.us-east-1.amazonaws.com/Prod/generate \
  -H "Content-Type: application/json" \
  -d '{"paper":"Your research paper text here..."}'
```

## Local Testing

To test locally before deploying:
```bash
# Run the regular server
go run main.go

# Or use SAM local
sam local start-api
```

## Monitoring

View logs:
```bash
# Via AWS CLI
aws logs tail /aws/lambda/PreSpeechLayer --follow

# Or via SAM
sam logs -n PreSpeechLayerFunction --tail
```

## Cost Optimization

- Lambda: ~$0.20 per 1M requests + compute time
- API Gateway: ~$1.00 per 1M requests
- Consider using Lambda Reserved Concurrency if you have consistent traffic
- Monitor CloudWatch Logs size (retention set to 7 days)

## Troubleshooting

### Timeout Errors
- Increase Lambda timeout (currently set to 5 minutes)
- Check OpenAI API response times

### Cold Start Issues
- Consider provisioned concurrency for production
- Use Lambda SnapStart (if available for custom runtime)

### Database Connection Issues
- Ensure Lambda has internet access (default) or VPC configuration
- Verify Supabase credentials and URL
- Check security group rules if using VPC

## Environment Variables Required

| Variable | Description | Example |
|----------|-------------|---------|
| `OPENAI_API_KEY` | OpenAI API Key | `sk-...` |
| `SUPABASE_URL` | Supabase Project URL | `https://xxx.supabase.co` |
| `SUPABASE_KEY` | Supabase API Key | `eyJ...` |
| `GIN_MODE` | Gin framework mode | `release` |

## API Endpoints

Once deployed, your endpoints will be:

- `POST https://YOUR_API_ID.execute-api.REGION.amazonaws.com/Prod/generate`
- `GET https://YOUR_API_ID.execute-api.REGION.amazonaws.com/Prod/session?id=1`
