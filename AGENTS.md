# AGENTS.md

## Purpose

This repository is a starter for a subagent-based migration workflow.
It defines a spec-first pipeline implemented in OpenCode.

## Repository structure

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

## Phase model

- Phase 1 `Rekonstruktion`: `Analyzer -> Spec-Writer`
- Phase 2 `Transformationsplanung`: `Architect -> Planner`
- Phase 3 `Neuimplementierung`: `Task Decomposer -> Builder`

The Orchestrator spans all three phases, coordinates handoffs, and enforces the transition gates between them.

## Role intent

- `Analyzer`: reconstruct the bounded source module from repository structure, symbols, dependencies, build metadata, and supporting docs, and emit a complete functionality inventory with traceable checklist items
- `Spec-Writer`: condense one Analyzer functionality slice into one dedicated textual specification artifact in Markdown
- `Architect`: derive a Go target architecture, framework/library choices, and implementation structure from the specification using Context7-backed best practices
- `Planner`: translate specification and architecture decisions into an ordered Go reimplementation plan with work packages
- `Task Decomposer`: break the migration plan into small, implementation-ready Go coding tasks for the builder
- `Builder`: reimplement target code in Go from the specification and architecture decisions, then record technical verification evidence

## Artifact ownership

Each agent owns exactly the artifacts it produces. No agent may write another agent's artifacts.

| Agent | Produces | Location |
|---|---|---|
| Analyzer | `ANALYSIS.md`, `USE-CASES.md` | `target/specs/` |
| Spec-Writer | `SPEC.md` (one per ready item) | `target/specs/<slug>/` |
| Architect | `ARCHITECTURE.md` | `target/specs/` |
| Planner | `PLAN.md` (one per use-case) | `target/specs/<slug>/` |
| Task Decomposer | `TASKS.md` (one per use-case) | `target/specs/<slug>/` |
| Builder | Go source files + `BUILD.md` | `target/` + `target/specs/<slug>/` |
| Orchestrator | `SPECS.md` (coordination only) | `target/specs/` |
| Verifier | `VERIFIER REPORT` (ephemeral, not persisted) | — |

## MCP usage

- `Sourcebot`: primary MCP for reconstruction work — repository structure, symbol lookup, references, and code-intelligence tasks
- `Context7`: primary MCP for target-stack planning and implementation — framework, package, library, and best-practice documentation

Do not model these MCPs as skills. They are external context/access layers.

Keep secrets such as Sourcebot API keys out of tracked repo files. Use `opencode.local.jsonc` for local-only overrides. Use the env key `SOURCEBOT_BEARER_TOKEN` for Sourcebot auth.

## Working rules

- Keep the workflow spec-first: do not jump from analysis straight to implementation.
- Keep the migration slice bounded and explicit.
- Keep the Orchestrator focused on orchestration only: it may route work, enforce gates, and maintain coordination artifacts such as `SPECS.md`, but it must not draft stage artifacts on behalf of specialist roles.
- Have workflow artifacts live under `target/specs/` unless the user explicitly requests a different location.
- Store `SPEC.md`, `PLAN.md`, `TASKS.md`, and `BUILD.md` under `target/specs/<use-case-slug>/` and store shared `ARCHITECTURE.md` under `target/specs/`.
- Treat Go as the default target reimplementation language unless the user explicitly requests another language.
- Use Context7 wherever current framework, package, library, or best-practice documentation is needed for planning or implementation.
- When `USE-CASES.md` defines multiple functionality items, let the Orchestrator fan out one `Spec-Writer` per functionality item in the background.
- After shared architecture is available, let the Orchestrator fan out one `Planner` per use-case folder containing `SPEC.md` in the background.
- After use-case-scoped plans are available, let the Orchestrator fan out one `Task Decomposer` per use-case-scoped plan in the background.
- Treat `PLAN.md` as a living checklist artifact: preserve traceability and do not silently delete unfinished work.
- Treat `TASKS.md` as the execution checklist; Builder should update checkboxes as implementation completes.
- Give each agent only the context and MCP access needed for its current phase.
- Prefer documented evidence over intuition.

## Parity verification: in-flow vs. manual

For behavior-replacing rewrites, parity is verified in two distinct places —
keep them separate:

- **In-flow (automated, test-driven):** the Analyzer captures example/golden
  outputs once from the pinned source version and commits them as fixture
  files. The Builder builds test-driven against them and commits byte-level
  golden tests across the recorded mode × flag matrix and edge set; the
  Verifier confirms those tests exist, are byte-level, cover the matrix, and
  pass (`go test ./...`). **No agent installs or runs the live source system.**
- **Manual (out-of-flow, done by the maintainers at the end):** a live
  differential run of the actually-installed source system vs. the target
  binary across the full matrix. This is a deliberate human step performed
  once before release. It is **not** a phase gate and must never be wired into
  the agent flow. The committed golden outputs exist precisely to make this
  final manual run cheap and likely to pass on the first try.
- Intentional, declared differences live in a `## Conformance Matrix`
  (`identical` / `compatible` / `out-of-scope`); claiming "same results" while
  undeclared `compatible`/`out-of-scope` rows exist is a defect.

## Sourcebot bootstrap

```bash
curl -o docker-compose.yml https://raw.githubusercontent.com/sourcebot-dev/sourcebot/main/docker-compose.yml
docker compose up -d
```

## Context7 bootstrap

```bash
opencode mcp add context7 -- npx -y @upstash/context7-mcp
```
