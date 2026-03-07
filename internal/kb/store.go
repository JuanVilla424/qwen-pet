package kb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/JuanVilla424/qwen-pet/internal/config"
	chromem "github.com/philippgille/chromem-go"
)

// Result represents a KB query result.
type Result struct {
	ID         string
	Content    string
	Metadata   map[string]string
	Similarity float32
}

// Store wraps chromem-go for knowledge base operations.
type Store struct {
	db         *chromem.DB
	collection *chromem.Collection
}

// NewStore creates a new persistent knowledge base.
func NewStore(cfg config.KBCfg, embFunc chromem.EmbeddingFunc) (*Store, error) {
	if err := os.MkdirAll(cfg.Path, 0o755); err != nil {
		return nil, fmt.Errorf("create KB directory: %w", err)
	}

	db, err := chromem.NewPersistentDB(cfg.Path, cfg.Compress)
	if err != nil {
		return nil, fmt.Errorf("open persistent DB: %w", err)
	}

	collection, err := db.GetOrCreateCollection(cfg.Collection, nil, embFunc)
	if err != nil {
		return nil, fmt.Errorf("get or create collection: %w", err)
	}

	return &Store{
		db:         db,
		collection: collection,
	}, nil
}

// Add stores a document in the knowledge base.
func (s *Store) Add(ctx context.Context, id, content string, metadata map[string]string) error {
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata["added_at"] = time.Now().Format(time.RFC3339)

	doc := chromem.Document{
		ID:       id,
		Content:  content,
		Metadata: metadata,
	}

	if err := s.collection.AddDocument(ctx, doc); err != nil {
		return fmt.Errorf("add document: %w", err)
	}
	return nil
}

// Query performs a semantic similarity search.
func (s *Store) Query(ctx context.Context, query string, nResults int) ([]Result, error) {
	return s.QueryFiltered(ctx, query, nResults, nil)
}

// QueryFiltered performs a semantic search with metadata filters.
func (s *Store) QueryFiltered(ctx context.Context, query string, nResults int, where map[string]string) ([]Result, error) {
	if s.collection.Count() == 0 {
		return nil, nil
	}

	results, err := s.collection.Query(ctx, query, nResults, where, nil)
	if err != nil {
		return nil, fmt.Errorf("query collection: %w", err)
	}

	out := make([]Result, len(results))
	for i, r := range results {
		out[i] = Result{
			ID:         r.ID,
			Content:    r.Content,
			Metadata:   r.Metadata,
			Similarity: r.Similarity,
		}
	}
	return out, nil
}

// Delete removes documents by ID.
func (s *Store) Delete(ctx context.Context, ids ...string) error {
	if err := s.collection.Delete(ctx, nil, nil, ids...); err != nil {
		return fmt.Errorf("delete documents: %w", err)
	}
	return nil
}

// Count returns the number of documents in the collection.
func (s *Store) Count() int {
	return s.collection.Count()
}
