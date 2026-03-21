# Prompt Templates

## 1. Generate specification

Use when analysis artifacts already exist.

```text
Create a technical specification from the provided analysis artifacts.

Requirements:
- Define scope and non-goals.
- Describe inputs, outputs, and behavioral rules.
- Include data contracts, error handling, and edge cases.
- Add explicit acceptance criteria.
- Keep all important claims traceable to the source analysis.

Output format:
- spec.md structure with clear headings.
```

## 2. Generate implementation plan

Use after the specification is stable enough for design decisions.

```text
Create an implementation plan from the approved specification.

Requirements:
- Derive implementation phases from the specification.
- Identify sequencing, dependencies, and validation strategy.
- Keep the plan inside the approved scope.
- Explain how the plan satisfies the specification.

Output format:
- plan.md with ordered sections and rationale.
```

## 3. Generate task breakdown

Use after the implementation plan exists.

```text
Break the implementation plan into executable tasks.

Requirements:
- Keep tasks bounded and independently verifiable.
- Link each task back to relevant parts of the specification or plan.
- Include implementation, testing, and verification work where needed.

Output format:
- tasks.md as an ordered task list with traceability notes.
```

## 4. Refine spec after feedback

Use when verification or implementation discovers a mismatch.

```text
Refine the specification using the feedback provided.

Requirements:
- Update only the affected parts of the specification.
- Preserve traceability to the original analysis.
- Make the change explicit instead of silently rewriting intent.
- Note the downstream impact on plan or tasks if relevant.
```
