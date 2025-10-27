package qdrant

import (
	"context"
	"errors"
	"fmt"

	"deep-thinking-agent/pkg/vectorstore"

	"github.com/google/uuid"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Store implements the vectorstore.Store interface for Qdrant.
type Store struct {
	client      pb.PointsClient
	collections pb.CollectionsClient
	conn        *grpc.ClientConn
	config      *vectorstore.Config
}

// NewStore creates a new Qdrant vector store instance.
// address: Qdrant server address (e.g., "localhost:6334")
// config: Configuration options (can be nil for defaults)
func NewStore(address string, config *vectorstore.Config) (*Store, error) {
	if address == "" {
		return nil, errors.New("Qdrant address is required")
	}

	// Apply default config if not provided
	if config == nil {
		config = &vectorstore.Config{
			Type:              "qdrant",
			Address:           address,
			TimeoutSeconds:    30,
			DefaultCollection: "documents",
		}
	}

	// Create gRPC connection
	// Note: In production, use proper TLS credentials
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Qdrant: %w", err)
	}

	// Create clients
	pointsClient := pb.NewPointsClient(conn)
	collectionsClient := pb.NewCollectionsClient(conn)

	return &Store{
		client:      pointsClient,
		collections: collectionsClient,
		conn:        conn,
		config:      config,
	}, nil
}

// Insert adds documents to the vector store.
func (s *Store) Insert(ctx context.Context, req *vectorstore.InsertRequest) (*vectorstore.InsertResponse, error) {
	if req == nil {
		return nil, errors.New("insert request cannot be nil")
	}
	if len(req.Documents) == 0 {
		return nil, errors.New("no documents to insert")
	}

	collectionName := req.CollectionName
	if collectionName == "" {
		collectionName = s.config.DefaultCollection
	}

	// Convert documents to Qdrant points
	points := make([]*pb.PointStruct, 0, len(req.Documents))
	insertedIDs := make([]string, 0, len(req.Documents))

	for _, doc := range req.Documents {
		// Generate ID if not provided
		id := doc.ID
		if id == "" {
			id = uuid.New().String()
		}

		// Convert metadata to payload
		payload := make(map[string]*pb.Value)
		payload["content"] = &pb.Value{
			Kind: &pb.Value_StringValue{StringValue: doc.Content},
		}

		// Add all metadata fields
		for k, v := range doc.Metadata {
			payload[k] = convertToQdrantValue(v)
		}

		point := &pb.PointStruct{
			Id: &pb.PointId{
				PointIdOptions: &pb.PointId_Uuid{Uuid: id},
			},
			Vectors: &pb.Vectors{
				VectorsOptions: &pb.Vectors_Vector{
					Vector: &pb.Vector{Data: doc.Embedding},
				},
			},
			Payload: payload,
		}

		points = append(points, point)
		insertedIDs = append(insertedIDs, id)
	}

	// Upsert points
	_, err := s.client.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}

	return &vectorstore.InsertResponse{
		InsertedIDs: insertedIDs,
		Errors:      []vectorstore.InsertError{},
	}, nil
}

// Search performs a vector similarity search.
func (s *Store) Search(ctx context.Context, req *vectorstore.SearchRequest) (*vectorstore.SearchResponse, error) {
	if req == nil {
		return nil, errors.New("search request cannot be nil")
	}
	if len(req.Vector) == 0 {
		return nil, errors.New("search vector cannot be empty")
	}

	collectionName := s.config.DefaultCollection

	// Build search request
	searchReq := &pb.SearchPoints{
		CollectionName: collectionName,
		Vector:         req.Vector,
		Limit:          uint64(req.TopK),
		WithPayload:    &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
		ScoreThreshold: &req.MinScore,
	}

	// Add filter if provided
	if req.Filter != nil && len(req.Filter) > 0 {
		searchReq.Filter = convertToQdrantFilter(req.Filter)
	}

	// Execute search
	resp, err := s.client.Search(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	// Convert results
	documents := make([]vectorstore.Document, 0, len(resp.Result))
	for _, hit := range resp.Result {
		doc := vectorstore.Document{
			ID:       hit.Id.GetUuid(),
			Score:    hit.Score,
			Metadata: make(map[string]interface{}),
		}

		// Extract content and metadata from payload
		if hit.Payload != nil {
			if contentVal, ok := hit.Payload["content"]; ok {
				doc.Content = contentVal.GetStringValue()
			}

			// Convert all payload fields to metadata
			for k, v := range hit.Payload {
				if k != "content" {
					doc.Metadata[k] = convertFromQdrantValue(v)
				}
			}
		}

		documents = append(documents, doc)
	}

	return &vectorstore.SearchResponse{
		Documents:    documents,
		TotalResults: len(documents),
	}, nil
}

// Delete removes documents from the vector store.
func (s *Store) Delete(ctx context.Context, req *vectorstore.DeleteRequest) (*vectorstore.DeleteResponse, error) {
	if req == nil {
		return nil, errors.New("delete request cannot be nil")
	}

	collectionName := req.CollectionName
	if collectionName == "" {
		collectionName = s.config.DefaultCollection
	}

	var pointsSelector *pb.PointsSelector

	if len(req.IDs) > 0 {
		// Delete by IDs
		uuids := make([]string, len(req.IDs))
		copy(uuids, req.IDs)

		pointsSelector = &pb.PointsSelector{
			PointsSelectorOneOf: &pb.PointsSelector_Points{
				Points: &pb.PointsIdsList{
					Ids: convertToQdrantIDs(uuids),
				},
			},
		}
	} else if req.Filter != nil {
		// Delete by filter
		pointsSelector = &pb.PointsSelector{
			PointsSelectorOneOf: &pb.PointsSelector_Filter{
				Filter: convertToQdrantFilter(req.Filter),
			},
		}
	} else {
		return nil, errors.New("either IDs or Filter must be provided")
	}

	// Execute delete
	resp, err := s.client.Delete(ctx, &pb.DeletePoints{
		CollectionName: collectionName,
		Points:         pointsSelector,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to delete documents: %w", err)
	}

	return &vectorstore.DeleteResponse{
		DeletedCount: int(resp.Result.GetOperationId()),
	}, nil
}

// Get retrieves specific documents by ID.
func (s *Store) Get(ctx context.Context, collectionName string, ids []string) ([]vectorstore.Document, error) {
	if collectionName == "" {
		collectionName = s.config.DefaultCollection
	}
	if len(ids) == 0 {
		return []vectorstore.Document{}, nil
	}

	// Retrieve points
	resp, err := s.client.Get(ctx, &pb.GetPoints{
		CollectionName: collectionName,
		Ids:            convertToQdrantIDs(ids),
		WithPayload:    &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
		WithVectors:    &pb.WithVectorsSelector{SelectorOptions: &pb.WithVectorsSelector_Enable{Enable: true}},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	// Convert results
	documents := make([]vectorstore.Document, 0, len(resp.Result))
	for _, point := range resp.Result {
		doc := vectorstore.Document{
			ID:       point.Id.GetUuid(),
			Metadata: make(map[string]interface{}),
		}

		// Extract vector
		if vector := point.Vectors.GetVector(); vector != nil {
			doc.Embedding = vector.Data
		}

		// Extract content and metadata
		if point.Payload != nil {
			if contentVal, ok := point.Payload["content"]; ok {
				doc.Content = contentVal.GetStringValue()
			}

			for k, v := range point.Payload {
				if k != "content" {
					doc.Metadata[k] = convertFromQdrantValue(v)
				}
			}
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// CreateCollection creates a new collection/index with specified dimensions.
func (s *Store) CreateCollection(ctx context.Context, name string, dimension int, metadata map[string]interface{}) error {
	if name == "" {
		return errors.New("collection name is required")
	}
	if dimension <= 0 {
		return errors.New("dimension must be positive")
	}

	_, err := s.collections.Create(ctx, &pb.CreateCollection{
		CollectionName: name,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     uint64(dimension),
					Distance: pb.Distance_Cosine,
				},
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

// DeleteCollection removes an entire collection/index.
func (s *Store) DeleteCollection(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("collection name is required")
	}

	_, err := s.collections.Delete(ctx, &pb.DeleteCollection{
		CollectionName: name,
	})

	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	return nil
}

// ListCollections returns information about all collections.
func (s *Store) ListCollections(ctx context.Context) ([]vectorstore.CollectionInfo, error) {
	resp, err := s.collections.List(ctx, &pb.ListCollectionsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	collections := make([]vectorstore.CollectionInfo, 0, len(resp.Collections))
	for _, col := range resp.Collections {
		collections = append(collections, vectorstore.CollectionInfo{
			Name:     col.Name,
			Metadata: make(map[string]interface{}),
		})
	}

	return collections, nil
}

// GetCollection returns information about a specific collection.
func (s *Store) GetCollection(ctx context.Context, name string) (*vectorstore.CollectionInfo, error) {
	if name == "" {
		return nil, errors.New("collection name is required")
	}

	resp, err := s.collections.Get(ctx, &pb.GetCollectionInfoRequest{
		CollectionName: name,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	info := &vectorstore.CollectionInfo{
		Name:          name, // Use the requested name
		DocumentCount: int(*resp.Result.PointsCount),
		Metadata:      make(map[string]interface{}),
	}

	// Extract vector dimension
	if params := resp.Result.Config.Params.VectorsConfig.GetParams(); params != nil {
		info.VectorDimension = int(params.Size)
	}

	return info, nil
}

// Close closes the connection to the vector store.
func (s *Store) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// Name returns the vector store implementation name.
func (s *Store) Name() string {
	return "qdrant"
}

// Helper functions for type conversion

func convertToQdrantValue(v interface{}) *pb.Value {
	switch val := v.(type) {
	case string:
		return &pb.Value{Kind: &pb.Value_StringValue{StringValue: val}}
	case int:
		return &pb.Value{Kind: &pb.Value_IntegerValue{IntegerValue: int64(val)}}
	case int64:
		return &pb.Value{Kind: &pb.Value_IntegerValue{IntegerValue: val}}
	case float64:
		return &pb.Value{Kind: &pb.Value_DoubleValue{DoubleValue: val}}
	case bool:
		return &pb.Value{Kind: &pb.Value_BoolValue{BoolValue: val}}
	default:
		// Default to string representation
		return &pb.Value{Kind: &pb.Value_StringValue{StringValue: fmt.Sprintf("%v", val)}}
	}
}

func convertFromQdrantValue(v *pb.Value) interface{} {
	if v == nil {
		return nil
	}

	switch kind := v.Kind.(type) {
	case *pb.Value_StringValue:
		return kind.StringValue
	case *pb.Value_IntegerValue:
		return kind.IntegerValue
	case *pb.Value_DoubleValue:
		return kind.DoubleValue
	case *pb.Value_BoolValue:
		return kind.BoolValue
	default:
		return nil
	}
}

func convertToQdrantIDs(ids []string) []*pb.PointId {
	result := make([]*pb.PointId, len(ids))
	for i, id := range ids {
		result[i] = &pb.PointId{
			PointIdOptions: &pb.PointId_Uuid{Uuid: id},
		}
	}
	return result
}

func convertToQdrantFilter(filter vectorstore.Filter) *pb.Filter {
	// Basic filter conversion - can be extended for more complex filters
	conditions := make([]*pb.Condition, 0, len(filter))

	for key, value := range filter {
		condition := &pb.Condition{
			ConditionOneOf: &pb.Condition_Field{
				Field: &pb.FieldCondition{
					Key: key,
					Match: &pb.Match{
						MatchValue: &pb.Match_Keyword{
							Keyword: fmt.Sprintf("%v", value),
						},
					},
				},
			},
		}
		conditions = append(conditions, condition)
	}

	return &pb.Filter{
		Must: conditions,
	}
}
