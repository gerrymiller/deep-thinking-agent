#!/bin/bash
# Example 5: Advanced Document Ingestion Patterns

set -e

echo "=== Deep Thinking Agent - Ingestion Patterns Example ==="
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

echo "This example demonstrates different document ingestion patterns:"
echo "  1. Single file ingestion"
echo "  2. Batch ingestion from directory"
echo "  3. Recursive directory ingestion"
echo "  4. Custom collection organization"
echo ""
read -p "Press Enter to start..."
echo ""

# Pattern 1: Single file ingestion
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Pattern 1: Single File Ingestion"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Ingesting a single document with detailed output..."
./bin/deep-thinking-agent ingest -verbose examples/documents/sample_research.md
echo ""
echo "✓ Single file ingested to default 'documents' collection"
echo ""
read -p "Press Enter to continue..."
echo ""

# Pattern 2: Multiple specific files
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Pattern 2: Multiple Specific Files"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Ingesting multiple files at once..."
./bin/deep-thinking-agent ingest \
  examples/documents/sample_research.md \
  examples/documents/company_report.md
echo ""
echo "✓ Multiple files ingested in a single command"
echo ""
read -p "Press Enter to continue..."
echo ""

# Pattern 3: Recursive directory ingestion
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Pattern 3: Recursive Directory Ingestion"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Ingesting all documents in a directory tree..."
./bin/deep-thinking-agent ingest -recursive -verbose examples/documents
echo ""
echo "✓ All documents in directory tree ingested"
echo ""
read -p "Press Enter to continue..."
echo ""

# Pattern 4: Custom collection organization
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Pattern 4: Custom Collection Organization"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Organizing documents into separate collections..."
echo ""

echo "Creating 'research' collection for academic papers..."
./bin/deep-thinking-agent ingest \
  -collection research \
  -verbose \
  examples/documents/sample_research.md
echo ""

echo "Creating 'business' collection for company documents..."
./bin/deep-thinking-agent ingest \
  -collection business \
  -verbose \
  examples/documents/company_report.md
echo ""

echo "Creating 'technical' collection for technical guides..."
./bin/deep-thinking-agent ingest \
  -collection technical \
  -verbose \
  examples/documents/technical_guide.txt
echo ""

echo "✓ Documents organized into domain-specific collections"
echo ""

echo "=== Ingestion Patterns Complete ==="
echo ""
echo "Key takeaways:"
echo "  • Use single file ingestion for targeted updates"
echo "  • Use recursive ingestion for bulk loading"
echo "  • Use custom collections to organize by domain/type"
echo "  • Use -verbose flag to monitor progress"
echo ""
echo "Collection structure:"
echo "  • research   - Academic and research papers"
echo "  • business   - Business reports and financial docs"
echo "  • technical  - Technical documentation and guides"
echo "  • documents  - Default collection (mixed content)"
echo ""
echo "Query specific collections with the -collection flag:"
echo "  ./bin/deep-thinking-agent query -collection research \"your question\""
