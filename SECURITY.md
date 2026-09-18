# Security Policy

## Supported version

Security fixes are applied to the current `main` branch. This project is still pre-1.0, so older snapshots are not maintained as separate supported release lines.

## Reporting a vulnerability

Please report suspected vulnerabilities privately through GitHub's **Security** tab using **Report a vulnerability** when private vulnerability reporting is available for this repository. Do not include API keys, access tokens, Roblox credentials, private project data, or other secrets in a public issue, pull request, log, screenshot, or reproduction.

A useful report includes:

- the affected commit or version;
- the smallest safe reproduction you can provide;
- the security impact and expected behavior;
- relevant operating-system and Roblox Studio details; and
- whether the issue can expose credentials, alter files, execute unexpected tools, or escape the configured safety boundaries.

If private vulnerability reporting is unavailable, open a public issue containing only a non-sensitive request for a private contact path. Do not disclose exploit details or secrets in that issue.

## Credential handling

`arena-rbx` is designed to resolve the Arena API key from the environment (or a local ignored `.env` file) rather than from command-line arguments. Never commit a real API key. Use `.env.example` only as a placeholder template.

When sharing diagnostics, review them for secrets and private Roblox project content first. The CLI should not print authorization headers or API keys; a report that shows otherwise should be treated as a credential-exposure vulnerability.

## Security boundaries

Safe mode reduces accidental destructive actions, but it is not a sandbox. The configured MCP server and Roblox Studio can execute operations with the permissions of the local user and Studio session. Only connect MCP servers and use project files you trust.

High-risk or irreversible operations should require confirmation while safe mode is enabled. Reversible writes depend on the change journal having captured sufficient prior state; `/undo` is not a replacement for source control or backups.
