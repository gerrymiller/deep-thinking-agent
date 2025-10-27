// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package vectorstore

import "context"

// Document represents a document with its embedding and metadata stored in the vector store.
type Document struct {
	// ID is the unique identifier for this document
	ID string

	// Content is the text content of the document
	Content string

	// Embedding is the vector representation
	Embedding []float32

	// Metadata contains additional information about the document
	Metadata map[string]interface{}

	// Score is the similarity/relevance score (set during search)
	Score float32
}

// SearchRequest contains parameters for a vector similarity search.
type SearchRequest struct {
	// Vector is the query vector to search for
	Vector []float32

	// TopK is the number of results to return
	TopK int

	// Filter allows filtering results by metadata
	// Format depends on vector store implementation
	Filter Filter

	// MinScore filters results below this similarity threshold
	MinScore float32
}

// SearchResponse contains the results of a vector search.
type SearchResponse struct {
	// Documents are the search results ordered by relevance
	Documents []Document

	// TotalResults is the total number of matching documents (before TopK limit)
	TotalResults int
}

// Filter represents metadata filters for search queries.
// The exact structure depends on the vector store implementation,
// but generally supports equality, range, and logical operations.
type Filter map[string]interface{}

// InsertRequest contains documents to insert into the vector store.
type InsertRequest struct {
	// Documents to insert
	Documents []Document

	// CollectionName is the collection/index to insert into
	CollectionName string
}

// InsertResponse contains the results of an insert operation.
type InsertResponse struct {
	// InsertedIDs are the IDs of successfully inserted documents
	InsertedIDs []string

	// Errors contains any documents that failed to insert
	Errors []InsertError
}

// InsertError represents a failed document insertion.
type InsertError struct {
	DocumentID string
	Error      error
}

// DeleteRequest contains parameters for deleting documents.
type DeleteRequest struct {
	// IDs are the document IDs to delete
	IDs []string

	// CollectionName is the collection/index to delete from
	CollectionName string

	// Filter allows deleting by metadata criteria (alternative to IDs)
	Filter Filter
}

// DeleteResponse contains the results of a delete operation.
type DeleteResponse struct {
	// DeletedCount is the number of documents successfully deleted
	DeletedCount int
}

// CollectionInfo contains metadata about a collection/index.
type CollectionInfo struct {
	// Name of the collection
	Name string

	// VectorDimension is the dimensionality of vectors in this collection
	VectorDimension int

	// DocumentCount is the number of documents in the collection
	DocumentCount int

	// Metadata contains additional collection-specific information
	Metadata map[string]interface{}
}

// Store defines the interface that all vector store implementations must provide.
// This abstraction allows swapping between Qdrant, Weaviate, Milvus, etc.
type Store interface {
	// Insert adds documents to the vector store.
	Insert(ctx context.Context, req *InsertRequest) (*InsertResponse, error)

	// Search performs a vector similarity search.
	Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)

	// Delete removes documents from the vector store.
	Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error)

	// Get retrieves specific documents by ID.
	Get(ctx context.Context, collectionName string, ids []string) ([]Document, error)

	// CreateCollection creates a new collection/index with specified dimensions.
	CreateCollection(ctx context.Context, name string, dimension int, metadata map[string]interface{}) error

	// DeleteCollection removes an entire collection/index.
	DeleteCollection(ctx context.Context, name string) error

	// ListCollections returns information about all collections.
	ListCollections(ctx context.Context) ([]CollectionInfo, error)

	// GetCollection returns information about a specific collection.
	GetCollection(ctx context.Context, name string) (*CollectionInfo, error)

	// Close closes the connection to the vector store.
	Close() error

	// Name returns the vector store implementation name.
	Name() string
}

// Config contains common configuration for vector store implementations.
type Config struct {
	// Type specifies which vector store to use
	Type string

	// Address is the connection address (host:port)
	Address string

	// APIKey for authentication (if required)
	APIKey string

	// Timeout in seconds for operations
	TimeoutSeconds int

	// DefaultCollection is the collection to use when not specified
	DefaultCollection string

	// Additional provider-specific settings
	Extra map[string]interface{}
}
