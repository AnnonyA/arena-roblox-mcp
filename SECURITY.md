# Security Policy

## Reporting a vulnerability

Please do not open a public issue for a vulnerability that could expose credentials, execute unintended Studio actions, or otherwise put users at risk.

Instead, use GitHub's private vulnerability reporting for this repository when available. Include the affected version or commit, a concise reproduction, the expected and observed behavior, and any relevant logs with secrets removed.

Do not include Arena API keys, authorization headers, Roblox credentials, session tokens, or other private data in reports.

## Credential handling

`arena-rbx` is designed to read the Arena API key from the environment variable configured by `arena.apiKeyEnv`, with local `.env` support for development convenience. Real credentials must never be committed to the repository.

When sharing logs or command output, review them for secrets before posting. The project should not print API keys or authorization headers, but reports should still be sanitized before they are made public.

## Terminal output safety

Treat text originating outside the CLI itself as untrusted terminal data. Model IDs, Studio metadata, MCP tool descriptions, session history, diffs, configuration values, and propagated errors may contain control or Unicode format characters even when their upstream source is normally trusted.

Terminal-facing dynamic text should pass through the CLI display sanitizers rather than being written directly. Single-line fields remove control and Unicode format characters. Multiline output uses `cli.WriteSafeMultiline`, which preserves intentional newlines and tabs while removing other control and Unicode format characters. This prevents dynamic data from injecting terminal escape sequences or visually hiding/reordering content.

When adding a new output path, include a regression test with representative terminal-control and invisible-Unicode input. Fuzz coverage for the shared sanitizers should remain enabled. Static CLI-owned formatting may still be written directly when it contains no dynamic data.

## Safe testing

Security reports should use fake or disposable data whenever possible. Do not test destructive MCP or Roblox Studio operations against projects you cannot safely restore.

For high-risk tool behavior, prefer unit or integration tests with fakes/mocks. A real Roblox Studio test should only be used when necessary and when the affected place can be recovered.

## Security boundaries

Safe mode reduces accidental destructive actions, but it is not a sandbox. The configured MCP server and Roblox Studio MCP process operate with the permissions available to the local user and Studio session, so only connect MCP servers and open project files you trust.

High-risk or irreversible operations should require confirmation while safe mode is enabled. Reversible writes depend on the change journal capturing sufficient prior state; `/undo` is not a replacement for source control or backups.

## Supported versions

Security fixes are applied to the current `main` branch while the project is pre-1.0. Once tagged releases have a formal support window, this section will be updated with the supported release lines.
