# Atlas — Projects & Ops

**Command:** `/atlas-project` | **Color:** green | **Model:** Sonnet

Atlas is the project management agent — tracking Linear issues, GitHub pull requests, sprint progress, licensing telemetry, capacity planning, and status reports. Atlas keeps all development workstreams visible and helps identify blockers before they become problems.

## When to Use

- Checking project status, sprint progress, or blockers
- Reviewing open GitHub PRs and community issues
- Monitoring open source licensing growth and telemetry
- Planning team capacity and workload distribution
- Generating status reports with KPIs and risk flags

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `int-github-review` | Repo health — open PRs, issues, stars, releases |
| `int-linear-review` | Sprint and issue tracking in Linear |
| `int-licensing` | Open source telemetry — instances, geo, versions |
| `prod-licensing-daily` | Daily growth snapshot from licensing data |
| `prod-licensing-weekly` | Weekly trends in licensing telemetry |
| `prod-licensing-monthly` | Monthly deep analysis of licensing growth |
| `ops-capacity-plan` | Team workload analysis and utilization forecast |
| `ops-status-report` | KPI reports with risk flags and action items |
| `ops-change-request` | Change management with impact analysis and rollback |

## Example Interactions

```
/atlas-project what's the sprint status for this week?
/atlas-project check open PRs on evolution-api
/atlas-project how's licensing growth looking this month?
```

## Routines

| Time | Routine |
|------|---------|
| Mon/Wed/Fri 09:00 | Linear Review |
| Mon/Wed/Fri 09:15 | GitHub Review |
| Daily 18:30 | Licensing Daily |
| Friday 07:45 | Licensing Weekly |
| 1st of month | Licensing Monthly |

## Memory

Persistent memory at `.claude/agent-memory/atlas/`. Stores project context, user preferences, and feedback across sessions.
