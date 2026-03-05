// Package testutil provides mock subprocess helpers and NDJSON fixtures for testing.
package testutil

// NDJSON fixtures captured from real Claude CLI sessions.

// SystemInit is a system init message.
const SystemInit = `{"type":"system","subtype":"init","session_id":"550e8400-e29b-41d4-a716-446655440000","uuid":"sys-001","cwd":"/home/user/project","model":"claude-sonnet-4-6","tools":["Read","Write","Edit","Bash","Glob","Grep"],"mcp_servers":[],"permissionMode":"default","apiKeySource":"anthropic","slash_commands":[],"agents":{},"claude_code_version":"2.1.39"}`

// AssistantTextOnly is an assistant message with only text content.
const AssistantTextOnly = `{"type":"assistant","uuid":"msg_01XFG","session_id":"550e8400-e29b-41d4-a716-446655440000","parent_tool_use_id":null,"message":{"id":"msg_01XFG","type":"message","role":"assistant","model":"claude-sonnet-4-6","content":[{"type":"text","text":"The answer is 4."}],"stop_reason":"end_turn","usage":{"input_tokens":100,"output_tokens":20,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}}`

// AssistantToolUse is an assistant message with a tool use content block.
const AssistantToolUse = `{"type":"assistant","uuid":"msg_02ABC","session_id":"550e8400-e29b-41d4-a716-446655440000","parent_tool_use_id":null,"message":{"id":"msg_02ABC","type":"message","role":"assistant","model":"claude-sonnet-4-6","content":[{"type":"text","text":"I'll read the file for you."},{"type":"tool_use","id":"toolu_01ABC","name":"Read","input":{"file_path":"/home/user/project/main.go"}}],"stop_reason":"tool_use","usage":{"input_tokens":200,"output_tokens":50,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}}`

// UserToolResult is a user message with a tool result.
const UserToolResult = `{"type":"user","uuid":"usr_01DEF","session_id":"550e8400-e29b-41d4-a716-446655440000","parent_tool_use_id":null,"message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_01ABC","content":"package main\n\nfunc main() {}","is_error":false}]}}`

// ResultSuccess is a successful result message.
const ResultSuccess = `{"type":"result","subtype":"success","session_id":"550e8400-e29b-41d4-a716-446655440000","is_error":false,"result":"The answer is 4.","num_turns":1,"duration_ms":1234,"duration_api_ms":890,"total_cost_usd":0.0023,"usage":{"input_tokens":100,"output_tokens":20,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}`

// ResultMaxTurns is a result message that hit max turns.
const ResultMaxTurns = `{"type":"result","subtype":"error_max_turns","session_id":"550e8400-e29b-41d4-a716-446655440000","is_error":true,"result":"Reached maximum turns limit","num_turns":10,"duration_ms":5000,"duration_api_ms":4000}`

// StreamEventTextDelta is a stream event with a text delta.
const StreamEventTextDelta = `{"type":"stream_event","uuid":"evt_01GHI","session_id":"550e8400-e29b-41d4-a716-446655440000","parent_tool_use_id":null,"event":{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello "}}}`

// AssistantThinking is an assistant message with a thinking block.
const AssistantThinking = `{"type":"assistant","uuid":"msg_03THK","session_id":"550e8400-e29b-41d4-a716-446655440000","parent_tool_use_id":null,"message":{"id":"msg_03THK","type":"message","role":"assistant","model":"claude-sonnet-4-6","content":[{"type":"thinking","thinking":"Let me analyze this...","signature":"sig_abc"},{"type":"text","text":"Here is my analysis."}],"stop_reason":"end_turn","usage":{"input_tokens":300,"output_tokens":100,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}}`

// ControlRequestCanUseTool is a permission request from the CLI.
const ControlRequestCanUseTool = `{"type":"control_request","request_id":"req_2_def456","request":{"subtype":"can_use_tool","tool_name":"Bash","input":{"command":"ls -la","description":"List files"}}}`

// ControlRequestHookCallback is a hook callback request from the CLI.
const ControlRequestHookCallback = `{"type":"control_request","request_id":"req_3_ghi789","request":{"subtype":"hook_callback","callback_id":"hook_PreToolUse_0","input":{"session_id":"550e8400","transcript_path":"/tmp/transcript","cwd":"/home/user","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"ls"}}}}`

// SimpleSession is a complete simple session: system init → assistant → result.
var SimpleSession = []string{
	SystemInit,
	AssistantTextOnly,
	ResultSuccess,
}

// ToolUseSession is a session with tool use: init → assistant(tool) → user(result) → assistant → result.
var ToolUseSession = []string{
	SystemInit,
	AssistantToolUse,
	UserToolResult,
	AssistantTextOnly,
	ResultSuccess,
}
