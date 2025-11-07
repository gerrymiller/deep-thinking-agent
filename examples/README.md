# Deep Thinking Agent Examples

This directory contains examples demonstrating how to use the Deep Thinking Agent CLI.

## ⚠️ Cost Warning

These examples use OpenAI APIs and **will incur costs**. See [Cost Estimates](#cost-estimates) below for details. We recommend:
- Setting [spending limits](https://platform.openai.com/account/limits) in your OpenAI account
- Starting with simple examples (01-03) before running expensive ones (04-05)
- Using `--no-schema` flag during testing to reduce costs

## Automated Scripts

The easiest way to run examples:

| Script | Purpose | Est. Cost |
|--------|---------|-----------|
| `01_setup.sh` | Build and validate setup | Free |
| `02_ingest.sh` | Ingest 3 sample documents | $0.15-0.30 |
| `03_query.sh` | Run 4 simple queries | $0.24-0.60 |
| `04_advanced.sh` | Run 5 complex multi-hop queries | $0.75-1.55 |
| `05_ingestion_patterns.sh` | Test various ingestion patterns | $0.45-0.90 |
| `06_cleanup.sh` | Clean up all resources | Free |

**Total cost for all examples: $1.59-3.35**

Run them in order from the project root:
```bash
./examples/01_setup.sh
./examples/02_ingest.sh
./examples/03_query.sh
./examples/04_advanced.sh
./examples/05_ingestion_patterns.sh
# When done:
./examples/06_cleanup.sh
```

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

3. **Configure API key**:
   ```bash
   # Option 1: Using .env file (recommended)
   cp .env.example .env
   # Edit .env to add your OPENAI_API_KEY
   
   # Option 2: Using environment variable
   export OPENAI_API_KEY="your-api-key-here"
   ```

4. **Create configuration** (optional):
   ```bash
   cp examples/config.example.json config.json
   # Edit config.json if needed (CLI works without it using defaults)
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

## Cost Estimates

### Per Operation Costs (Approximate)

Based on OpenAI pricing as of 2025:

**Document Ingestion:**
- With schema analysis (default): $0.012-0.052 per document
- Without schema (`--no-schema`): $0.002-0.005 per document
- Depends on: document length, complexity

**Queries:**
- Simple query (1-3 iterations): $0.06-0.15
- Complex query (5-10 iterations): $0.15-0.31
- Interactive query session (10 queries): $0.60-2.00
- Depends on: query complexity, iterations needed, context length

**Example Scripts Breakdown:**
| Script | Operations | Estimated Cost |
|--------|-----------|----------------|
| 02_ingest.sh | 3 docs with schema | $0.15-0.30 |
| 03_query.sh | 4 simple queries | $0.24-0.60 |
| 04_advanced.sh | 5 complex queries (15-20 iter) | $0.75-1.55 |
| 05_ingestion_patterns.sh | Multiple ingestion modes | $0.45-0.90 |
| **Total** | | **$1.59-3.35** |

### Cost Reduction Tips

1. **Use `--no-schema` flag**
   ```bash
   ./bin/deep-thinking-agent ingest --no-schema document.txt
   ```
   Saves ~$0.01-0.05 per document by skipping LLM schema analysis

2. **Reduce max iterations**
   ```bash
   ./bin/deep-thinking-agent query -max-iterations 3 "simple question"
   ```
   Each iteration costs ~$0.03-0.06

3. **Use cheaper models** (edit config.json)
   ```json
   {
     "reasoning_llm": {"model": "gpt-4o-mini"},
     "fast_llm": {"model": "gpt-4o-mini"}
   }
   ```
   Reduces costs by ~60-80% but may affect quality

4. **Start with smaller documents**
   Test with short documents before ingesting large corpora

5. **Monitor usage**
   Check [platform.openai.com/usage](https://platform.openai.com/usage) regularly

### Setting Spending Limits

Protect yourself from unexpected charges:

1. Go to [OpenAI Billing Settings](https://platform.openai.com/account/billing/limits)
2. Set a monthly budget limit (e.g., $10 for testing)
3. Enable email notifications for:
   - 75% of limit reached
   - 90% of limit reached
   - Limit reached

## After Running Examples

### Cleanup

When you're done with the examples, clean up resources:

```bash
# Automated cleanup (recommended)
./06_cleanup.sh

# Or with no prompts
./06_cleanup.sh --force

# Remove everything including config and binary
./06_cleanup.sh --all
```

This will:
- Stop and remove Qdrant Docker container
- Delete all test collections
- Optionally remove generated files

For detailed cleanup instructions, see [CLEANUP.md](../CLEANUP.md).

### Manual Cleanup

If you prefer manual cleanup:

```bash
# Stop Qdrant
docker rm -f qdrant

# Delete collections
curl -X DELETE http://localhost:6333/collections/documents
curl -X DELETE http://localhost:6333/collections/research
curl -X DELETE http://localhost:6333/collections/business
curl -X DELETE http://localhost:6333/collections/technical

# Remove files (optional)
rm config.json
rm bin/deep-thinking-agent
```

### Check Your Costs

After running examples:

1. Visit [OpenAI Usage Dashboard](https://platform.openai.com/usage)
2. Review costs by date and model
3. Verify charges match expectations (~$2-5 for all examples)

## Next Steps

- Review the main [README](../README.md) for architecture details
- Check [SETUP.md](../SETUP.md) for detailed installation guide
- See [CLEANUP.md](../CLEANUP.md) for complete cleanup instructions
- Read [AGENTS.md](../AGENTS.md) for development guidelines
- Explore the source code in `pkg/` and `cmd/` directories
- Experiment with different query types and document collections
