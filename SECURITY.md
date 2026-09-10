# Security Policy

## Reporting a vulnerability

Please do not open a public issue for a vulnerability that could expose credentials, execute unintended Studio actions, or otherwise put users at risk.

Instead, use GitHub's private vulnerability reporting for this repository when available. Include the affected version or commit, a concise reproduction, the expected and observed behavior, and any relevant logs with secrets removed.

Do not include Arena API keys, authorization headers, Roblox credentials, session tokens, or other private data in reports.

## Credential handling

`arena-rbx` is designed to read the Arena API key from the environment variable configured by `arena.apiKeyEnv`, with local `.env` support for development convenience. Real credentials must never be committed to the repository.

When sharing logs or command output, review them for secrets before posting. The project should not print API keys or authorization headers, but reports should still be sanitized before they are made public.

## Safe testing

Security reports should use fake or disposable data whenever possible. Do not test destructive MCP or Roblox Studio operations against projects you cannot safely restore.

For high-risk tool behavior, prefer unit or integration tests with fakes/mocks. A real Roblox Studio test should only be used when necessary and when the affected place can be recovered.

## Supported versions

Security fixes are applied to the current `main` branch while the project is pre-1.0. Once tagged releases have a formal support window, this section will be updated with the supported release lines.
