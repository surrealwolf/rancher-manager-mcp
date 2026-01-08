package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type ToolHandler func(ctx context.Context, args map[string]interface{}) (interface{}, error)

type Server struct {
	name        string
	version     string
	tools       map[string]Tool
	toolHandlers map[string]ToolHandler
	mu          sync.RWMutex
}

func NewServer(name, version string) *Server {
	return &Server{
		name:         name,
		version:      version,
		tools:        make(map[string]Tool),
		toolHandlers: make(map[string]ToolHandler),
	}
}

func (s *Server) RegisterTool(name, description string, handler ToolHandler) {
	s.RegisterToolWithSchema(name, description, nil, handler)
}

func (s *Server) RegisterToolWithSchema(name, description string, inputSchema map[string]interface{}, handler ToolHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if inputSchema == nil {
		inputSchema = map[string]interface{}{
			"type":       "object",
			"properties": make(map[string]interface{}),
		}
	}

	s.tools[name] = Tool{
		Name:        name,
		Description: description,
		InputSchema: inputSchema,
	}
	s.toolHandlers[name] = handler
}

func (s *Server) HandleRequest(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	return s.handleRequest(ctx, req)
}

func (s *Server) Serve(ctx context.Context, stdin io.Reader, stdout io.Writer) error {
	decoder := json.NewDecoder(stdin)
	encoder := json.NewEncoder(stdout)
	encoder.SetEscapeHTML(false)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var req JSONRPCRequest
			if err := decoder.Decode(&req); err != nil {
				if err == io.EOF {
					return nil
				}
				// Send error response for malformed requests
				if req.ID != nil {
					errorResp := &JSONRPCResponse{
						JSONRPC: "2.0",
						ID:      req.ID,
						Error: &JSONRPCError{
							Code:    -32700,
							Message: "Parse error",
						},
					}
					encoder.Encode(errorResp)
				}
				continue
			}

			// Skip notifications (requests without ID)
			if req.ID == nil {
				continue
			}

			resp := s.handleRequest(ctx, &req)
			// Ensure ID is always set in response
			if resp.ID == nil {
				resp.ID = req.ID
			}
			if err := encoder.Encode(resp); err != nil {
				// Can't log to stderr in stdio mode, just continue
				continue
			}
		}
	}
}

func (s *Server) handleRequest(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	// Ensure ID is always preserved
	responseID := req.ID
	if responseID == nil {
		// Use null for notifications (though we shouldn't get here)
		responseID = nil
	}

	switch req.Method {
	case "initialize":
		resp := s.handleInitialize(req)
		resp.ID = responseID
		return resp
	case "tools/list":
		resp := s.handleToolsList(req)
		resp.ID = responseID
		return resp
	case "tools/call":
		resp := s.handleToolsCall(ctx, req)
		resp.ID = responseID
		return resp
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      responseID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

func (s *Server) handleInitialize(req *JSONRPCRequest) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: InitializeResponse{
			ProtocolVersion: "2024-11-05",
			Capabilities: ServerCapabilities{
				Tools: &ToolsCapability{
					ListChanged: true,
				},
			},
			ServerInfo: ServerInfo{
				Name:    s.name,
				Version: s.version,
			},
		},
	}
}

func (s *Server) handleToolsList(req *JSONRPCRequest) *JSONRPCResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tools := make([]Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: ToolListResponse{
			Tools: tools,
		},
	}
}

func (s *Server) handleToolsCall(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	var callReq CallToolRequest
	if req.Params == nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	name, ok := req.Params["name"].(string)
	if !ok {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid params: name is required",
			},
		}
	}

	callReq.Name = name
	if args, ok := req.Params["arguments"].(map[string]interface{}); ok {
		callReq.Arguments = args
	}

	s.mu.RLock()
	handler, exists := s.toolHandlers[name]
	s.mu.RUnlock()

	if !exists {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: fmt.Sprintf("Tool not found: %s", name),
			},
		}
	}

	result, err := handler(ctx, callReq.Arguments)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: CallToolResponse{
				Content: []Content{
					{
						Type: "text",
						Text: fmt.Sprintf("Error: %v", err),
					},
				},
				IsError: true,
			},
		}
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Internal error: %v", err),
			},
		}
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: CallToolResponse{
			Content: []Content{
				{
					Type: "text",
					Text: string(resultJSON),
				},
			},
		},
	}
}
