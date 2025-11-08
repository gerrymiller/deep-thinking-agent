# Cleanup Guide

Complete guide for cleaning up Deep Thinking Agent resources, stopping services, and removing data.

## Table of Contents

1. [Quick Cleanup](#quick-cleanup)
2. [Manual Cleanup Steps](#manual-cleanup-steps)
3. [Cleanup by Component](#cleanup-by-component)
4. [Cost Tracking](#cost-tracking)
5. [Troubleshooting Cleanup](#troubleshooting-cleanup)

## Quick Cleanup

The fastest way to clean up everything:

```bash
# Run the automated cleanup script
cd examples
./06_cleanup.sh

# Or with no prompts (automatic mode)
./06_cleanup.sh --force
```

This script will:
- Stop and remove Qdrant Docker container
- Delete all test collections
- Optionally remove generated files
- Show summary of what was cleaned

## Manual Cleanup Steps

If you prefer to clean up manually or the script doesn't work:

### Step 1: Stop Qdrant

**If using Docker:**
```bash
# Find the Qdrant container
docker ps | grep qdrant

# Stop the container
docker stop qdrant

# Remove the container
docker rm qdrant

# Or stop and remove in one command
docker rm -f qdrant
```

**If using Qdrant Cloud:**
- No action needed - cloud service continues running
- Delete collections if you want to clean data (see Step 2)
- Or delete the entire cluster from cloud.qdrant.io

**If using local binary:**
```bash
# Find the process
ps aux | grep qdrant

# Kill it
pkill qdrant

# Or if you know the PID
kill <pid>
```

### Step 2: Remove Collections

Delete all test collections created by the examples:

```bash
# Default collection
curl -X DELETE http://localhost:6333/collections/documents

# Collections from 05_ingestion_patterns.sh
curl -X DELETE http://localhost:6333/collections/research
curl -X DELETE http://localhost:6333/collections/business
curl -X DELETE http://localhost:6333/collections/technical

# Verify all collections are removed
curl http://localhost:6333/collections
# Should return: {"result":{"collections":[]},"status":"ok"}
```

**For Qdrant Cloud:**
```bash
# Replace with your cluster URL and API key
curl -X DELETE https://your-cluster.qdrant.io:6333/collections/documents \
  -H "api-key: your-api-key"
```

### Step 3: Remove Docker Volumes

If you used persistent storage with Qdrant:

```bash
# List Docker volumes
docker volume ls | grep qdrant

# Remove Qdrant volumes
docker volume rm qdrant_storage

# Or remove all unused volumes (careful!)
docker volume prune -f
```

If you mounted a local directory:

```bash
# Remove the local storage directory
rm -rf ./qdrant_storage
```

### Step 4: Remove Generated Files

**Configuration file:**
```bash
# Remove config.json (if you want to start fresh)
rm config.json

# Or keep it for future use
```

**Built binary:**
```bash
# Remove the CLI binary
rm bin/deep-thinking-agent

# Remove entire bin directory
rm -rf bin/
```

**Temporary files:**
```bash
# Remove any test documents you created
rm test_doc.txt

# Remove verification script if created
rm verify_setup.sh
```

### Step 5: Unset Environment Variables

**Linux/macOS:**
```bash
# Unset for current session
unset OPENAI_API_KEY

# Remove from shell config (if you added it permanently)
# Edit ~/.bashrc or ~/.zshrc and remove the export line
```

**Windows (PowerShell):**
```powershell
# Remove from current session
Remove-Item Env:OPENAI_API_KEY

# Remove permanently
[System.Environment]::SetEnvironmentVariable('OPENAI_API_KEY', $null, 'User')
```

## Cleanup by Component

### Minimal Cleanup (Keep Infrastructure)

If you just want to clear data but keep services running:

```bash
# Only delete collections
curl -X DELETE http://localhost:6333/collections/documents
curl -X DELETE http://localhost:6333/collections/research
curl -X DELETE http://localhost:6333/collections/business
curl -X DELETE http://localhost:6333/collections/technical

# Qdrant continues running for future use
```

### Standard Cleanup (Stop Services, Keep Files)

Stop services but keep configuration and binary:

```bash
# Stop Qdrant
docker stop qdrant

# Delete collections (optional)
curl -X DELETE http://localhost:6333/collections/documents

# Keep config.json and binary for next time
```

### Complete Cleanup (Remove Everything)

Full teardown for completely removing the project:

```bash
# Stop and remove Qdrant
docker rm -f qdrant
docker volume prune -f

# Remove all generated files
rm config.json
rm -rf bin/
rm -rf qdrant_storage/

# Unset environment variables
unset OPENAI_API_KEY

# Optionally remove the repository
cd ..
rm -rf deep-thinking-agent/
```

## Cost Tracking

### Check Your OpenAI Usage

1. **View Usage Dashboard**
   - Go to [platform.openai.com/usage](https://platform.openai.com/usage)
   - View costs by day, model, and operation

2. **Export Usage Data**
   ```bash
   # Using OpenAI API (requires jq)
   curl https://api.openai.com/v1/usage \
     -H "Authorization: Bearer $OPENAI_API_KEY" \
     | jq
   ```

3. **Set Up Alerts**
   - Go to [Billing Settings](https://platform.openai.com/account/billing/limits)
   - Configure email alerts for usage thresholds
   - Set hard limits to prevent unexpected charges

### Estimated Costs for Examples

If you ran all the example scripts:

| Script | Estimated Cost |
|--------|---------------|
| 02_ingest.sh | $0.15-0.30 (3 documents with schema analysis) |
| 03_query.sh | $0.24-0.60 (4 simple queries) |
| 04_advanced.sh | $0.75-1.55 (5 complex multi-hop queries) |
| 05_ingestion_patterns.sh | $0.45-0.90 (additional documents) |
| **Total** | **$1.59-3.35** |

Actual costs depend on:
- Document length
- Query complexity
- Number of reasoning iterations
- Model pricing (subject to change)

### Cost Reduction Tips

For future runs to minimize costs:

1. **Use `--no-schema` flag**
   ```bash
   ./bin/deep-thinking-agent ingest --no-schema document.txt
   # Skips LLM-based schema analysis, uses simple chunking
   ```

2. **Reduce max iterations**
   ```bash
   # Edit config.json
   "workflow": {
     "max_iterations": 3  # Instead of 10
   }
   ```

3. **Use cheaper models**
   ```bash
   # Edit config.json - use gpt-5-mini for everything
   "reasoning_llm": {
     "model": "gpt-5-mini"  # Instead of gpt-5
   }
   ```

4. **Test with smaller documents**
   ```bash
   # Create minimal test documents
   echo "Short test content" > test.txt
   ./bin/deep-thinking-agent ingest test.txt
   ```

## Troubleshooting Cleanup

### "Container not found" when trying to remove

**Solution:** Container might have a different name
```bash
# List all containers
docker ps -a

# Remove by ID
docker rm -f <container_id>
```

### "Cannot remove volume: volume is in use"

**Solution:** Container is still running
```bash
# Stop and remove container first
docker rm -f qdrant

# Then remove volume
docker volume rm qdrant_storage
```

### "Connection refused" when deleting collections

**Solution:** Qdrant is not running or wrong address
```bash
# Check if Qdrant is accessible
curl http://localhost:6333/collections

# If not accessible, collections are already gone
# (they only exist while Qdrant is running unless using persistent storage)
```

### Collections still exist after deletion

**Solution:** Using persistent storage, need to delete volume
```bash
# Delete the volume
docker volume rm qdrant_storage

# Or delete the local directory
rm -rf ./qdrant_storage
```

### "Permission denied" when removing files

**Solution:** Files owned by Docker or root
```bash
# Use sudo
sudo rm -rf qdrant_storage/

# Or change ownership
sudo chown -R $USER:$USER qdrant_storage/
rm -rf qdrant_storage/
```

## When to Clean Up

### Always Clean Up
- **After testing** - Don't leave services running unnecessarily
- **Before uninstalling** - Remove all traces of the project
- **When switching environments** - Clean local before moving to cloud
- **To reset state** - Start fresh with clean collections

### Keep Running
- **During active development** - Restart services takes time
- **For production use** - Keep Qdrant running continuously
- **When testing multiple queries** - Avoid re-ingesting documents

### Partial Cleanup
- **Between test runs** - Delete collections but keep Qdrant running
- **When changing document set** - Clear old collections, keep infrastructure
- **After examples** - Remove test data but keep configuration

## Cleanup Checklist

Before considering cleanup complete, verify:

- [ ] Qdrant container stopped and removed
- [ ] All Docker volumes removed (if desired)
- [ ] Collections deleted (verify with curl)
- [ ] config.json removed or backed up
- [ ] Binary removed from bin/
- [ ] Environment variables unset
- [ ] No zombie processes running (`ps aux | grep qdrant`)
- [ ] OpenAI usage checked and within expectations
- [ ] Local storage directory removed (if used)

## Restore After Cleanup

To start fresh after complete cleanup:

```bash
# 1. Restart Qdrant
docker run -d --name qdrant -p 6333:6333 -p 6334:6334 qdrant/qdrant

# 2. Rebuild CLI
go build -o bin/deep-thinking-agent ./cmd/cli

# 3. Reinitialize config
export OPENAI_API_KEY="sk-your-key"
./bin/deep-thinking-agent config init

# 4. Verify setup
./examples/01_setup.sh

# Ready to use again!
```

## Support

If you encounter issues during cleanup:

1. Check this guide's troubleshooting section
2. Review [SETUP.md](SETUP.md) for original configuration
3. See [examples/README.md](examples/README.md) for additional help
4. Open an issue at [GitHub Issues](https://github.com/yourusername/deep-thinking-agent/issues)

## Summary

**Quick cleanup:**
```bash
./examples/06_cleanup.sh
```

**Manual cleanup:**
```bash
docker rm -f qdrant
curl -X DELETE http://localhost:6333/collections/documents
rm config.json
rm -rf bin/
```

**Cost tracking:**
- Check [platform.openai.com/usage](https://platform.openai.com/usage)
- Review costs after running examples (~$2-5 total)
- Set spending limits to avoid surprises

Now you know exactly how to clean up after using Deep Thinking Agent!
