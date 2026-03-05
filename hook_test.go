package claude

import "testing"

func TestHookEventConstants(t *testing.T) {
	// Verify all hook event constants have expected values.
	events := map[HookEvent]string{
		HookPreToolUse:         "PreToolUse",
		HookPostToolUse:        "PostToolUse",
		HookPostToolUseFailure: "PostToolUseFailure",
		HookStop:               "Stop",
		HookSubagentStart:      "SubagentStart",
		HookSubagentStop:       "SubagentStop",
		HookUserPromptSubmit:   "UserPromptSubmit",
		HookPreCompact:         "PreCompact",
		HookNotification:       "Notification",
	}

	for event, want := range events {
		if string(event) != want {
			t.Errorf("HookEvent %q != %q", event, want)
		}
	}
}

func TestHookMatcherConfig(t *testing.T) {
	matcher := "Bash|Edit"
	timeout := 30.0
	cfg := hookMatcherConfig{
		Matcher:         &matcher,
		HookCallbackIDs: []string{"hook_1", "hook_2"},
		Timeout:         &timeout,
	}
	if *cfg.Matcher != "Bash|Edit" {
		t.Errorf("Matcher = %q", *cfg.Matcher)
	}
	if len(cfg.HookCallbackIDs) != 2 {
		t.Errorf("HookCallbackIDs length = %d", len(cfg.HookCallbackIDs))
	}
	if *cfg.Timeout != 30.0 {
		t.Errorf("Timeout = %f", *cfg.Timeout)
	}
}

func TestHookMatcherNilMatcher(t *testing.T) {
	cfg := hookMatcherConfig{
		Matcher:         nil, // matches all
		HookCallbackIDs: []string{"hook_1"},
	}
	if cfg.Matcher != nil {
		t.Error("Matcher should be nil (match all)")
	}
}
