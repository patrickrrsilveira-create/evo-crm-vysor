# Sage — Strategy

**Command:** `/sage-strategy` | **Color:** orange | **Model:** Sonnet

Sage is the senior strategist agent. It handles OKR tracking, competitive analysis, roadmap planning, decision frameworks, and weekly strategy digests that synthesize financial, product, and community data into actionable strategic insights. Use Sage when you need to think about the big picture, weigh trade-offs, or make informed decisions about company direction.

## When to Use

- You need to decide whether to prioritize initiative X over Y
- You want to review OKR progress or define new quarterly objectives
- You need a competitive landscape analysis or market positioning review
- You want a strategy digest that consolidates data across all domains
- You need a decision framework for a partnership, monetization approach, or negotiation

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `sage-strategy-digest` | Weekly strategic synthesis of financial, product, community, and market data |
| `sage-okr-review` | Track OKR progress, update key results, and flag at-risk objectives |
| `sage-competitive-analysis` | Analyze competitors and evaluate market positioning |

## Example Interactions

```
/sage-strategy should we prioritize the self-hosted plan or the managed cloud offering?
/sage-strategy OKR progress — how are we tracking for Q2?
/sage-strategy competitive analysis of the AI CRM space
```

## Routines

| Day | Routine | Command |
|-----|---------|---------|
| Friday 09:00 | Strategy Digest | `make strategy` |

## Memory

Persistent memory at `.claude/agent-memory/sage/`. Stores project context, user preferences, and feedback across sessions.
