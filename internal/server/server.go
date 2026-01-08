package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
	"github.com/rancher/rancher-manager-mcp/internal/server/handlers"
)

type Server struct {
	rancherURL   string
	rancherToken string
	client       *client.RancherClient
	mcpServer    *mcp.Server
}

func NewServer(rancherURL, rancherToken string, insecureSkipVerify bool) *Server {
	s := &Server{
		rancherURL:   rancherURL,
		rancherToken: rancherToken,
	}

	// Initialize Rancher client
	if rancherURL != "" && rancherToken != "" {
		s.client = client.NewRancherClient(rancherURL, rancherToken, insecureSkipVerify)
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
		"status":             "ok",
		"rancher_configured": s.client != nil,
	}
	json.NewEncoder(w).Encode(status)
}

func (s *Server) registerTools() {
	// Register tools from handler modules
	handlers.RegisterClusterTools(s.mcpServer, s.client)
	handlers.RegisterUserTools(s.mcpServer, s.client)
	handlers.RegisterProjectTools(s.mcpServer, s.client)
	handlers.RegisterAuditPolicyTools(s.mcpServer, s.client)
	handlers.RegisterKubeconfigTools(s.mcpServer, s.client)
	handlers.RegisterTokenTools(s.mcpServer, s.client)
	handlers.RegisterGlobalRoleTools(s.mcpServer, s.client)
	handlers.RegisterGlobalRoleBindingTools(s.mcpServer, s.client)
	handlers.RegisterRoleTemplateTools(s.mcpServer, s.client)
	handlers.RegisterClusterRoleTemplateBindingTools(s.mcpServer, s.client)
	handlers.RegisterProjectRoleTemplateBindingTools(s.mcpServer, s.client)

	// Register status tools
	handlers.RegisterClusterStatusTools(s.mcpServer, s.client)
	handlers.RegisterUserStatusTools(s.mcpServer, s.client)
	handlers.RegisterProjectStatusTools(s.mcpServer, s.client)
	handlers.RegisterRoleStatusTools(s.mcpServer, s.client)
	handlers.RegisterBindingStatusTools(s.mcpServer, s.client)
	handlers.RegisterAuditPolicyStatusTools(s.mcpServer, s.client)

	// Register create and update tools
	handlers.RegisterClusterCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterUserCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterProjectCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterAuditPolicyCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterKubeconfigCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterTokenCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterGlobalRoleCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterGlobalRoleBindingCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterRoleTemplateCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterClusterRoleTemplateBindingCreateUpdateTools(s.mcpServer, s.client)
	handlers.RegisterProjectRoleTemplateBindingCreateUpdateTools(s.mcpServer, s.client)

	// Register delete tools
	handlers.RegisterClusterDeleteTools(s.mcpServer, s.client)
	handlers.RegisterUserDeleteTools(s.mcpServer, s.client)
	handlers.RegisterProjectDeleteTools(s.mcpServer, s.client)
	handlers.RegisterAuditPolicyDeleteTools(s.mcpServer, s.client)
	handlers.RegisterKubeconfigDeleteTools(s.mcpServer, s.client)
	handlers.RegisterTokenDeleteTools(s.mcpServer, s.client)
	handlers.RegisterGlobalRoleDeleteTools(s.mcpServer, s.client)
	handlers.RegisterGlobalRoleBindingDeleteTools(s.mcpServer, s.client)
	handlers.RegisterRoleTemplateDeleteTools(s.mcpServer, s.client)
	handlers.RegisterClusterRoleTemplateBindingDeleteTools(s.mcpServer, s.client)
	handlers.RegisterProjectRoleTemplateBindingDeleteTools(s.mcpServer, s.client)
}

func (s *Server) registerResources() {
	// Register resources if needed
}
