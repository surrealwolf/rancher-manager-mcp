package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

type Server struct {
	rancherURL   string
	rancherToken string
	client       *RancherClient
	mcpServer    *mcp.Server
}

func NewServer(rancherURL, rancherToken string, insecureSkipVerify bool) *Server {
	s := &Server{
		rancherURL:   rancherURL,
		rancherToken: rancherToken,
	}

	// Initialize Rancher client
	if rancherURL != "" && rancherToken != "" {
		s.client = NewRancherClient(rancherURL, rancherToken, insecureSkipVerify)
	}

	// Initialize MCP server
	s.mcpServer = mcp.NewServer("rancher-manager-mcp", "1.0.0")
	s.registerTools()
	s.registerResources()

	return s
}

func (s *Server) ServeStdio(ctx context.Context) error {
	return s.mcpServer.Serve(ctx, os.Stdin, os.Stdout)
}

func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", s.handleMCPRequest)
	mux.HandleFunc("/health", s.handleHealth)
	return mux
}

func (s *Server) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request: %v", err), http.StatusBadRequest)
		return
	}

	var mcpRequest mcp.JSONRPCRequest
	if err := json.Unmarshal(body, &mcpRequest); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request: %v", err), http.StatusBadRequest)
		return
	}

	// Handle MCP request
	response := s.mcpServer.HandleRequest(r.Context(), &mcpRequest)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := map[string]interface{}{
		"status": "ok",
		"rancher_configured": s.client != nil,
	}
	json.NewEncoder(w).Encode(status)
}

func (s *Server) registerTools() {
	// List clusters
	s.mcpServer.RegisterTool("list_clusters", "List all Rancher clusters", s.listClusters)

	// Get cluster details
	s.mcpServer.RegisterTool("get_cluster", "Get details of a specific cluster", s.getCluster)

	// List users
	s.mcpServer.RegisterTool("list_users", "List all Rancher users", s.listUsers)

	// Get user details
	s.mcpServer.RegisterTool("get_user", "Get details of a specific user", s.getUser)

	// List projects
	s.mcpServer.RegisterTool("list_projects", "List all Rancher projects", s.listProjects)

	// Get project details
	s.mcpServer.RegisterTool("get_project", "Get details of a specific project", s.getProject)
}

func (s *Server) registerResources() {
	// Register resources if needed
}

// Tool handlers
func (s *Server) listClusters(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	return s.client.ListClusters(ctx)
}

func (s *Server) getCluster(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return s.client.GetCluster(ctx, name)
}

func (s *Server) listUsers(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	return s.client.ListUsers(ctx)
}

func (s *Server) getUser(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return s.client.GetUser(ctx, name)
}

func (s *Server) listProjects(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	return s.client.ListProjects(ctx)
}

func (s *Server) getProject(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	namespace, _ := args["namespace"].(string)
	return s.client.GetProject(ctx, name, namespace)
}
