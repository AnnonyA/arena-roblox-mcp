# arena-roblox-mcp

`arena-roblox-mcp` is a lightweight Go CLI that connects Arena.ai models with Roblox Studio through the official MCP workflow.

> **Status:** v0.1 implementation plan complete; ongoing work focuses on incremental quality, performance, reliability, tests, documentation, security, and CLI experience within the approved design. The project is still pre-release, so review the current implementation before relying on it for important Roblox projects.

## Requirements

- Go 1.23+
- Windows for the primary Roblox Studio workflow
- Roblox Studio with its MCP integration available
- An Arena API key in `ARENA_API_KEY`

## Build

```powershell
git clone https://github.com/AnnonyA/arena-roblox-mcp.git
cd arena-roblox-mcp
go build -o arena-rbx.exe ./cmd/arena-rbx
```

Run it with:

```powershell
$env:ARENA_API_KEY="your-key-here"
.\arena-rbx.exe
```

Do not commit real API keys. The repository includes `.env.example` only; `.env` is ignored by Git.

## Configuration

`arena-rbx` reads `arena-rbx.json` from the current directory when present. If the file is absent, built-in defaults are used.

Example:

```json
{
  "arena": {
    "apiKeyEnv": "ARENA_API_KEY",
    "model": "",
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

You can also choose a model for one run:

```powershell
.\arena-rbx.exe --model <arena-model-id>
```

Without a configured model, use `/models` to discover models from Arena and `/model <id>` to select one.

## CLI commands

```text
/model              choose an Arena model
/models             list Arena models
/studio             select a Roblox Studio session
/status             show Arena/MCP/Studio state
/tools              show discovered MCP tools
/history            show session actions
/diff               show recorded changes
/undo               revert the latest supported reversible change
/clear              clear conversational context
/config              show effective non-secret configuration
/help                show commands
/exit                exit cleanly
```

Any non-empty line that does not start with `/` is treated as a task.

## Security

- Keep `ARENA_API_KEY` outside source control.
- Never put a real key in `arena-rbx.json`, issues, logs, screenshots, or commits.
- Safe mode defaults to enabled.
- Read operations may run automatically. Reversible writes are journaled with prior state so supported changes can be inspected with `/diff` and reverted with `/undo`.
- High-risk operations require explicit confirmation while safe mode is enabled; `--allow-dangerous` is an explicit process-lifetime opt-out and is never enabled by default or persisted silently.
- `/undo` only guarantees the latest journaled action explicitly marked reversible; it does not imply that irreversible actions can be restored.

## Development

Run the same core checks used by CI:

```bash
go test ./...
go vet ./...
go test -race ./...
```

Build the Windows CLI with:

```powershell
go build -o arena-rbx.exe ./cmd/arena-rbx
```

The project intentionally keeps the Arena, MCP, Roblox, agent, session, and CLI layers separated so they can be tested mostly with fakes instead of requiring a live Arena account or Roblox Studio for every test.

## Project direction

The approved v0.1 target is a bounded local agent loop that can stream Arena responses, discover and invoke Roblox Studio MCP tools, iterate through tool results, keep context compact, journal reversible edits, and stop safely on cancellation, errors, confirmation boundaries, or the configured tool-round limit.

See `docs/superpowers/specs/2026-09-02-arena-roblox-mcp-design.md` and `docs/superpowers/plans/2026-09-02-arena-roblox-mcp-implementation.md` for the current design and implementation plan.