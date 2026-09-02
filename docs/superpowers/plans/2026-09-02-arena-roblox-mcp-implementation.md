# Arena Roblox MCP Implementation Plan

**Goal:** Build the approved lightweight Go CLI connecting Arena.ai to Roblox Studio through MCP.

**Spec:** `docs/superpowers/specs/2026-09-02-arena-roblox-mcp-design.md`

## Tasks

- [ ] 1. Bootstrap Go module and secure configuration with tests.
- [ ] 2. Add Arena HTTP/model/streaming client with fake-server tests.
- [ ] 3. Add persistent official MCP Go SDK wrapper and cached tool registry.
- [ ] 4. Add Roblox Studio adapter and Studio selection.
- [ ] 5. Add bounded session history, diffs, and reversible edit records.
- [ ] 6. Add bounded agent loop, context compaction, and tool dispatch.
- [ ] 7. Add interactive CLI commands, renderer, and lifecycle handling.
- [ ] 8. Add fake end-to-end integration tests, CI, README, benchmarks, and Windows build verification.

Each task follows TDD: write a failing test, verify the failure, implement the smallest passing change, run tests and vet, then commit a small change. Release verification also runs the race detector where supported.
