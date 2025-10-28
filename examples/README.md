# Deep Thinking Agent Examples

This directory contains examples demonstrating how to use the Deep Thinking Agent CLI.

## Prerequisites

1. **Build the CLI**:
   ```bash
   cd /path/to/deep-thinking-agent
   go build -o bin/deep-thinking-agent ./cmd/cli
   ```

2. **Set up Qdrant** (vector database):
   ```bash
   docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
   ```

3. **Set environment variables**:
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```

4. **Create configuration**:
   ```bash
   cp examples/config.example.json config.json
   # Edit config.json if needed
   ```

## Example 1: Basic Setup

Initialize a default configuration file:

```bash
./bin/deep-thinking-agent config init
```

Validate your configuration:

```bash
./bin/deep-thinking-agent config validate config.json
```

Show current configuration:

```bash
./bin/deep-thinking-agent config show
```

## Example 2: Document Ingestion

Ingest a single document:

```bash
./bin/deep-thinking-agent ingest examples/documents/sample_research.md
```

Ingest all documents in a directory:

```bash
./bin/deep-thinking-agent ingest -recursive examples/documents
```

Ingest with verbose output:

```bash
./bin/deep-thinking-agent ingest -verbose examples/documents/company_report.md
```

Ingest to a custom collection:

```bash
./bin/deep-thinking-agent ingest -collection research_papers examples/documents/sample_research.md
```

## Example 3: Simple Queries

Execute a single query:

```bash
./bin/deep-thinking-agent query "What are the main applications of machine learning in healthcare?"
```

Query with verbose output (shows execution steps):

```bash
./bin/deep-thinking-agent query -verbose "What ethical considerations are mentioned regarding AI in healthcare?"
```

Query with custom max iterations:

```bash
./bin/deep-thinking-agent query -max-iterations 5 "Summarize the key findings about diagnostic imaging"
```

## Example 4: Interactive Mode

Start interactive query mode:

```bash
./bin/deep-thinking-agent query -interactive
```

Then type queries at the prompt:

```
Query> What were TechCorp's revenue numbers?
Query> How does their cloud division perform?
Query> What are their main risk factors?
Query> exit
```

## Example 5: Complete Workflow

Here's a complete workflow from setup to querying:

```bash
# 1. Build the CLI
go build -o bin/deep-thinking-agent ./cmd/cli

# 2. Initialize configuration
./bin/deep-thinking-agent config init

# 3. Edit config.json to add your OpenAI API key
# (or set OPENAI_API_KEY environment variable)

# 4. Start Qdrant
docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant

# 5. Ingest documents
./bin/deep-thinking-agent ingest -recursive -verbose examples/documents

# 6. Run queries
./bin/deep-thinking-agent query "What are the main topics covered in the ingested documents?"

./bin/deep-thinking-agent query "Compare the revenue growth mentioned in the company report with industry standards"

./bin/deep-thinking-agent query "What are best practices for Kubernetes deployments?"
```

## Example 6: Advanced Multi-Hop Queries

The deep thinking system excels at complex, multi-hop queries that require reasoning across multiple pieces of information:

```bash
# Multi-step reasoning query
./bin/deep-thinking-agent query -verbose "Based on the healthcare research, what are the main challenges in AI deployment, and how might TechCorp's AI initiatives address these?"

# Comparative analysis query
./bin/deep-thinking-agent query "Compare the approaches to AI deployment in healthcare versus enterprise software based on the documents"

# Synthesis query
./bin/deep-thinking-agent query "What common themes exist across all three documents regarding technology adoption and challenges?"
```

## Example 7: Working with Different Document Types

The system supports multiple document formats:

```bash
# Ingest text files
./bin/deep-thinking-agent ingest examples/documents/technical_guide.txt

# Ingest markdown files
./bin/deep-thinking-agent ingest examples/documents/sample_research.md

# Mix of formats
./bin/deep-thinking-agent ingest -recursive examples/documents
```

## Troubleshooting

### "Failed to connect to vector store"
- Ensure Qdrant is running: `docker ps | grep qdrant`
- Check the address in config.json matches Qdrant's address (default: localhost:6334)

### "Failed to initialize LLM"
- Verify OPENAI_API_KEY is set: `echo $OPENAI_API_KEY`
- Check your OpenAI API key is valid and has credits

### "No documents found" when querying
- Ensure documents were ingested successfully
- Check the collection name matches between ingest and query operations
- Verify Qdrant contains the data: `curl http://localhost:6333/collections`

### "Context deadline exceeded"
- The query may be too complex for the default timeout
- Try a simpler query first
- Check Qdrant and OpenAI API are responding

## Configuration Options

### LLM Configuration

```json
{
  "reasoning_llm": {
    "provider": "openai",
    "model": "gpt-4o",              // Use "gpt-4o" for best results
    "default_temperature": 0.7      // Higher = more creative
  },
  "fast_llm": {
    "provider": "openai",
    "model": "gpt-4o-mini",      // Faster, cheaper for simple tasks
    "default_temperature": 0.5
  }
}
```

### Workflow Configuration

```json
{
  "workflow": {
    "max_iterations": 10,          // Maximum reasoning steps
    "top_k_retrieval": 10,         // Number of documents to retrieve
    "top_n_reranking": 3          // Number of documents after reranking
  }
}
```

## Performance Tips

1. **Use appropriate models**: gpt-4o for complex reasoning, gpt-4o-mini for simple tasks
2. **Adjust max_iterations**: Lower for simple queries (faster), higher for complex reasoning
3. **Tune retrieval parameters**: Increase top_k_retrieval for better recall
4. **Batch ingestion**: Use `-recursive` flag to ingest multiple documents at once
5. **Collection organization**: Use different collections for different document types

## Next Steps

- Review the main [README](../README.md) for architecture details
- Check [CLAUDE.md](../CLAUDE.md) for development guidelines
- Explore the source code in `pkg/` and `cmd/` directories
- Experiment with different query types and document collections
