#!/bin/bash
# Example 2: Document Ingestion with Validation

set -e

echo "=== Deep Thinking Agent - Document Ingestion Example ==="
echo ""

# Check if CLI exists
if [ ! -f "bin/deep-thinking-agent" ]; then
    echo "Error: CLI not found. Run ./examples/01_setup.sh first"
    exit 1
fi

# Check if Qdrant is running
echo "Checking Qdrant connection..."
if ! curl -s http://localhost:6333/collections >/dev/null 2>&1; then
    echo "Error: Qdrant is not running"
    echo "Start it with: docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant"
    exit 1
fi
echo "✅ Qdrant is running"
echo ""

# Check collection status before ingestion
echo "Checking collection status..."
if curl -s http://localhost:6333/collections/documents 2>&1 | grep -q "points_count"; then
    POINTS_BEFORE=$(curl -s http://localhost:6333/collections/documents | grep -o '"points_count":[0-9]*' | cut -d: -f2)
    echo "⚠️  Collection 'documents' exists with $POINTS_BEFORE points"
else
    echo "Collection 'documents' does not exist - will be created"
fi
echo ""

echo "Ingesting sample documents..."
echo ""

# Ingest documents
for i in 1 2 3; do
    case $i in
        1) FILE="examples/documents/sample_research.md"; DESC="Healthcare ML research" ;;
        2) FILE="examples/documents/company_report.md"; DESC="Financial report" ;;
        3) FILE="examples/documents/technical_guide.txt"; DESC="Kubernetes guide" ;;
    esac
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "$i. $DESC"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "File: $FILE ($(wc -c < "$FILE" | tr -d ' ') bytes)"
    
    if ./bin/deep-thinking-agent ingest -verbose "$FILE"; then
        echo "✅ Success"
    else
        echo "❌ Failed"
        exit 1
    fi
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Verifying Ingestion"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
sleep 2

if curl -s http://localhost:6333/collections/documents >/dev/null 2>&1; then
    INFO=$(curl -s http://localhost:6333/collections/documents)
    POINTS=$(echo "$INFO" | grep -o '"points_count":[0-9]*' | cut -d: -f2)
    DIMS=$(echo "$INFO" | grep -o '"size":[0-9]*' | head -1 | cut -d: -f2)
    
    echo "✅ Collection 'documents' verified"
    echo "📊 Total chunks: $POINTS"
    echo "📏 Vector size: $DIMS"
    
    if [ "$POINTS" -gt 0 ]; then
        echo ""
        echo "✅ Ingestion successful!"
        echo ""
        echo "Next: ./examples/03_query.sh"
    else
        echo ""
        echo "⚠️  Collection has 0 points!"
        exit 1
    fi
else
    echo "❌ Collection 'documents' not found!"
    echo "Check: docker logs \$(docker ps -q --filter ancestor=qdrant/qdrant)"
    exit 1
fi
