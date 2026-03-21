# theseus

Starter repository for the Theseus migration setup.

## MCP Setup

The repo uses `Context7` as an MCP-backed documentation source for framework and library context.

Add it with:

```bash
codex mcp add context7 -- npx -y @upstash/context7-mcp
```

Current project assets:

- `.codex/agents/` for Codex agent definitions such as `analyzer.toml`, `spec-writer.toml`, and `architect.toml`
- `.agents/skills/` for reusable Agent Skills directories (`skill-name/SKILL.md`)
- `docker-compose.yml` and `config.json` for the Sourcebot bootstrap setup

## Sourcebot check

Start the local stack:

```bash
docker compose up -d
```

Run the repository health check:

```bash
./scripts/check-sourcebot.sh
```

Exit codes:

- `0`: Sourcebot reachable at `http://localhost:3000`
- `2`: Docker CLI not available
- `3`: `sourcebot` container not running
- `4`: endpoint not reachable while running (or missing `curl`)

## Current subagent flow

`Analyzer -> Spec-Writer -> Verifier -> Architect -> Builder`

- `Analyzer` reads source code, README guidance, and dependency information.
- `Architect` plans the target architecture and selects suitable packages for the Builder with Context7-backed documentation.
