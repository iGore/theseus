# AGENTS.md

## Purpose

This repository is a starter for a subagent-based migration workflow.
It defines a small, spec-first pipeline and the supporting assets needed to run it in Codex-style and OpenCode environments.

## Repository structure

- `.codex/agents/`: TOML agent profiles
- `.opencode/agents/`: Markdown agent profiles with YAML frontmatter
- `opencode.jsonc`: OpenCode configuration for MCP servers and default agent selection
- `.agents/skills/`: Agent Skills directories with `SKILL.md`
- `docker-compose.yml`: Sourcebot bootstrap
- `config.json`: Sourcebot example configuration
- `README.md`: quickstart and MCP setup notes

## Canonical role order

Use the following flow unless there is a strong reason to deviate:

1. `Analyzer`
2. `Spec-Writer`
3. `Architect`
4. `Planner`
5. `Task Decomposer`
6. `Builder`

The `Orchestrator` coordinates handoffs, artifact gates, and retry decisions across the full flow. It must remain coordination-only and must not perform Analyzer, Spec-Writer, Architect, Planner, Task Decomposer, or Builder work itself.

This role set now matches the full conceptual model from the report: `Orchestrator`, `Analyzer`, `Spec-Writer`, `Architect`, `Planner`, `Task Decomposer`, and `Builder`.

## Phase model

- Phase 1 `Rekonstruktion`: `Analyzer -> Spec-Writer`
- Phase 2 `Transformationsplanung`: `Architect -> Planner`
- Phase 3 `Neuimplementierung`: `Task Decomposer -> Builder`

The Orchestrator spans all three phases, coordinates handoffs, and enforces the transition gates between them.

## Role intent

- `Analyzer`: reconstruct the bounded source module from repository structure, symbols, dependencies, build metadata, and supporting docs
- `Spec-Writer`: condense Analyzer findings into the central textual specification artifact in Markdown
- `Architect`: derive a Go target architecture, framework/library choices, and implementation structure from the specification using Context7-backed best practices
- `Planner`: translate specification and architecture decisions into an ordered Go reimplementation plan with work packages
- `Task Decomposer`: break the migration plan into small, implementation-ready Go coding tasks for the builder
- `Builder`: reimplement target code in Go from the specification and architecture decisions, then record technical verification evidence

## MCP usage

- `Sourcebot`: primary MCP for reconstruction work - repository structure, symbol lookup, references, and code-intelligence tasks
- `Context7`: primary MCP for target-stack planning and implementation - framework, package, library, and best-practice documentation

Do not model these MCPs as skills. They are external context/access layers.

Keep secrets such as Sourcebot API keys out of tracked repo files. Use local-only files such as `.codex/config.local.toml`, `.codex/secrets.toml`, or `opencode.local.jsonc`. A tracked template is available as `opencode.local.jsonc.example`. Use the shared env key `SOURCEBOT_BEARER_TOKEN` when wiring Sourcebot auth across Codex and OpenCode.

## Working rules

- Keep the workflow spec-first: do not jump from analysis straight to implementation.
- Keep the migration slice bounded and explicit.
- Keep the Orchestrator focused on orchestration only: it may route work, enforce gates, and maintain coordination artifacts such as `SPEC-INDEX.md`, but it must not draft stage artifacts on behalf of specialist roles.
- Have workflow artifacts live under `target/specs/` unless the user explicitly requests a different location; this includes Analyzer outputs such as `ANALYSIS.md` and `OVERVIEW.md`, while generated implementation code can still live under `target/` when needed.
- Store `SPEC.md`, `PLAN.md`, `TASKS.md`, and `BUILD.md` under `target/specs/<use-case-slug>/` and store shared `ARCHITECTURE` under `target/specs/` unless the user explicitly requests a different location.
- Treat Go as the default target reimplementation language unless the user explicitly requests another language.
- Use Context7 wherever current framework, package, library, or best-practice documentation is needed for planning or implementation.
- When `OVERVIEW.md` defines multiple requirements, tasks, or use cases, let the Orchestrator fan out one `Spec-Writer` per requirement in the background and store each requirement's artifacts under its own use-case subfolder in `target/specs/`.
- After shared architecture is available, let the Orchestrator fan out one `Planner` per use-case folder containing `SPEC.md` in the background.
- After use-case-scoped plans are available, let the Orchestrator fan out one `Task Decomposer` per use-case-scoped plan in the background.
- When use-case-specific planning or build artifacts are produced, store `PLAN.md`, `TASKS.md`, and `BUILD.md` in the same `target/specs/<use-case-slug>/` folder as the related `SPEC.md` file.
- Treat the textual specification as the main handoff artifact between reconstruction and implementation.
- Record architecture and package decisions before build work starts.
- Record a migration plan before detailed implementation work starts.
- Use task decomposition to keep the builder on bounded, implementation-ready work packages.
- Give each agent only the context and MCP access needed for its current phase.
- Prefer documented evidence over intuition.
- Keep agent names aligned with the report: `Analyzer`, `Spec-Writer`, `Architect`, `Planner`, `Task Decomposer`, `Builder`, `Orchestrator`.

## Sourcebot bootstrap

To fetch the Sourcebot compose file:

```bash
curl -o docker-compose.yml https://raw.githubusercontent.com/sourcebot-dev/sourcebot/main/docker-compose.yml
```

`Sourcebot` is configured for both Codex and OpenCode against `http://localhost:3000/api/mcp`.

## Context7 bootstrap

To register Context7 for Codex:

```bash
codex mcp add context7 -- npx -y @upstash/context7-mcp
```

To register Context7 for OpenCode:

```bash
opencode mcp add context7 -- npx -y @upstash/context7-mcp
```
