#!/bin/bash
# Example 3: Simple Queries

set -e

echo "=== Deep Thinking Agent - Query Examples ==="
echo ""

# Check if CLI exists
if [ ! -f "bin/deep-thinking-agent" ]; then
    echo "Error: CLI not found. Run ./examples/01_setup.sh first"
    exit 1
fi

echo "Running example queries against ingested documents..."
echo ""
echo "Tip: Add -verbose flag to see the reasoning process"
echo ""

# Query 1: Simple factual query
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Query 1: Healthcare ML Applications"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./bin/deep-thinking-agent query "What are the main applications of machine learning in healthcare?"
echo ""
read -p "Press Enter to continue to next query..."
echo ""

# Query 2: Specific data retrieval
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Query 2: Financial Metrics"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./bin/deep-thinking-agent query "What was TechCorp's revenue performance in Q4 2024?"
echo ""
read -p "Press Enter to continue to next query..."
echo ""

# Query 3: Technical information
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Query 3: Kubernetes Best Practices"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./bin/deep-thinking-agent query "What are the best practices for Kubernetes deployments?"
echo ""
read -p "Press Enter to continue to next query..."
echo ""

# Query 4: Ethical considerations
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Query 4: Ethical Considerations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./bin/deep-thinking-agent query "What ethical considerations are mentioned regarding AI in healthcare?"
echo ""

echo "=== Query Examples Complete ==="
echo ""
echo "Try more queries:"
echo "  - Interactive mode: ./bin/deep-thinking-agent query -interactive"
echo "  - Verbose mode: ./bin/deep-thinking-agent query -verbose \"your question\""
echo "  - Advanced queries: ./examples/04_advanced.sh"
