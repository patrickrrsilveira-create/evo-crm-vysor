# Aria — HR & People

**Command:** `/aria-hr` | **Color:** pink | **Model:** Sonnet

Aria is the HR and People Operations agent — it manages the recruiting pipeline, performance reviews, onboarding, compensation benchmarking, org planning, and policy lookup. Aria helps you build and maintain a healthy organization by keeping people processes structured, documented, and on schedule.

## When to Use

- Check the current recruiting pipeline status and candidate progress
- Generate an onboarding checklist for a new hire starting next week
- Kick off the quarterly performance review cycle
- Run a compensation benchmark for a specific role or level
- Look up the remote work policy or any other company policy

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `hr-recruiting-pipeline` | Track and manage hiring pipeline stages and candidate status |
| `hr-performance-review` | Set up review cycles with self-assessment and manager templates |
| `hr-onboarding` | Generate first-week checklists and onboarding plans |
| `hr-comp-analysis` | Benchmark salaries and model compensation packages |
| `hr-draft-offer` | Draft offer letters with compensation details and terms |
| `hr-interview-prep` | Create structured interview plans with competency-based questions |
| `hr-org-planning` | Plan headcount, org design, and team structure |
| `hr-people-report` | Generate headcount, attrition, and organizational health reports |
| `hr-policy-lookup` | Find and explain company policies in plain language |

## Example Interactions

```
/aria-hr show me the recruiting pipeline status
/aria-hr create an onboarding checklist for a new backend developer
/aria-hr run a compensation benchmark for senior engineers in Brazil
```

## Routines

This agent is not used in any scheduled routines.

## Memory

Persistent memory at `.claude/agent-memory/aria/`. Stores project context, user preferences, and feedback across sessions.
