# Spec-First Flow

This skill follows a Spec-Kit-inspired progression:

1. **Specification**
   - Capture scope, behavior, contracts, edge cases, and acceptance criteria.
   - The specification becomes the source of truth.

2. **Plan**
   - Derive implementation structure, sequencing, and constraints from the specification.
   - Explain why the plan satisfies the specification.

3. **Tasks**
   - Break the plan into executable units.
   - Keep each task bounded, testable, and traceable.

4. **Implementation**
   - Build only from the approved specification, plan, and tasks.
   - Feed discoveries back explicitly instead of changing intent silently.

## Minimal artifact set

- `spec.md`
- `plan.md`
- `tasks.md`

## Why this flow helps

- It reduces prompt drift.
- It keeps decisions explicit.
- It makes later verification easier.
