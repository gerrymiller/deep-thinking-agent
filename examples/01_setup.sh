#!/bin/bash
# Example 1: Basic Setup

set -e

echo "=== Deep Thinking Agent - Setup Example ==="
echo ""

# Check if running from project root
if [ ! -f "go.mod" ]; then
    echo "Error: Please run this script from the project root directory"
    exit 1
fi

# Build the CLI
echo "Step 1: Building CLI..."
go build -o bin/deep-thinking-agent ./cmd/cli
echo "✓ Built: bin/deep-thinking-agent"
echo ""

# Check for OpenAI API key
if [ -z "$OPENAI_API_KEY" ]; then
    echo "Warning: OPENAI_API_KEY environment variable not set"
    echo ""
    echo "You can configure it using a .env file (recommended):"
    echo "  cp .env.example .env"
    echo "  # Edit .env to add your API key"
    echo ""
    echo "Or set it as an environment variable:"
    echo "  export OPENAI_API_KEY='your-key-here'"
    echo ""
fi

# Initialize config if it doesn't exist
if [ ! -f "config.json" ]; then
    echo "Step 2: Creating default configuration..."
    cp examples/config.example.json config.json
    echo "✓ Created: config.json"
    echo ""
else
    echo "Step 2: Configuration already exists (config.json)"
    echo ""
fi

# Validate configuration
echo "Step 3: Validating configuration..."
if ./bin/deep-thinking-agent config validate config.json; then
    echo "✓ Configuration is valid"
else
    echo "✗ Configuration validation failed"
    echo "  Please check your config.json file"
fi
echo ""

# Check if Qdrant is running
echo "Step 4: Checking Qdrant vector database..."
if curl -s http://localhost:6333/collections >/dev/null 2>&1; then
    echo "✓ Qdrant is running on localhost:6333"
else
    echo "⚠ Qdrant is not running"
    echo "  Start it with: docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant"
fi
echo ""

# Show version
echo "Step 5: Verifying installation..."
./bin/deep-thinking-agent version
echo ""

echo "=== Setup Complete ==="
echo ""
echo "Next steps:"
echo "  1. Ensure OPENAI_API_KEY is set"
echo "  2. Start Qdrant if not running"
echo "  3. Run: ./examples/02_ingest.sh to load sample documents"
echo "  4. Run: ./examples/03_query.sh to test queries"
