#!/bin/bash
# Example 2: Document Ingestion

set -e

echo "=== Deep Thinking Agent - Document Ingestion Example ==="
echo ""

# Check if CLI exists
if [ ! -f "bin/deep-thinking-agent" ]; then
    echo "Error: CLI not found. Run ./examples/01_setup.sh first"
    exit 1
fi

# Check if Qdrant is running
if ! curl -s http://localhost:6333/collections >/dev/null 2>&1; then
    echo "Error: Qdrant is not running"
    echo "Start it with: docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant"
    exit 1
fi

echo "Ingesting sample documents..."
echo ""

# Ingest healthcare research document
echo "1. Ingesting sample_research.md (Healthcare ML research)..."
./bin/deep-thinking-agent ingest -verbose examples/documents/sample_research.md
echo ""

# Ingest company report
echo "2. Ingesting company_report.md (Financial report)..."
./bin/deep-thinking-agent ingest -verbose examples/documents/company_report.md
echo ""

# Ingest technical guide
echo "3. Ingesting technical_guide.txt (Kubernetes guide)..."
./bin/deep-thinking-agent ingest -verbose examples/documents/technical_guide.txt
echo ""

echo "=== Ingestion Complete ==="
echo ""
echo "Documents have been processed and stored in the vector database."
echo "You can now run queries with: ./examples/03_query.sh"
echo ""
echo "Or try manual queries:"
echo "  ./bin/deep-thinking-agent query \"What are the main topics in the documents?\""
