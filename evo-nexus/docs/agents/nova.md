# Nova — Product Management

**Command:** `/nova-product` | **Color:** blue | **Model:** Sonnet

Nova is the Product Management agent, covering feature specs and PRDs, product metrics analysis, roadmap management, brainstorming sessions, stakeholder updates, and user research synthesis. She turns ideas into structured documents, keeps the roadmap current, and ensures stakeholders stay informed with data-driven updates.

## When to Use

- You need to write a PRD or feature spec from a problem statement or idea
- You want to review sprint or product metrics and spot trends
- The roadmap needs updating after new priorities or market changes
- You want to brainstorm ideas for a new feature or product direction
- You need to prepare a stakeholder update on quarterly progress

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `pm-write-spec` | Creates feature specs and PRDs from problem statements |
| `pm-metrics-review` | Analyzes product metrics with trend detection and insights |
| `pm-roadmap-update` | Manages roadmap priorities, additions, and trade-offs |
| `pm-product-brainstorming` | Facilitates structured ideation and problem exploration |
| `pm-stakeholder-update` | Generates audience-tailored status updates |
| `pm-synthesize-research` | Synthesizes interviews, surveys, and feedback into insights |
| `int-linear-review` | Tracks sprint progress, issues, and blockers in Linear |
| `mkt-competitive-brief` | Competitive positioning and messaging comparison |
| `evo-sprint-planning` | Sprint planning and task generation from epics |

## Example Interactions

```
/nova-product write a PRD for the new agent builder feature
/nova-product review sprint metrics and highlight risks
/nova-product brainstorm ideas for improving the onboarding flow
```

## Routines

This agent is not used in any scheduled routines.

## Memory

Persistent memory at `.claude/agent-memory/nova/`. Stores project context, user preferences, and feedback across sessions.
