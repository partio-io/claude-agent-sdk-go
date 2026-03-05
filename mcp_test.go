package claude

import "testing"

func TestMcpServerJSON_Stdio(t *testing.T) {
	srv := &MCPStdioServer{
		Command: "npx",
		Args:    []string{"@modelcontextprotocol/server-postgres", "postgres://localhost/db"},
		Env:     map[string]string{"PG_USER": "admin"},
		Cwd:     "/tmp",
	}

	m := mcpServerJSON(srv)

	if m["type"] != "stdio" {
		t.Errorf("type = %v", m["type"])
	}
	if m["command"] != "npx" {
		t.Errorf("command = %v", m["command"])
	}
	args, ok := m["args"].([]string)
	if !ok || len(args) != 2 {
		t.Errorf("args = %v", m["args"])
	}
	env, ok := m["env"].(map[string]string)
	if !ok || env["PG_USER"] != "admin" {
		t.Errorf("env = %v", m["env"])
	}
	if m["cwd"] != "/tmp" {
		t.Errorf("cwd = %v", m["cwd"])
	}
}

func TestMcpServerJSON_SSE(t *testing.T) {
	srv := &MCPSSEServer{
		URL:     "http://localhost:3000/sse",
		Headers: map[string]string{"Authorization": "Bearer token"},
	}

	m := mcpServerJSON(srv)

	if m["type"] != "sse" {
		t.Errorf("type = %v", m["type"])
	}
	if m["url"] != "http://localhost:3000/sse" {
		t.Errorf("url = %v", m["url"])
	}
}

func TestMcpServerJSON_HTTP(t *testing.T) {
	srv := &MCPHTTPServer{
		URL: "http://localhost:3000",
	}

	m := mcpServerJSON(srv)

	if m["type"] != "http" {
		t.Errorf("type = %v", m["type"])
	}
	if m["url"] != "http://localhost:3000" {
		t.Errorf("url = %v", m["url"])
	}
	// No headers set, so headers key should not exist.
	if _, ok := m["headers"]; ok {
		t.Error("headers should not be set")
	}
}

func TestMcpServerJSON_StdioMinimal(t *testing.T) {
	srv := &MCPStdioServer{Command: "my-server"}
	m := mcpServerJSON(srv)

	if m["type"] != "stdio" {
		t.Errorf("type = %v", m["type"])
	}
	// Optional fields should be absent.
	if _, ok := m["args"]; ok {
		t.Error("args should not be set")
	}
	if _, ok := m["env"]; ok {
		t.Error("env should not be set")
	}
	if _, ok := m["cwd"]; ok {
		t.Error("cwd should not be set")
	}
}

func TestMCPServerConfigSealedInterface(t *testing.T) {
	var _ MCPServerConfig = (*MCPStdioServer)(nil)
	var _ MCPServerConfig = (*MCPSSEServer)(nil)
	var _ MCPServerConfig = (*MCPHTTPServer)(nil)
}
