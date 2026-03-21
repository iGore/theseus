# AGENTS.md

## Purpose

This repository is a starter for a subagent-based migration workflow.
It defines a small, spec-first pipeline and the supporting assets needed to run it in Codex-style environments.

## Repository structure

- `.codex/agents/`: TOML agent profiles
- `.agents/skills/`: Agent Skills directories with `SKILL.md`
- `docker-compose.yml`: Sourcebot bootstrap
- `config.json`: Sourcebot example configuration
- `README.md`: quickstart and MCP setup notes

## Canonical role order

Use the following flow unless there is a strong reason to deviate:

1. `Analyzer`
2. `Spec-Writer`
3. `Verifier`
4. `Architect`
5. `Builder`

The `Orchestrator` coordinates handoffs, artifacts, and retry decisions across the full flow.

## Role intent

- `Analyzer`: inspect source code, `README.md`, dependency manifests, and relevant project metadata
- `Spec-Writer`: turn analyzed findings into a specification-first artifact
- `Verifier`: check traceability, consistency, and acceptance readiness
- `Architect`: plan target architecture and package selection for the builder
- `Builder`: implement from verified specification and architecture outputs

## MCP usage

- `Sourcebot`: use for repository structure, symbol lookup, references, and code-intelligence tasks
- `Context7`: use for framework, package, library, and best-practice documentation

Do not model these MCPs as skills. They are external context/access layers.

## Skill usage

Skills live under `.agents/skills/<skill-name>/SKILL.md`.

Use only the skills that are actually reflected in the report and starter flow:

- `reverse-engineering`
- `code-mapping`
- `dependency-inspection`
- `spec-driven-design`
- `specification-drafting`
- `contract-writing`
- `edge-case-structuring`
- `traceability-check`
- `consistency-review`
- `acceptance-validation`
- `evaluation-rubric`
- `ecosystem-mapping`
- `package-selection`
- `spec-to-code-translation`
- `test-scaffolding`
- `build-verification`

## Working rules

- Keep the workflow spec-first: do not jump from analysis straight to implementation.
- Keep the migration slice bounded and explicit.
- Record architecture and package decisions before build work starts.
- Prefer documented evidence over intuition.
- Keep agent names aligned with the report: `Analyzer`, `Spec-Writer`, `Verifier`, `Architect`, `Builder`, `Orchestrator`.

## Sourcebot bootstrap

To fetch the Sourcebot compose file:

```bash
curl -o docker-compose.yml https://raw.githubusercontent.com/sourcebot-dev/sourcebot/main/docker-compose.yml
```

## Context7 bootstrap

To register Context7 for Codex:

```bash
codex mcp add context7 -- npx -y @upstash/context7-mcp
```
