#!/bin/bash
# Example 4: Advanced Multi-Hop Queries

set -e

echo "=== Deep Thinking Agent - Advanced Query Examples ==="
echo ""
echo "These examples demonstrate the deep thinking system's ability to:"
echo "  - Reason across multiple documents"
echo "  - Synthesize information from different sources"
echo "  - Perform multi-step analysis"
echo ""

# Check if CLI exists
if [ ! -f "bin/deep-thinking-agent" ]; then
    echo "Error: CLI not found. Run ./examples/01_setup.sh first"
    exit 1
fi

echo "Note: Add -verbose flag to see the reasoning steps"
echo ""
read -p "Press Enter to start advanced queries..."
echo ""

# Query 1: Cross-document synthesis
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Advanced Query 1: Cross-Document Synthesis"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Question: What are the common themes across all documents regarding"
echo "          technology adoption and challenges?"
echo ""
./bin/deep-thinking-agent query -max-iterations 15 \
  "What are the common themes across all documents regarding technology adoption and challenges?"
echo ""
read -p "Press Enter to continue..."
echo ""

# Query 2: Multi-step reasoning
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Advanced Query 2: Multi-Step Reasoning"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Question: Based on the healthcare research, what are the main"
echo "          challenges in AI deployment, and how might TechCorp's"
echo "          AI initiatives address these challenges?"
echo ""
./bin/deep-thinking-agent query -max-iterations 15 \
  "Based on the healthcare research, what are the main challenges in AI deployment, and how might TechCorp's AI initiatives address these challenges?"
echo ""
read -p "Press Enter to continue..."
echo ""

# Query 3: Comparative analysis
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Advanced Query 3: Comparative Analysis"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Question: Compare the deployment best practices mentioned for"
echo "          Kubernetes with TechCorp's cloud infrastructure approach"
echo ""
./bin/deep-thinking-agent query -max-iterations 15 \
  "Compare the deployment best practices mentioned for Kubernetes with TechCorp's cloud infrastructure approach"
echo ""
read -p "Press Enter to continue..."
echo ""

# Query 4: Risk analysis
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Advanced Query 4: Risk Analysis Across Domains"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Question: What risk factors are mentioned across all documents,"
echo "          and how do they relate to each other?"
echo ""
./bin/deep-thinking-agent query -max-iterations 15 \
  "What risk factors are mentioned across all documents, and how do they relate to each other?"
echo ""
read -p "Press Enter to continue..."
echo ""

# Query 5: Strategic synthesis
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Advanced Query 5: Strategic Synthesis"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Question: If TechCorp wanted to enter the healthcare AI market,"
echo "          what would be their key advantages and challenges based"
echo "          on the available information?"
echo ""
./bin/deep-thinking-agent query -max-iterations 20 \
  "If TechCorp wanted to enter the healthcare AI market, what would be their key advantages and challenges based on the available information?"
echo ""

echo "=== Advanced Query Examples Complete ==="
echo ""
echo "These queries demonstrate:"
echo "  ✓ Multi-document reasoning"
echo "  ✓ Cross-domain synthesis"
echo "  ✓ Strategic analysis"
echo "  ✓ Complex comparative reasoning"
echo ""
echo "Experiment with your own complex queries in interactive mode:"
echo "  ./bin/deep-thinking-agent query -interactive -max-iterations 20"
