#!/bin/bash

# Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
#
# Licensed under the MIT License.
# See LICENSE file in the project root for full license information.

# Example 6: Cleanup and Teardown
#
# This script cleans up all resources created by the example scripts:
# - Stops and removes Qdrant Docker container
# - Deletes all test collections
# - Optionally removes generated files

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Parse arguments
FORCE=false
REMOVE_ALL=false

while [[ $# -gt 0 ]]; do
  case $1 in
    --force)
      FORCE=true
      shift
      ;;
    --all)
      REMOVE_ALL=true
      shift
      ;;
    --help)
      echo "Usage: $0 [--force] [--all]"
      echo ""
      echo "Options:"
      echo "  --force    Skip confirmation prompts"
      echo "  --all      Also remove config.json and binary"
      echo "  --help     Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

echo "=== Deep Thinking Agent - Cleanup ==="
echo ""

# Confirmation prompt unless --force
if [ "$FORCE" = false ]; then
  echo -e "${YELLOW}This will clean up:${NC}"
  echo "  - Qdrant Docker container"
  echo "  - All test collections"
  if [ "$REMOVE_ALL" = true ]; then
    echo "  - config.json"
    echo "  - CLI binary"
  fi
  echo ""
  read -p "Continue? (y/N) " -n 1 -r
  echo ""
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cleanup cancelled."
    exit 0
  fi
fi

# Track what was cleaned
CLEANED=()
ERRORS=()

# Step 1: Stop and remove Qdrant container
echo "Step 1: Stopping Qdrant Docker container..."
if docker ps -q -f name=qdrant > /dev/null 2>&1; then
  if docker rm -f qdrant > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Qdrant container removed"
    CLEANED+=("Qdrant container")
  else
    echo -e "${RED}✗${NC} Failed to remove Qdrant container"
    ERRORS+=("Qdrant container removal failed")
  fi
elif docker ps -a -q -f name=qdrant > /dev/null 2>&1; then
  # Container exists but is stopped
  if docker rm qdrant > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Qdrant container (stopped) removed"
    CLEANED+=("Qdrant container")
  else
    echo -e "${RED}✗${NC} Failed to remove Qdrant container"
    ERRORS+=("Qdrant container removal failed")
  fi
else
  echo -e "${YELLOW}⊘${NC} No Qdrant container found"
fi

# Step 2: Delete collections (if Qdrant was running or is still accessible)
echo ""
echo "Step 2: Deleting test collections..."

# Collections to delete
COLLECTIONS=("documents" "research" "business" "technical")

# Check if Qdrant is still accessible (maybe running elsewhere)
if curl -s http://localhost:6333/collections > /dev/null 2>&1; then
  for collection in "${COLLECTIONS[@]}"; do
    # Check if collection exists first
    if curl -s "http://localhost:6333/collections/$collection" 2>/dev/null | grep -q '"status":"ok"'; then
      if curl -s -X DELETE "http://localhost:6333/collections/$collection" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} Deleted collection: $collection"
        CLEANED+=("Collection: $collection")
      else
        echo -e "${RED}✗${NC} Failed to delete collection: $collection"
        ERRORS+=("Collection deletion failed: $collection")
      fi
    else
      echo -e "${YELLOW}⊘${NC} Collection not found: $collection"
    fi
  done
else
  echo -e "${YELLOW}⊘${NC} Qdrant not accessible (collections already removed with container)"
fi

# Step 3: Remove Docker volumes (optional)
echo ""
echo "Step 3: Checking for Docker volumes..."
if docker volume ls | grep -q qdrant; then
  if [ "$FORCE" = false ]; then
    read -p "Remove Qdrant Docker volumes? (y/N) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
      docker volume prune -f > /dev/null 2>&1
      echo -e "${GREEN}✓${NC} Docker volumes removed"
      CLEANED+=("Docker volumes")
    else
      echo -e "${YELLOW}⊘${NC} Docker volumes kept"
    fi
  else
    docker volume prune -f > /dev/null 2>&1
    echo -e "${GREEN}✓${NC} Docker volumes removed"
    CLEANED+=("Docker volumes")
  fi
else
  echo -e "${YELLOW}⊘${NC} No Qdrant volumes found"
fi

# Step 4: Remove local storage directory if it exists
if [ -d "../qdrant_storage" ]; then
  echo ""
  if [ "$FORCE" = false ]; then
    read -p "Remove local Qdrant storage directory? (y/N) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
      rm -rf ../qdrant_storage
      echo -e "${GREEN}✓${NC} Local storage directory removed"
      CLEANED+=("Local storage directory")
    else
      echo -e "${YELLOW}⊘${NC} Local storage directory kept"
    fi
  else
    rm -rf ../qdrant_storage
    echo -e "${GREEN}✓${NC} Local storage directory removed"
    CLEANED+=("Local storage directory")
  fi
fi

# Step 5: Remove generated files (if --all flag used)
if [ "$REMOVE_ALL" = true ]; then
  echo ""
  echo "Step 4: Removing generated files..."

  # Remove config.json
  if [ -f "../config.json" ]; then
    rm ../config.json
    echo -e "${GREEN}✓${NC} config.json removed"
    CLEANED+=("config.json")
  else
    echo -e "${YELLOW}⊘${NC} config.json not found"
  fi

  # Remove binary
  if [ -f "../bin/deep-thinking-agent" ]; then
    rm ../bin/deep-thinking-agent
    echo -e "${GREEN}✓${NC} CLI binary removed"
    CLEANED+=("CLI binary")

    # Remove bin directory if empty
    if [ -z "$(ls -A ../bin)" ]; then
      rmdir ../bin
      echo -e "${GREEN}✓${NC} bin directory removed"
    fi
  else
    echo -e "${YELLOW}⊘${NC} CLI binary not found"
  fi
fi

# Summary
echo ""
echo "=== Cleanup Summary ==="
echo ""

if [ ${#CLEANED[@]} -gt 0 ]; then
  echo -e "${GREEN}Cleaned up:${NC}"
  for item in "${CLEANED[@]}"; do
    echo "  ✓ $item"
  done
fi

if [ ${#ERRORS[@]} -gt 0 ]; then
  echo ""
  echo -e "${RED}Errors:${NC}"
  for error in "${ERRORS[@]}"; do
    echo "  ✗ $error"
  done
fi

if [ ${#CLEANED[@]} -eq 0 ] && [ ${#ERRORS[@]} -eq 0 ]; then
  echo -e "${YELLOW}Nothing to clean up${NC}"
fi

echo ""
echo -e "${GREEN}Cleanup complete!${NC}"
echo ""

# Helpful next steps
if [ "$REMOVE_ALL" = false ]; then
  echo "Note: config.json and CLI binary were preserved."
  echo "To remove them too, run: $0 --all"
  echo ""
fi

echo "To start fresh:"
echo "  1. Start Qdrant: docker run -d --name qdrant -p 6333:6333 -p 6334:6334 qdrant/qdrant"
echo "  2. Build CLI: go build -o bin/deep-thinking-agent ./cmd/cli"
echo "  3. Configure: ./bin/deep-thinking-agent config init"
echo "  4. Run examples: ./01_setup.sh"
echo ""

# Check OpenAI usage
echo "Remember to check your OpenAI usage:"
echo "  https://platform.openai.com/usage"
echo ""
