package mcp

import (
	"context"
	"log/slog"

	"github.com/JuanVilla424/qwen-pet/internal/ai"
	"github.com/JuanVilla424/qwen-pet/internal/kb"
	"github.com/JuanVilla424/qwen-pet/internal/pet"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server with qwen-pet dependencies.
type Server struct {
	engine   *ai.Engine
	store    *kb.Store
	animal   *pet.Animal
	petState *pet.State
}

// NewServer creates a new MCP server wrapper.
func NewServer(engine *ai.Engine, store *kb.Store, animal *pet.Animal, petState *pet.State) *Server {
	return &Server{
		engine:   engine,
		store:    store,
		animal:   animal,
		petState: petState,
	}
}

// Run starts the MCP server on stdio transport (blocks until client disconnects).
func (s *Server) Run(ctx context.Context, version string) error {
	server := mcpsdk.NewServer(
		&mcpsdk.Implementation{
			Name:    "qwen-pet",
			Version: version,
		},
		nil,
	)

	s.registerTools(server)

	slog.Info("starting MCP server on stdio")
	return server.Run(ctx, &mcpsdk.StdioTransport{})
}
