# Arena Roblox MCP — Design Specification

**Date:** 2026-09-02  
**Repository:** `AnnonyA/arena-roblox-mcp`  
**Status:** Approved design, pre-implementation  
**Primary language:** Go

## 1. Purpose

`arena-roblox-mcp` is a lightweight local CLI agent that uses Arena.ai models as the reasoning layer and Roblox Studio's official MCP server as the tool/execution layer.

The user interacts with one local binary, `arena-rbx`, writes natural-language tasks, and the agent can inspect Roblox Studio, read and search Luau code, edit scripts, run playtests, inspect console output, and continue iterating until the task is complete or a configured safety/iteration limit is reached.

The first release prioritizes Windows and Roblox Studio, while keeping the internal Arena and MCP layers sufficiently decoupled to avoid architectural lock-in.

## 2. Primary Goals

The v0.1 architecture must:

1. Run as a single lightweight Go binary with no Node.js, Python, Electron, or other runtime requirement.
2. Connect to Arena.ai through its OpenAI-compatible API.
3. Discover Arena models dynamically rather than hardcoding a model list.
4. Stream model output to the terminal as soon as it arrives.
5. Connect to Roblox Studio's official MCP server over persistent stdio.
6. Discover and cache MCP tool schemas at session startup.
7. Execute multi-round agentic tool-calling loops.
8. Read, search, inspect, edit, and debug Roblox/Luau projects through MCP.
9. Support playtest-driven iteration and console feedback.
10. Keep user credentials out of source control, output, and logs.
11. Record reversible edits so `/diff` and `/undo` can work safely.
12. Degrade gracefully on API, MCP, model, tool, or Studio failures.
13. Remain responsive on low-end hardware by keeping local CPU/RAM overhead low.
14. Be testable mostly without a real Arena account or a running Roblox Studio instance.

## 3. Explicit Non-Goals for v0.1

The first version will not include:

- A desktop GUI or Electron application.
- A web UI.
- Multi-agent orchestration.
- Judge-model arbitration.
- Automatic model routing based on task classification.
- A generic multi-provider gateway beyond Arena.ai.
- Full project snapshots or source-control replacement.
- Auto-updating installers.
- A Roblox Studio plugin of our own.
- A requirement to support every operating system equally on day one.

The architecture may leave room for these capabilities, but they do not block v0.1.

## 4. High-Level Architecture

```text
User
  |
  v
CLI
  |
  v
Context / Session Manager
  |
  v
Agent Loop <----------------------------+
  |                                     |
  +--> Arena Client --> Arena.ai         |
  |          |                          |
  |          +-- stream text/tool calls |
  |                                     |
  +--> Tool Dispatcher                  |
             |                          |
             v                          |
        MCP Client                      |
             |                          |
             v                          |
     Roblox Studio MCP                  |
             |                          |
             v                          |
        Roblox Studio ------------------+
```

The CLI is only a presentation and command layer. Business logic lives in isolated packages so that future frontends can reuse the same agent core.

## 5. Core Components

### 5.1 Arena Client

Responsibilities:

- Read the Arena API base URL and API key configuration.
- Reuse one HTTP client and transport for the process lifetime.
- Fetch `/v1/models` dynamically.
- Submit chat-completion requests.
- Support streaming responses.
- Parse text deltas and tool-call deltas incrementally.
- Send MCP-derived tools in Arena-compatible tool definitions.
- Handle 401, 429, 5xx, transport failures, malformed stream chunks, cancellation, and timeouts.
- Respect retry metadata such as `Retry-After` when available.
- Capture resolved-model/fallback metadata where Arena exposes it.
- Never log authorization headers or the API key.

The Arena package must not depend on Roblox-specific code.

### 5.2 MCP Client

Responsibilities:

- Spawn and maintain the configured MCP server process.
- Use a persistent stdio transport.
- Perform MCP initialization and capability negotiation.
- Discover available tools once per Studio connection/session.
- Cache tool schemas.
- Execute tools with structured arguments.
- Return structured results and errors to the agent.
- Shut down cleanly on exit and cancellation.
- Avoid process leaks/zombie child processes.

The MCP package should remain generic enough that Roblox-specific behavior sits in a separate adapter layer.

### 5.3 Roblox Adapter

Responsibilities:

- Supply the default Windows MCP command for Roblox Studio.
- Understand Studio-specific connection/session information.
- Manage `studio_id` when one or multiple Studio instances are available.
- Classify Roblox tools by purpose and risk.
- Provide tool-family filtering hints to the agent.
- Recognize editing operations that can be snapshotted and reversed.
- Support playtest and console-oriented agent workflows.

When exactly one Studio is available, selection should be automatic. When multiple Studios are available, the CLI must expose selection rather than guessing silently.

### 5.4 Agent Loop

Responsibilities:

- Accept a user task and current session context.
- Ask Arena for the next response/action.
- Stream ordinary text to the CLI.
- Detect tool calls.
- Validate tool name and arguments.
- Dispatch tools through MCP.
- Append compact tool results to the conversation.
- Continue until:
  - the model returns a final answer,
  - the user cancels,
  - the maximum tool-round limit is reached,
  - a non-recoverable error occurs,
  - or a safety confirmation is declined.

Default maximum tool rounds: **12**, configurable.

The loop must not allow unbounded autonomous iteration.

### 5.5 Context Manager

Responsibilities:

- Maintain conversation state for the active CLI session.
- Keep recent user/model/tool events.
- Prevent giant tool outputs from accumulating indefinitely.
- Prefer targeted script reads/searches instead of reading entire projects.
- Compact or omit stale large payloads while keeping decisions and relevant findings.
- Clear context on `/clear` without destroying edit history needed for `/diff` or `/undo`.

v0.1 will favor simple, low-overhead budgeting rather than a heavyweight always-resident tokenizer.

### 5.6 Session / Change Journal

Responsibilities:

- Record tool actions performed by the agent.
- Record reversible mutations with enough previous state to undo them.
- Track affected script/instance identifiers and before/after content where appropriate.
- Produce human-readable diffs.
- Avoid storing full-game snapshots.

On Windows, session data should default under a user-local application directory such as:

```text
%LOCALAPPDATA%\arena-rbx\sessions\
```

Session journals must never contain the Arena API key.

## 6. CLI Design

Binary name:

```text
arena-rbx.exe
```

Normal startup:

```text
arena-rbx
```

Expected startup presentation:

```text
Arena Roblox MCP
────────────────────────────
Arena      connected
Studio     connected
Model      <selected model>
Session    default

> _
```

No full-screen TUI is required. ANSI formatting may be used when supported, with a plain-text fallback.

### 6.1 Slash Commands

The initial CLI command surface is:

```text
/model              change model
/models             list Arena models
/studio             select Studio instance
/status             show Arena/MCP/Studio state
/tools              show available MCP tools
/history            show session actions/tool calls
/diff               show recorded changes
/undo               revert the latest supported reversible change
/clear              clear conversational context
/config              show effective non-secret configuration
/help                show commands
/exit                exit cleanly
```

Any non-empty line that does not start with `/` is treated as an agent task.

### 6.2 Cancellation

`Ctrl+C` while a request is active should cancel the current Arena/tool operation and return to the prompt when safe.

A subsequent `Ctrl+C` or explicit `/exit` may terminate the process cleanly.

Child MCP processes must not be left running accidentally.

## 7. Configuration

Default config filename:

```text
arena-rbx.json
```

Example shape:

```json
{
  "arena": {
    "apiKeyEnv": "ARENA_API_KEY",
    "model": "auto",
    "fallbacks": [],
    "stream": true
  },
  "agent": {
    "maxToolRounds": 12,
    "autoPlaytest": true,
    "contextBudget": "balanced",
    "safeMode": true
  },
  "mcpServers": {
    "Roblox_Studio": {
      "command": "cmd.exe",
      "args": [
        "/c",
        "%LOCALAPPDATA%\\Roblox\\mcp.bat"
      ]
    }
  }
}
```

Configuration precedence:

```text
CLI flags
  > environment variables
  > arena-rbx.json
  > internal defaults
```

### 7.1 Credentials

The repository will contain `.env.example`, never a real `.env`.

Example:

```text
ARENA_API_KEY=put_your_arena_api_key_here
```

`.env` must be ignored by Git.

The program should also support ordinary environment-variable configuration, for example in PowerShell:

```powershell
$env:ARENA_API_KEY="..."
arena-rbx.exe
```

No command such as `/config`, `/status`, debug logging, HTTP error formatting, or panic path may print the API key.

## 8. Dynamic Model Discovery

`/models` must query Arena rather than relying on a hardcoded table.

The selected model can come from config, CLI flags, or `/model`.

A future automatic router may use `/model auto`, but in v0.1 `auto` must not imply a complex custom model-ranking subsystem unless Arena itself provides a directly usable automatic route. If no usable automatic route is available during implementation, the CLI must require a concrete model selection rather than invent hidden routing logic.

## 9. MCP Tool Discovery and Tool Filtering

MCP tool schemas are discovered once per active connection and cached.

The implementation should classify tools into logical families, for example:

- code search/read
- game-tree inspection
- editing
- Luau execution
- playtest
- console/debugging
- screenshots
- assets or other heavier operations

The agent should avoid repeatedly rebuilding or rediscovering identical schemas.

Where the Arena request format allows it, the agent should send only the tool families relevant to the task or current phase, while always preserving enough core tools to recover from an incorrect initial classification.

Correctness is more important than over-aggressive tool suppression.

## 10. Performance Design

The user requirement is a CLI that feels fast and adds minimal local overhead, including on modest PCs.

The project cannot guarantee zero latency on every computer because model-provider latency, internet quality, Roblox Studio load, and project size are external factors. The design target is therefore: **our client must not add material avoidable latency**.

Required practices:

- One reusable Arena HTTP client/transport.
- Persistent MCP stdio connection.
- Streaming output as soon as deltas arrive.
- No Electron, browser runtime, Node, or Python dependency.
- Lazy initialization where practical.
- MCP schemas cached after discovery.
- Bounded buffers.
- No aggressive polling loops.
- Goroutines used for clear concurrency benefits, not gratuitously.
- Logging optional and bounded.
- Context kept compact.
- No reconnect-to-Studio cycle for every prompt.

Performance benchmarks should include at least:

```text
BenchmarkCLIStartup
BenchmarkToolRegistry
BenchmarkContextProcessing
BenchmarkStreamParser
```

Benchmark numbers should inform regression detection rather than use unrealistic CI-hardcoded absolute thresholds across heterogeneous runners.

## 11. Safety Model

### 11.1 Tool Risk Classes

Tools/actions are classified into three broad levels:

**READ**

Examples: script reads, searches, grep, game-tree inspection, console reads.

These may execute automatically.

**REVERSIBLE_WRITE**

Examples: script edits or property edits for which the previous state can be reliably captured.

These may execute automatically in safe mode only after the change journal records sufficient pre-change state.

**HIGH_RISK**

Examples: destructive deletes, broad irreversible operations, mutation-heavy arbitrary Luau execution, or operations whose previous state cannot be restored reliably.

These require explicit confirmation while safe mode is enabled.

### 11.2 Safe Mode

Safe mode defaults to `true`.

Default behavior:

```text
reads                     automatic
reversible edits          automatic after snapshot
playtests                 automatic
high-risk operations      confirmation required
```

A future or implementation-time flag such as `--allow-dangerous` may disable selected confirmations for advanced users, but must not be the default.

### 11.3 Undo Semantics

`/undo` guarantees only the latest journaled action that is explicitly marked reversible.

The CLI must never imply that an irreversible action can be restored.

If the latest action cannot be reversed, the user must be told clearly.

## 12. Diff and Change Journal

Before a reversible edit:

```text
read/capture prior state
  -> write journal entry
  -> execute mutation
  -> capture resulting state
  -> finalize diff entry
```

`/diff` shows session mutations using human-readable unified-style diffs where appropriate.

Only changed resources are journaled. The application does not maintain a giant shadow copy of the whole Roblox place.

## 13. Error Handling

Errors should be concise, actionable, and non-secret-bearing.

Examples:

### Arena key absent

```text
Arena API key not configured. Set ARENA_API_KEY or configure arena.apiKeyEnv.
```

### Arena authentication failure

```text
Arena authentication failed. Check your API key.
```

The raw key or Authorization header must never be included.

### Rate limiting

Respect server-provided retry timing when available, perform bounded retries, then surface a clear failure.

### MCP / Roblox Studio unavailable

```text
No Roblox Studio MCP session detected.
```

The CLI remains usable enough to show status/help rather than crashing.

### Tool failure

Tool errors are returned to the agent in structured form so the model can revise its approach when safe and within the tool-round limit.

### Iteration limit

When `maxToolRounds` is reached, the loop stops and reports what was completed and what remains unresolved.

### Cancellation

Cancellation propagates through HTTP, agent, and MCP layers using Go contexts where possible.

## 14. Repository Layout

Planned structure:

```text
arena-roblox-mcp/
|
├── cmd/
│   └── arena-rbx/
│       └── main.go
|
├── internal/
│   ├── arena/
│   │   ├── client.go
│   │   ├── models.go
│   │   ├── chat.go
│   │   └── stream.go
│   │
│   ├── agent/
│   │   ├── agent.go
│   │   ├── loop.go
│   │   ├── context.go
│   │   └── tools.go
│   │
│   ├── mcp/
│   │   ├── client.go
│   │   ├── transport.go
│   │   └── registry.go
│   │
│   ├── roblox/
│   │   ├── studio.go
│   │   ├── discovery.go
│   │   └── tools.go
│   │
│   ├── cli/
│   │   ├── cli.go
│   │   ├── commands.go
│   │   └── renderer.go
│   │
│   ├── config/
│   │   ├── config.go
│   │   └── env.go
│   │
│   └── session/
│       ├── session.go
│       ├── history.go
│       └── changes.go
|
├── tests/
├── docs/
│   └── superpowers/
│       └── specs/
|
├── arena-rbx.example.json
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

`main.go` should contain composition/startup logic only, not application business logic.

## 15. Testing Strategy

Most behavior must be testable without real external services.

### 15.1 Fake Arena Server

Use an in-process HTTP test server to cover:

- authentication header handling
- model discovery
- normal chat responses
- streaming
- fragmented/malformed chunks
- tool-call deltas
- 401
- 429
- Retry-After
- 5xx
- retries
- cancellation
- resolved-model/fallback metadata
- secret redaction

### 15.2 Fake MCP Server

Use a deterministic test transport/server to cover:

- initialize
- tool discovery
- schema caching
- tool execution
- tool errors
- disconnects
- reconnect/shutdown behavior
- cancellation

### 15.3 Agent Tests

Cover:

- single response without tools
- one tool call
- multiple sequential tool calls
- playtest/console feedback loop
- max tool rounds
- malformed tool arguments
- unknown tool requests
- context compaction
- tool-result truncation/compaction
- cancellation
- model/tool loop prevention

### 15.4 CLI Tests

Cover:

- slash-command parsing
- `/model`
- `/models`
- `/studio`
- `/status`
- `/tools`
- `/history`
- `/diff`
- `/undo`
- `/clear`
- `/config`
- `/help`
- `/exit`
- plain-text output without ANSI
- configuration precedence
- graceful Ctrl+C behavior where practical to test

### 15.5 Safety Tests

Cover:

- snapshot occurs before reversible writes
- reversible changes produce diffs
- `/undo` restores prior state
- irreversible operations are not mislabeled reversible
- safe-mode confirmation gates high-risk actions
- API keys are redacted from logs/errors/status/config output

### 15.6 Required Verification Commands

At appropriate implementation checkpoints:

```text
go test ./...
go vet ./...
go test -race ./...
```

Race tests may be limited to platforms/configurations where supported reliably, but concurrency-sensitive packages must be exercised with the race detector during development/CI where possible.

## 16. Manual Integration Test Flow

The intended user test flow is:

1. Clone/download the repo.
2. Copy `.env.example` to `.env` or set `ARENA_API_KEY` through the shell.
3. Insert the user's own Arena API key locally.
4. Open Roblox Studio.
5. Start `arena-rbx.exe`.
6. Run `/status`.
7. Run `/models`.
8. Select a model.
9. Ask the agent to identify scripts related to a known feature.
10. Ask the agent to read and explain a target script.
11. Request a small reversible modification.
12. Inspect `/diff`.
13. Test `/undo`.
14. Request a modification followed by playtest/debugging.
15. Run an end-to-end autonomous bug-fix task.

The user's API key remains local and is never committed to GitHub.

## 17. v0.1.0 Acceptance Criteria

v0.1.0 is not considered complete until all of these are satisfied:

- [ ] Builds as a standalone Windows executable.
- [ ] Does not require Node.js or Python.
- [ ] Arena API connectivity works.
- [ ] Arena model discovery is dynamic.
- [ ] Model selection works.
- [ ] Streaming output works.
- [ ] Roblox MCP connection works.
- [ ] MCP tool discovery and caching work.
- [ ] Multi-round tool calling works.
- [ ] Script reading/searching works through Studio MCP.
- [ ] Reversible script editing works.
- [ ] Playtest invocation works.
- [ ] Console output can be fed back to the model.
- [ ] `/status` works.
- [ ] `/models` works.
- [ ] `/model` works.
- [ ] `/diff` works for supported edits.
- [ ] `/undo` works for supported edits.
- [ ] `/clear` works.
- [ ] Errors are understandable and do not dump secrets.
- [ ] Ctrl+C/shutdown does not leave orphaned MCP child processes.
- [ ] Arena API keys never appear in normal logs/status/config output.
- [ ] Unit/integration tests pass.
- [ ] Race-sensitive code has been checked with the Go race detector where supported.
- [ ] Windows release build passes.
- [ ] README documents setup from a clean machine.

A required real-world demonstration before calling the release viable:

```text
user task
  -> Arena reasons
  -> agent investigates project
  -> identifies issue
  -> edits Luau
  -> launches playtest
  -> reads console/result
  -> revises if necessary
  -> completes successfully
```

## 18. Release and Build Direction

v0.1 prioritizes Windows because Roblox Studio MCP is the immediate target.

The Go codebase should avoid unnecessary platform coupling so later CI can produce builds such as:

```text
arena-rbx-windows-amd64.exe
arena-rbx-windows-arm64.exe
arena-rbx-linux-amd64
arena-rbx-linux-arm64
arena-rbx-darwin-amd64
arena-rbx-darwin-arm64
```

Cross-platform artifacts are not a v0.1 completion requirement unless Roblox/MCP support on those platforms is verified during implementation.

## 19. Roadmap After v0.1

### v0.2 direction

- stronger context compaction
- persistent/resumable sessions
- `/resume`
- improved diff presentation
- smarter tool-family selection

### v0.3 direction

- intelligent `/model auto`
- task-aware model selection
- richer fallback policies
- coding/debugging/architecture profiles

### Later direction

- multi-agent investigation
- judge model
- parallel analysis
- checkpoints
- richer multi-Studio workflows
- optional GUI built on the same core

These are roadmap directions, not commitments required by the initial implementation plan.

## 20. Design Principles

Implementation decisions should preserve the following principles:

1. **Lightweight first.** Local work should be cheap; network/model/Studio operations are the expected expensive parts.
2. **One responsibility per package.** Arena, MCP, Roblox, agent, CLI, config, and session concerns remain separable.
3. **Persistent connections over repeated setup.** Avoid reconnect/reinitialize overhead when possible.
4. **Stream everything practical.** The user should see useful output early.
5. **Safe autonomy.** Reads and reversible changes may flow automatically; destructive operations need explicit handling.
6. **No secret leakage.** Credentials are never treated as ordinary config output.
7. **Test through interfaces.** Fake Arena and fake MCP layers make agent behavior reproducible.
8. **No speculative complexity.** Multi-agent, GUI, custom routing, and unrelated providers wait until the core is proven.
9. **Roblox-focused product, decoupled internals.** The public project is for Roblox Studio, while internal boundaries avoid unnecessary Roblox/Arena entanglement.
10. **Evidence before release.** v0.1 requires automated verification plus a real Arena → MCP → Studio bug-fix workflow.
