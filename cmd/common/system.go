// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package common

import (
	"context"
	"fmt"
	"strings"

	"deep-thinking-agent/pkg/agent"
	"deep-thinking-agent/pkg/document/chunker"
	"deep-thinking-agent/pkg/embedding"
	"deep-thinking-agent/pkg/llm"
	"deep-thinking-agent/pkg/llm/openai"
	"deep-thinking-agent/pkg/nodes"
	"deep-thinking-agent/pkg/schema"
	"deep-thinking-agent/pkg/vectorstore"
	"deep-thinking-agent/pkg/vectorstore/qdrant"
	"deep-thinking-agent/pkg/workflow"

	"github.com/google/uuid"
)

// System encapsulates all components of the deep thinking agent.
type System struct {
	Config         *Config
	ReasoningLLM   llm.Provider
	FastLLM        llm.Provider
	Embedder       embedding.Embedder
	VectorStore    vectorstore.Store
	SchemaResolver *schema.Resolver
	Executor       *workflow.Executor
}

// InitializeSystem creates and initializes all system components based on configuration.
func InitializeSystem(config *Config) (*System, error) {
	sys := &System{
		Config: config,
	}

	// Initialize LLM providers
	if err := sys.initLLMs(); err != nil {
		return nil, fmt.Errorf("failed to initialize LLMs: %w", err)
	}

	// Initialize embedder
	if err := sys.initEmbedder(); err != nil {
		return nil, fmt.Errorf("failed to initialize embedder: %w", err)
	}

	// Initialize vector store
	if err := sys.initVectorStore(); err != nil {
		return nil, fmt.Errorf("failed to initialize vector store: %w", err)
	}

	// Initialize schema resolver
	if err := sys.initSchemaResolver(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema resolver: %w", err)
	}

	// Initialize workflow executor
	if err := sys.initWorkflow(); err != nil {
		return nil, fmt.Errorf("failed to initialize workflow: %w", err)
	}

	return sys, nil
}

func (s *System) initLLMs() error {
	// Initialize reasoning LLM
	switch s.Config.LLM.ReasoningLLM.Provider {
	case "openai":
		provider, err := openai.NewProvider(
			s.Config.LLM.ReasoningLLM.APIKey,
			s.Config.LLM.ReasoningLLM.Model,
			&llm.Config{
				DefaultTemperature: s.Config.LLM.ReasoningLLM.DefaultTemperature,
				DefaultMaxTokens:   2000,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create reasoning LLM: %w", err)
		}
		s.ReasoningLLM = provider
	default:
		return fmt.Errorf("unsupported reasoning LLM provider: %s", s.Config.LLM.ReasoningLLM.Provider)
	}

	// Initialize fast LLM
	switch s.Config.LLM.FastLLM.Provider {
	case "openai":
		provider, err := openai.NewProvider(
			s.Config.LLM.FastLLM.APIKey,
			s.Config.LLM.FastLLM.Model,
			&llm.Config{
				DefaultTemperature: s.Config.LLM.FastLLM.DefaultTemperature,
				DefaultMaxTokens:   1000,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create fast LLM: %w", err)
		}
		s.FastLLM = provider
	default:
		return fmt.Errorf("unsupported fast LLM provider: %s", s.Config.LLM.FastLLM.Provider)
	}

	return nil
}

func (s *System) initEmbedder() error {
	switch s.Config.Embedding.Provider {
	case "openai":
		embedder, err := embedding.NewOpenAIEmbedder(
			s.Config.Embedding.APIKey,
			s.Config.Embedding.Model,
			&embedding.Config{
				BatchSize: 100,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create embedder: %w", err)
		}
		s.Embedder = embedder
	default:
		return fmt.Errorf("unsupported embedding provider: %s", s.Config.Embedding.Provider)
	}

	return nil
}

func (s *System) initVectorStore() error {
	switch s.Config.VectorStore.Type {
	case "qdrant":
		store, err := qdrant.NewStore(
			s.Config.VectorStore.Address,
			&vectorstore.Config{
				DefaultCollection: s.Config.VectorStore.DefaultCollection,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create vector store: %w", err)
		}
		s.VectorStore = store
	default:
		return fmt.Errorf("unsupported vector store type: %s", s.Config.VectorStore.Type)
	}

	return nil
}

func (s *System) initSchemaResolver() error {
	// Create resolver
	s.SchemaResolver = schema.NewResolver(s.ReasoningLLM, &schema.ResolverConfig{
		EnablePatternMatching: true,
		EnableLLMAnalysis:     true,
		EnableCaching:         true,
		CacheTTL:              3600000000000, // 1 hour in nanoseconds
	})

	return nil
}

func (s *System) initWorkflow() error {
	ctx := context.Background()

	// Create agents
	// For gpt-5 reasoning models, MaxTokens includes reasoning + output tokens
	// Need much higher limit to allow space for output after reasoning
	plannerMaxTokens := 2000
	if strings.HasPrefix(s.Config.LLM.ReasoningLLM.Model, "gpt-5") {
		plannerMaxTokens = 16000 // Reasoning models need more space
	}

	planner := agent.NewPlanner(s.ReasoningLLM, &agent.PlannerConfig{
		Temperature: s.Config.LLM.ReasoningLLM.DefaultTemperature,
		MaxTokens:   plannerMaxTokens,
	})

	rewriter := agent.NewRewriter(s.FastLLM, &agent.RewriterConfig{
		Temperature: 0.5,
		MaxTokens:   500,
	})

	supervisor := agent.NewSupervisor(s.FastLLM, &agent.SupervisorConfig{
		Temperature: 0.3,
		MaxTokens:   300,
	})

	retrieverAgent := agent.NewRetriever(
		s.VectorStore,
		s.Embedder,
		&agent.RetrieverConfig{
			DefaultTopK: s.Config.Workflow.TopKRetrieval,
		},
	)

	reranker := agent.NewReranker(&agent.RerankerConfig{
		TopN: s.Config.Workflow.TopNReranking,
	})

	// Fast LLM agents - increase tokens if using gpt-5-mini (reasoning model)
	distillerMaxTokens := 1000
	reflectorMaxTokens := 500
	policyMaxTokens := 300
	if strings.HasPrefix(s.Config.LLM.FastLLM.Model, "gpt-5") {
		distillerMaxTokens = 5000
		reflectorMaxTokens = 2500
		policyMaxTokens = 1500
	}

	distiller := agent.NewDistiller(s.FastLLM, &agent.DistillerConfig{
		Temperature: 0.3,
		MaxTokens:   distillerMaxTokens,
	})

	reflector := agent.NewReflector(s.FastLLM, &agent.ReflectorConfig{
		Temperature: 0.5,
		MaxTokens:   reflectorMaxTokens,
	})

	policy := agent.NewPolicy(s.FastLLM, &agent.PolicyConfig{
		Temperature: 0.3,
		MaxTokens:   policyMaxTokens,
	})

	// Create workflow nodes
	nodeMap := map[string]workflow.Node{
		"planner":    nodes.NewPlannerNode(ctx, planner),
		"rewriter":   nodes.NewRewriterNode(ctx, rewriter),
		"supervisor": nodes.NewSupervisorNode(ctx, supervisor),
		"retriever":  nodes.NewRetrieverNode(ctx, retrieverAgent),
		"reranker":   nodes.NewRerankerNode(ctx, reranker),
		"distiller":  nodes.NewDistillerNode(ctx, distiller),
		"reflector":  nodes.NewReflectorNode(ctx, reflector),
		"policy":     nodes.NewPolicyNode(ctx, policy),
	}

	// Build workflow graph
	graph, err := workflow.BuildDeepThinkingGraph(nodeMap)
	if err != nil {
		return fmt.Errorf("failed to build workflow graph: %w", err)
	}

	// Create executor
	s.Executor = workflow.NewExecutor(graph, &workflow.ExecutorConfig{
		Timeout: 300000000000, // 5 minutes in nanoseconds
	})

	return nil
}

// IngestDocument processes and ingests a document into the vector store.
// If deriveSchema is true, uses schema-aware chunking; otherwise uses simple paragraph chunking.
func (s *System) IngestDocument(ctx context.Context, docID string, content string, deriveSchema bool) (int, error) {
	var chunks []string
	var chunkMetadata []map[string]interface{}

	if deriveSchema && s.SchemaResolver != nil {
		// Use schema-aware chunking
		resolutionResult, err := s.SchemaResolver.Resolve(ctx, docID, content, "text/plain", nil)
		if err != nil {
			// Fall back to simple chunking if schema resolution fails
			chunks = splitIntoChunks(content, 512)
			chunkMetadata = make([]map[string]interface{}, len(chunks))
			for i := range chunks {
				chunkMetadata[i] = map[string]interface{}{"doc_id": docID}
			}
		} else {
			// Use schema-aware chunker
			chunkerConfig := chunker.DefaultConfig()
			chunkResults, err := chunker.ChunkDocument(content, resolutionResult.Schema, chunkerConfig)
			if err != nil {
				// Fall back to simple chunking
				chunks = splitIntoChunks(content, 512)
				chunkMetadata = make([]map[string]interface{}, len(chunks))
				for i := range chunks {
					chunkMetadata[i] = map[string]interface{}{"doc_id": docID}
				}
			} else {
				// Extract chunks and their metadata
				chunks = make([]string, len(chunkResults))
				chunkMetadata = make([]map[string]interface{}, len(chunkResults))
				for i, chunkResult := range chunkResults {
					chunks[i] = chunkResult.Text
					metadata := map[string]interface{}{"doc_id": docID}
					if chunkResult.Metadata != nil {
						metadata["section_id"] = chunkResult.Metadata.SectionID
						metadata["section_type"] = chunkResult.Metadata.SectionType
						metadata["hierarchy"] = chunkResult.Metadata.HierarchyPath
					}
					chunkMetadata[i] = metadata
				}
			}
		}
	} else {
		// Simple chunking: split by paragraphs
		chunks = splitIntoChunks(content, 512)
		chunkMetadata = make([]map[string]interface{}, len(chunks))
		for i := range chunks {
			chunkMetadata[i] = map[string]interface{}{"doc_id": docID}
		}
	}

	// Generate embeddings
	embedResp, err := s.Embedder.Embed(ctx, &embedding.EmbedRequest{
		Texts: chunks,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Insert into vector store
	docs := make([]vectorstore.Document, len(chunks))
	for i, chunk := range chunks {
		docs[i] = vectorstore.Document{
			ID:        uuid.New().String(), // Generate valid UUID for Qdrant
			Content:   chunk,
			Embedding: embedResp.Vectors[i].Embedding,
			Metadata:  chunkMetadata[i],
		}
	}

	_, err = s.VectorStore.Insert(ctx, &vectorstore.InsertRequest{
		CollectionName: s.Config.VectorStore.DefaultCollection,
		Documents:      docs,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert chunks: %w", err)
	}

	return len(chunks), nil
}

// splitIntoChunks splits text into chunks of approximately maxSize characters
func splitIntoChunks(text string, maxSize int) []string {
	var chunks []string
	var currentChunk string

	// Split by newlines first
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if len(currentChunk)+len(line)+1 > maxSize && len(currentChunk) > 0 {
			chunks = append(chunks, strings.TrimSpace(currentChunk))
			currentChunk = line
		} else {
			if len(currentChunk) > 0 {
				currentChunk += "\n"
			}
			currentChunk += line
		}
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk))
	}

	return chunks
}

// Close releases all system resources.
func (s *System) Close() error {
	if s.VectorStore != nil {
		return s.VectorStore.Close()
	}
	return nil
}
