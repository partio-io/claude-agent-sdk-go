package claude

import (
	"context"
	"testing"
)

func TestApplyOptions(t *testing.T) {
	cfg := applyOptions([]Option{
		WithModel("claude-opus-4-6"),
		WithSystemPrompt("You are a helpful assistant"),
		WithCwd("/tmp"),
		WithVerbose(true),
		WithMaxTurns(10),
		WithMaxBudgetUSD(1.50),
		WithPermissionMode("bypassPermissions"),
		WithIncludePartialMessages(true),
		WithNoSessionPersistence(true),
	})

	if cfg.model != "claude-opus-4-6" {
		t.Errorf("model = %q", cfg.model)
	}
	if cfg.systemPrompt != "You are a helpful assistant" {
		t.Errorf("systemPrompt = %q", cfg.systemPrompt)
	}
	if cfg.cwd != "/tmp" {
		t.Errorf("cwd = %q", cfg.cwd)
	}
	if !cfg.verbose {
		t.Error("verbose should be true")
	}
	if cfg.maxTurns != 10 {
		t.Errorf("maxTurns = %d", cfg.maxTurns)
	}
	if cfg.maxBudgetUSD != 1.50 {
		t.Errorf("maxBudgetUSD = %f", cfg.maxBudgetUSD)
	}
	if cfg.permissionMode != "bypassPermissions" {
		t.Errorf("permissionMode = %q", cfg.permissionMode)
	}
	if !cfg.includePartialMessages {
		t.Error("includePartialMessages should be true")
	}
	if !cfg.noSessionPersistence {
		t.Error("noSessionPersistence should be true")
	}
}

func TestWithEnv(t *testing.T) {
	cfg := applyOptions([]Option{
		WithEnv("FOO", "bar"),
		WithEnv("BAZ", "qux"),
	})
	if cfg.env["FOO"] != "bar" {
		t.Errorf("env[FOO] = %q", cfg.env["FOO"])
	}
	if cfg.env["BAZ"] != "qux" {
		t.Errorf("env[BAZ] = %q", cfg.env["BAZ"])
	}
}

func TestWithAllowedTools(t *testing.T) {
	cfg := applyOptions([]Option{
		WithAllowedTools("Read", "Write"),
		WithAllowedTools("Bash"),
	})
	if len(cfg.allowedTools) != 3 {
		t.Errorf("allowedTools length = %d, want 3", len(cfg.allowedTools))
	}
}

func TestWithDisallowedTools(t *testing.T) {
	cfg := applyOptions([]Option{
		WithDisallowedTools("Bash"),
	})
	if len(cfg.disallowedTools) != 1 {
		t.Errorf("disallowedTools length = %d, want 1", len(cfg.disallowedTools))
	}
}

func TestWithMCPServer(t *testing.T) {
	cfg := applyOptions([]Option{
		WithMCPServer("postgres", &MCPStdioServer{
			Command: "npx",
			Args:    []string{"@modelcontextprotocol/server-postgres"},
		}),
	})
	if cfg.mcpServers == nil {
		t.Fatal("mcpServers is nil")
	}
	srv, ok := cfg.mcpServers["postgres"]
	if !ok {
		t.Fatal("postgres server not found")
	}
	stdio, ok := srv.(*MCPStdioServer)
	if !ok {
		t.Fatalf("expected *MCPStdioServer, got %T", srv)
	}
	if stdio.Command != "npx" {
		t.Errorf("Command = %q", stdio.Command)
	}
}

func TestWithHook(t *testing.T) {
	matcher := "Bash"
	cfg := applyOptions([]Option{
		WithHook(HookPreToolUse, HookMatcher{
			Matcher: &matcher,
			Handler: func(_ context.Context, input HookCallbackInput) (HookOutput, error) {
				return HookOutput{Decision: "allow"}, nil
			},
		}),
	})
	if cfg.hooks == nil {
		t.Fatal("hooks is nil")
	}
	matchers, ok := cfg.hooks[HookPreToolUse]
	if !ok {
		t.Fatal("PreToolUse hooks not found")
	}
	if len(matchers) != 1 {
		t.Errorf("matchers length = %d, want 1", len(matchers))
	}
}

func TestWithAgent(t *testing.T) {
	maxTurns := 5
	cfg := applyOptions([]Option{
		WithAgent("searcher", AgentDefinition{
			Model:        "claude-haiku-4-5",
			SystemPrompt: "You search code",
			MaxTurns:     &maxTurns,
		}),
	})
	if cfg.agents == nil {
		t.Fatal("agents is nil")
	}
	agent, ok := cfg.agents["searcher"]
	if !ok {
		t.Fatal("searcher agent not found")
	}
	if agent.Model != "claude-haiku-4-5" {
		t.Errorf("Model = %q", agent.Model)
	}
}

func TestWithResume(t *testing.T) {
	cfg := applyOptions([]Option{
		WithResume("session-123"),
		WithForkSession(true),
	})
	if cfg.resume != "session-123" {
		t.Errorf("resume = %q", cfg.resume)
	}
	if !cfg.forkSession {
		t.Error("forkSession should be true")
	}
}

func TestWithCLIPath(t *testing.T) {
	cfg := applyOptions([]Option{
		WithCLIPath("/usr/local/bin/claude"),
	})
	if cfg.cliPath != "/usr/local/bin/claude" {
		t.Errorf("cliPath = %q", cfg.cliPath)
	}
}
