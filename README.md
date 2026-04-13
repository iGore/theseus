# theseus

Starter repository for the Theseus migration setup.

## MCP Setup

The repo uses two MCPs across Codex and OpenCode:

- `Context7` for framework and library documentation
- `Sourcebot` for repository and code-intelligence access via `http://localhost:3000/api/mcp`

Add `Context7` for Codex with:

```bash
codex mcp add context7 -- npx -y @upstash/context7-mcp
```

Current project assets:

- `.codex/agents/` for Codex agent definitions such as `analyzer.toml`, `spec-writer.toml`, and `architect.toml`
- `.opencode/agents/` for OpenCode agent definitions that mirror the Codex role set in Markdown frontmatter format
- `opencode.jsonc` for OpenCode MCP and default-agent configuration
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

`Analyzer -> Spec-Writer -> Verifier -> Architect -> Planner -> Task Decomposer -> Builder`

- This now mirrors the full conceptual model from the report: a migration pipeline coordinated by one `Orchestrator` across reconstruction, transformation planning, and re-implementation.
- `Analyzer` reconstructs the bounded source module from code, repository metadata, and supporting docs.
- `Spec-Writer` produces the central textual specification artifact.
- `Verifier` gates the pipeline by checking whether the specification is traceable back to source-context evidence.
- `Architect` plans the target architecture and selects suitable packages with Context7-backed documentation.
- `Planner` derives an ordered migration plan with work packages from the verified specification and architecture.
- `Task Decomposer` turns that plan into implementation-ready coding tasks.
- `Builder` implements from the verified specification, architecture outputs, and decomposed work packages, then records build/test evidence.

## OpenCode setup

The repo now ships the same migration roles for OpenCode under `.opencode/agents/`:

- `orchestrator.md` as the primary agent
- `analyzer.md` as a subagent
- `spec-writer.md` as a subagent
- `verifier.md` as a subagent
- `architect.md` as a subagent
- `planner.md` as a subagent
- `task-decomposer.md` as a subagent
- `builder.md` as a subagent

The textual specification is the main handoff artifact between reconstruction and implementation. `AGENTS.md` acts as the shared rule layer across all roles, matching the thesis' emphasis on a persistent instruction artifact for build hints, conventions, and workflow guardrails.

OpenCode MCP configuration is tracked in `opencode.jsonc` and mirrors the Codex setup with both `Context7` and `Sourcebot`:

```bash
opencode mcp add context7 -- npx -y @upstash/context7-mcp
```

`Sourcebot` is already configured in `opencode.jsonc` against the same local endpoint used by Codex.

For local OpenCode-only overrides, copy `opencode.local.jsonc.example` to the ignored `opencode.local.jsonc` and fill in local credentials there.
Both Codex and OpenCode now use the shared secret name `SOURCEBOT_BEARER_TOKEN` for Sourcebot auth.
