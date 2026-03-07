package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type askPetInput struct {
	Question string `json:"question" jsonschema:"The question to ask the pet"`
	Context  string `json:"context,omitempty" jsonschema:"Additional context for the question"`
}

type checkPreferenceInput struct {
	Category string `json:"category" jsonschema:"Preference category (e.g. stack, conventions, ui)"`
	Key      string `json:"key" jsonschema:"Specific preference key to look up"`
}

type storeDecisionInput struct {
	Category string `json:"category" jsonschema:"Decision category (e.g. architecture, patterns, preferences)"`
	Key      string `json:"key" jsonschema:"Decision identifier"`
	Value    string `json:"value" jsonschema:"The decision or preference value"`
	Reason   string `json:"reason,omitempty" jsonschema:"Why this decision was made"`
}

type approveArtifactInput struct {
	Type        string `json:"type" jsonschema:"Artifact type (e.g. wireframe, architecture, schema)"`
	Description string `json:"description" jsonschema:"Description of what needs approval"`
	FilePath    string `json:"file_path,omitempty" jsonschema:"Path to the artifact file"`
}

func (s *Server) registerTools(server *mcpsdk.Server) {
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "ask_pet",
		Description: "Ask the pet a question. It searches its knowledge base, reasons with AI, or escalates via Telegram.",
	}, s.handleAskPet)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "check_preference",
		Description: "Check a specific user preference or convention from the knowledge base.",
	}, s.handleCheckPreference)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "store_decision",
		Description: "Store a decision or preference in the knowledge base for future reference.",
	}, s.handleStoreDecision)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "approve_artifact",
		Description: "Request user approval for an artifact via Telegram.",
	}, s.handleApproveArtifact)
}

func (s *Server) handleAskPet(ctx context.Context, _ *mcpsdk.CallToolRequest, in askPetInput) (*mcpsdk.CallToolResult, any, error) {
	if in.Question == "" {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: "Error: question is required"}},
			IsError: true,
		}, nil, nil
	}

	slog.Info("ask_pet", "question", in.Question)

	answer, err := s.engine.Decide(ctx, in.Question, in.Context)
	if err != nil {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
			IsError: true,
		}, nil, nil
	}

	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: answer.Text}},
	}, nil, nil
}

func (s *Server) handleCheckPreference(ctx context.Context, _ *mcpsdk.CallToolRequest, in checkPreferenceInput) (*mcpsdk.CallToolResult, any, error) {
	if in.Category == "" || in.Key == "" {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: "Error: category and key are required"}},
			IsError: true,
		}, nil, nil
	}

	query := fmt.Sprintf("%s %s", in.Category, in.Key)
	where := map[string]string{"category": in.Category}

	results, err := s.store.QueryFiltered(ctx, query, 3, where)
	if err != nil {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
			IsError: true,
		}, nil, nil
	}

	if len(results) == 0 {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: fmt.Sprintf("No preference found for %s/%s", in.Category, in.Key)}},
		}, nil, nil
	}

	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: results[0].Content}},
	}, nil, nil
}

func (s *Server) handleStoreDecision(ctx context.Context, _ *mcpsdk.CallToolRequest, in storeDecisionInput) (*mcpsdk.CallToolResult, any, error) {
	if in.Category == "" || in.Key == "" || in.Value == "" {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: "Error: category, key, and value are required"}},
			IsError: true,
		}, nil, nil
	}

	id := fmt.Sprintf("decision-%s-%s-%d", in.Category, in.Key, time.Now().UnixNano())
	content := fmt.Sprintf("%s/%s: %s", in.Category, in.Key, in.Value)
	if in.Reason != "" {
		content += fmt.Sprintf(" (reason: %s)", in.Reason)
	}

	meta := map[string]string{
		"category": in.Category,
		"key":      in.Key,
		"type":     "decision",
	}

	if err := s.store.Add(ctx, id, content, meta); err != nil {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: fmt.Sprintf("Error storing: %v", err)}},
			IsError: true,
		}, nil, nil
	}

	emoji := s.animal.Emoji
	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{&mcpsdk.TextContent{
			Text: fmt.Sprintf("%s Stored! %s/%s = %s", emoji, in.Category, in.Key, in.Value),
		}},
	}, nil, nil
}

func (s *Server) handleApproveArtifact(ctx context.Context, _ *mcpsdk.CallToolRequest, in approveArtifactInput) (*mcpsdk.CallToolResult, any, error) {
	if in.Type == "" || in.Description == "" {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: "Error: type and description are required"}},
			IsError: true,
		}, nil, nil
	}

	// This tool always escalates to Telegram for human approval
	question := fmt.Sprintf("[APPROVAL REQUEST]\nType: %s\nDescription: %s", in.Type, in.Description)
	if in.FilePath != "" {
		question += fmt.Sprintf("\nFile: %s", in.FilePath)
	}

	answer, err := s.engine.Decide(ctx, question, "")
	if err != nil {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
			IsError: true,
		}, nil, nil
	}

	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: answer.Text}},
	}, nil, nil
}
