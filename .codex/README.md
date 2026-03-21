# Theseus Codex Starter

Minimal starter for an agentic migration setup for the project PoC.

## Structure

- `.codex/agents/`: agent profiles as TOML (`*.toml`)
- `.agents/skills/`: reusable cross-agent Agent Skills directories

## Included Roles

- `orchestrator.toml`: coordinates stage order and artifact flow
- `analyzer.toml`: reverse engineering and structural analysis
- `architect.toml`: target ecosystem mapping and package selection
- `spec-writer.toml`: spec creation for a spec-first workflow
- `verifier.toml`: consistency, traceability, and release-gate checks
- `builder.toml`: target-code implementation from the validated specification

## Included Skills

- `.agents/skills/`
- `reverse-engineering/`
- `code-mapping/`
- `dependency-inspection/`
- `ecosystem-mapping/`
- `package-selection/`
- `spec-driven-design/`
- `specification-drafting/`
- `contract-writing/`
- `edge-case-structuring/`
- `traceability-check/`
- `consistency-review/`
- `acceptance-validation/`
- `spec-to-code-translation/`
- `test-scaffolding/`
- `build-verification/`
- `evaluation-rubric/`

## Quick Start

1. Describe the source system and target system in one task.
2. Choose a bounded demonstration module for the PoC.
3. Start the orchestrator and keep the order Analyzer -> Spec-Writer -> Verifier -> Architect -> Builder.
4. Persist each stage result as an artifact (Markdown or JSON).

## Configured MCPs

- `Context7`: documentation and library reference context for authoring, orchestration, and implementation support
- `Sourcebot`: code-intelligence and repository context via the local bootstrap setup

Add `Context7` with:

```bash
codex mcp add context7 -- npx -y @upstash/context7-mcp
```

## Sourcebot Bootstrap

For the Sourcebot-based code-intelligence setup, fetch the official Docker Compose file with:

```bash
curl -o docker-compose.yml https://raw.githubusercontent.com/sourcebot-dev/sourcebot/main/docker-compose.yml
```
