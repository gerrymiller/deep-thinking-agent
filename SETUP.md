# Setup Guide

Complete installation and configuration guide for Deep Thinking Agent.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Go Installation](#go-installation)
3. [Qdrant Installation](#qdrant-installation)
4. [OpenAI API Setup](#openai-api-setup)
5. [Building the CLI](#building-the-cli)
6. [Configuration](#configuration)
7. [Verification](#verification)
8. [Troubleshooting](#troubleshooting)

## Prerequisites

Before starting, ensure you have:

- **Go 1.25.3 or higher** - For building the application
- **Docker** (recommended) or access to Qdrant Cloud - For vector database
- **OpenAI API account** - For LLM and embedding operations
- **~8GB RAM** - Recommended for running Qdrant and processing
- **Internet connection** - Required for API calls

### Quick Prerequisites Check

```bash
# Check Go version
go version  # Should show go1.25.3 or higher

# Check Docker
docker --version  # Should show Docker version info

# Check internet connectivity
curl -I https://api.openai.com  # Should return 200 OK
```

## Go Installation

### Verify Current Installation

```bash
go version
```

If you see `go version go1.25.3` or higher, you're ready. Otherwise, install or upgrade:

### Installation by Platform

#### macOS

```bash
# Using Homebrew
brew install go

# Or download from golang.org
curl -OL https://go.dev/dl/go1.25.3.darwin-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.3.darwin-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

#### Linux

```bash
# Download and install
curl -OL https://go.dev/dl/go1.25.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.3.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc for persistence)
export PATH=$PATH:/usr/local/go/bin
```

#### Windows

Download the installer from [golang.org/dl](https://go.dev/dl/) and run it.

### Verify Installation

```bash
go version
go env GOPATH  # Should show your Go workspace path
```

## Qdrant Installation

Choose one of three options based on your needs:

### Option A: Docker (Recommended for Testing)

**Advantages:** Quick setup, easy cleanup, isolated environment

**Installation:**

```bash
# Start Qdrant with persistent storage
docker run -d \
  --name qdrant \
  -p 6333:6333 \
  -p 6334:6334 \
  -v $(pwd)/qdrant_storage:/qdrant/storage \
  qdrant/qdrant

# Verify it's running
docker ps | grep qdrant
curl http://localhost:6333/collections
```

**Managing Qdrant:**

```bash
# Check status
docker ps | grep qdrant

# View logs
docker logs qdrant

# Stop Qdrant
docker stop qdrant

# Start existing container
docker start qdrant

# Remove container
docker rm -f qdrant

# Remove stored data
rm -rf ./qdrant_storage
```

### Option B: Qdrant Cloud (Recommended for Production)

**Advantages:** Managed service, no local resources, automatic backups

**Setup:**

1. Go to [cloud.qdrant.io](https://cloud.qdrant.io)
2. Create a free account
3. Create a new cluster
4. Note your cluster URL and API key
5. Update your `config.json`:

```json
{
  "vector_store": {
    "type": "qdrant",
    "address": "https://your-cluster.qdrant.io:6333",
    "api_key": "your-api-key-here",
    "default_collection": "documents"
  }
}
```

### Option C: Local Binary

**Advantages:** No Docker required, native performance

**Installation:**

```bash
# macOS/Linux
curl -sSL https://github.com/qdrant/qdrant/releases/latest/download/qdrant-x86_64-apple-darwin.tar.gz | tar xz
./qdrant &

# Or build from source
git clone https://github.com/qdrant/qdrant.git
cd qdrant
cargo build --release
./target/release/qdrant &
```

### Verify Qdrant Installation

```bash
# Check HTTP endpoint
curl http://localhost:6333/collections

# Should return:
# {"result":{"collections":[]},"status":"ok","time":0.000123}

# Check gRPC endpoint (used by the CLI)
# If this fails, check that port 6334 is accessible
telnet localhost 6334
```

## OpenAI API Setup

### Getting an API Key

1. **Create an OpenAI Account**
   - Go to [platform.openai.com](https://platform.openai.com)
   - Sign up or log in

2. **Generate API Key**
   - Navigate to [API Keys](https://platform.openai.com/api-keys)
   - Click "Create new secret key"
   - **Copy the key immediately** (you won't be able to see it again)
   - Store it securely (password manager recommended)

3. **Set Spending Limits** (Highly Recommended)
   - Go to [Billing Settings](https://platform.openai.com/account/billing/limits)
   - Set a monthly budget limit (e.g., $10 for testing)
   - Enable email notifications for usage alerts

### Setting the API Key

#### Using .env File (Recommended)

The easiest way to configure your API key is using a `.env` file:

```bash
# Copy the example template
cp .env.example .env

# Edit .env and add your API key
# The file should contain:
# OPENAI_API_KEY=sk-your-key-here
```

The CLI will automatically load `.env` and `.env.local` files when you run commands. You can also override specific settings using `.env.local` (useful for testing different configurations):

```bash
# .env - Your main configuration
OPENAI_API_KEY=sk-your-key-here

# .env.local - Local overrides (optional)
REASONING_LLM_MODEL=gpt-4o-mini
```

**Note**: Both `.env` and `.env.local` are already in `.gitignore` to protect your API keys.

#### Using Environment Variables (Alternative)

Alternatively, you can set environment variables directly:

**Linux/macOS:**

```bash
# Temporary (current session only)
export OPENAI_API_KEY="sk-your-key-here"

# Permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export OPENAI_API_KEY="sk-your-key-here"' >> ~/.bashrc
source ~/.bashrc
```

**Windows (PowerShell):**

```powershell
# Temporary
$env:OPENAI_API_KEY="sk-your-key-here"

# Permanent (system-wide)
[System.Environment]::SetEnvironmentVariable('OPENAI_API_KEY', 'sk-your-key-here', 'User')
```

**Windows (Command Prompt):**

```cmd
# Temporary
set OPENAI_API_KEY=sk-your-key-here

# Permanent
setx OPENAI_API_KEY "sk-your-key-here"
```

### Verify API Key

```bash
# If using .env file, the CLI loads it automatically
# You can verify with:
grep OPENAI_API_KEY .env

# If using environment variables, check it's set:
echo $OPENAI_API_KEY  # Should print your key

# Test API access (requires curl and jq)
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY" | jq '.data[0].id'
```

### Cost Warnings and Optimization

**Typical Costs (as of 2025):**
- GPT-4o: ~$0.005 per 1K input tokens, ~$0.015 per 1K output tokens
- GPT-4o-mini: ~$0.00015 per 1K input tokens, ~$0.0006 per 1K output tokens
- text-embedding-3-small: ~$0.00002 per 1K tokens

**Per Operation Estimates:**
- Document ingestion: $0.012-0.052 per document
- Simple query: $0.06-0.15
- Complex multi-hop query: $0.15-0.31
- Running all examples: ~$2-5 total

**Cost Reduction Tips:**
- Use `--no-schema` flag during testing to skip LLM-based schema analysis
- Start with smaller documents
- Use `gpt-4o-mini` for all LLM operations (edit config.json)
- Set `max_iterations` lower in config (e.g., 3 instead of 10)
- Monitor usage at [platform.openai.com/usage](https://platform.openai.com/usage)

## Building the CLI

### Clone Repository

```bash
# Clone the repository
git clone https://github.com/yourusername/deep-thinking-agent.git
cd deep-thinking-agent

# Or if you have the source already
cd /path/to/deep-thinking-agent
```

### Download Dependencies

```bash
# Download all Go module dependencies
go mod download

# Verify dependencies
go mod verify
```

### Build the CLI Binary

```bash
# Build for your current platform
go build -o bin/deep-thinking-agent ./cmd/cli

# Verify the binary
./bin/deep-thinking-agent --version
```

### Platform-Specific Builds

```bash
# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o bin/deep-thinking-agent-macos ./cmd/cli

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/deep-thinking-agent-linux ./cmd/cli

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/deep-thinking-agent.exe ./cmd/cli
```

### Install Globally (Optional)

```bash
# Copy to system PATH
sudo cp bin/deep-thinking-agent /usr/local/bin/

# Or add bin directory to PATH
export PATH=$PATH:$(pwd)/bin
echo 'export PATH=$PATH:/path/to/deep-thinking-agent/bin' >> ~/.bashrc
```

## Configuration

### Initialize Default Configuration

```bash
# Create config.json with defaults
./bin/deep-thinking-agent config init

# This creates config.json in the current directory
```

### Configuration File Structure

The `config.json` file controls all system behavior:

```json
{
  "llm": {
    "reasoning_llm": {
      "provider": "openai",
      "model": "gpt-4o",
      "api_key": "${OPENAI_API_KEY}",
      "default_temperature": 0.7
    },
    "fast_llm": {
      "provider": "openai",
      "model": "gpt-4o-mini",
      "api_key": "${OPENAI_API_KEY}",
      "default_temperature": 0.5
    }
  },
  "embedding": {
    "provider": "openai",
    "model": "text-embedding-3-small",
    "api_key": "${OPENAI_API_KEY}"
  },
  "vector_store": {
    "type": "qdrant",
    "address": "localhost:6334",
    "default_collection": "documents"
  },
  "workflow": {
    "max_iterations": 10,
    "top_k_retrieval": 10,
    "top_n_reranking": 3,
    "default_strategy": "hybrid"
  }
}
```

### Configuration Options Explained

**LLM Configuration:**
- `reasoning_llm`: Used for complex tasks (planning, reflection, policy decisions)
  - Recommended: `gpt-4o` for quality, `gpt-4o-mini` for cost savings
- `fast_llm`: Used for quick tasks (rewriting, distillation, supervision)
  - Recommended: `gpt-4o-mini` for speed and cost
- `temperature`: Controls randomness (0.0 = deterministic, 1.0 = creative)
  - Reasoning: 0.7 (balanced)
  - Fast: 0.5 (more focused)

**Embedding Configuration:**
- `model`: Embedding model for vector search
  - Recommended: `text-embedding-3-small` (cost-effective, 1536 dimensions)
  - Alternative: `text-embedding-3-large` (better quality, 3072 dimensions, more expensive)

**Vector Store Configuration:**
- `address`: Qdrant server address
  - Local: `localhost:6334` (gRPC port)
  - Cloud: `https://your-cluster.qdrant.io:6333`
- `default_collection`: Default collection name for documents

**Workflow Configuration:**
- `max_iterations`: Maximum reasoning loop iterations (default: 10)
  - Lower = faster, cheaper, less thorough
  - Higher = slower, costlier, more thorough
- `top_k_retrieval`: Number of documents to retrieve initially (default: 10)
- `top_n_reranking`: Number of documents after reranking (default: 3)
- `default_strategy`: Retrieval strategy (`vector`, `keyword`, or `hybrid`)

### Environment Variable Overrides

Environment variables take precedence over config file:

```bash
# LLM settings
export REASONING_LLM_MODEL=gpt-4o-mini
export FAST_LLM_MODEL=gpt-4o-mini
export OPENAI_API_KEY=sk-your-key

# Vector store settings
export VECTOR_STORE_ADDRESS=localhost:6334
export VECTOR_STORE_DEFAULT_COLLECTION=my_documents

# Workflow settings
export MAX_ITERATIONS=5
export TOP_K_RETRIEVAL=20
```

### Validate Configuration

```bash
# Validate config syntax and connectivity
./bin/deep-thinking-agent config validate config.json

# Show current active configuration
./bin/deep-thinking-agent config show
```

## Verification

Verify your entire setup is working correctly:

### Step 1: Check All Prerequisites

```bash
# Create verification script
cat > verify_setup.sh << 'EOF'
#!/bin/bash

echo "=== Deep Thinking Agent Setup Verification ==="
echo ""

# Check Go
echo -n "Go version: "
go version || echo "❌ Go not found"

# Check Docker
echo -n "Docker: "
docker --version || echo "❌ Docker not found"

# Check Qdrant
echo -n "Qdrant: "
curl -s http://localhost:6333/collections > /dev/null && echo "✓ Running" || echo "❌ Not accessible"

# Check OpenAI API key
echo -n "OpenAI API key: "
[ -n "$OPENAI_API_KEY" ] && echo "✓ Set" || echo "❌ Not set"

# Check CLI binary
echo -n "CLI binary: "
[ -f ./bin/deep-thinking-agent ] && echo "✓ Built" || echo "❌ Not built"

# Check config
echo -n "Configuration: "
[ -f config.json ] && echo "✓ Exists" || echo "❌ Not found"

echo ""
echo "Setup verification complete"
EOF

chmod +x verify_setup.sh
./verify_setup.sh
```

### Step 2: Run Quick Test

```bash
# Initialize config
./bin/deep-thinking-agent config init

# Validate config
./bin/deep-thinking-agent config validate config.json

# If all checks pass, you're ready!
```

### Step 3: Test End-to-End

```bash
# Create a test document
echo "Machine learning is a subset of artificial intelligence." > test_doc.txt

# Ingest it
./bin/deep-thinking-agent ingest test_doc.txt

# Query it
./bin/deep-thinking-agent query "What is machine learning?"

# Clean up
curl -X DELETE http://localhost:6333/collections/documents
rm test_doc.txt
```

If you see a reasonable answer about machine learning, your setup is complete!

## Troubleshooting

### Common Issues and Solutions

#### "go: command not found"

**Solution:** Go is not installed or not in PATH
```bash
# Add Go to PATH
export PATH=$PATH:/usr/local/go/bin
# Or reinstall Go (see Go Installation section)
```

#### "docker: command not found"

**Solution:** Docker is not installed
```bash
# Install Docker Desktop from docker.com
# Or on Linux: sudo apt-get install docker.io
```

#### "Cannot connect to the Docker daemon"

**Solution:** Docker daemon is not running
```bash
# Start Docker Desktop (macOS/Windows)
# Or on Linux: sudo systemctl start docker
```

#### "curl: (7) Failed to connect to localhost port 6333"

**Solution:** Qdrant is not running
```bash
# Check if Qdrant container exists
docker ps -a | grep qdrant

# Start it if stopped
docker start qdrant

# Or create new container
docker run -d --name qdrant -p 6333:6333 -p 6334:6334 qdrant/qdrant
```

#### "failed to initialize vector store: connection refused"

**Solution:** Wrong Qdrant address in config
```bash
# Check config.json has correct address
# For local Docker: "localhost:6334"
# Note: gRPC port is 6334, HTTP is 6333
```

#### "401 Unauthorized" from OpenAI API

**Solution:** Invalid or missing API key
```bash
# Verify API key is set
echo $OPENAI_API_KEY

# Test API key directly
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"

# Get new API key from platform.openai.com if invalid
```

#### "429 Rate limit exceeded"

**Solution:** Too many API requests
- Wait a few minutes and try again
- Check your OpenAI usage limits
- Consider upgrading your OpenAI plan

#### "insufficient_quota" error

**Solution:** OpenAI account has no credits
- Add payment method at platform.openai.com/account/billing
- Purchase credits or set up automatic recharge
- Check your spending limits

#### Build errors: "package X is not in GOROOT"

**Solution:** Dependencies not downloaded
```bash
go mod download
go mod tidy
```

#### "collection already exists" error

**Solution:** Collection from previous run still exists
```bash
# Delete the collection
curl -X DELETE http://localhost:6333/collections/documents

# Or use a different collection name
./bin/deep-thinking-agent ingest -collection test_docs document.txt
```

### Getting Help

If you encounter issues not covered here:

1. Check [examples/README.md](examples/README.md) for more troubleshooting
2. Search existing [GitHub Issues](https://github.com/yourusername/deep-thinking-agent/issues)
3. Open a new issue with:
   - Your OS and versions (Go, Docker, etc.)
   - Complete error message
   - Steps to reproduce
   - Output of `./verify_setup.sh`

### Debug Mode

For detailed debugging information:

```bash
# Run with verbose output
./bin/deep-thinking-agent query -verbose "your question"

# Check Qdrant logs
docker logs qdrant

# Enable Go debugging
export GODEBUG=http2debug=2
```

## Next Steps

Now that your setup is complete:

1. **Run the examples**: `./examples/01_setup.sh`
2. **Read the examples guide**: [examples/README.md](examples/README.md)
3. **Try interactive mode**: `./bin/deep-thinking-agent query -interactive`
4. **Explore advanced features**: Schema analysis, custom collections, hybrid search

When you're done:
- See [CLEANUP.md](CLEANUP.md) for teardown instructions
- Run `./examples/06_cleanup.sh` to clean up automatically
