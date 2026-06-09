# Clawdia — Operations & COO

**Command:** `/clawdia-assistant` | **Color:** cyan | **Model:** Sonnet

Clawdia is the operational hub and default agent — a COO-like partner that manages your calendar, emails, tasks, meetings, prioritization, decisions, research, and documentation. She keeps the daily rhythm running smoothly, from morning briefings to end-of-day consolidation, ensuring nothing falls through the cracks.

## When to Use

- Starting your day with "good morning" to get a full briefing of agenda, emails, and tasks
- Checking, triaging, or drafting email replies
- Reviewing and organizing Todoist tasks
- Summarizing a meeting recording or syncing Fathom transcripts
- Making a decision between options and need a structured analysis

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `prod-good-morning` | Morning briefing with agenda, emails, tasks, and priorities |
| `prod-end-of-day` | Daily consolidation of work done, decisions, and learnings |
| `prod-dashboard` | 360 dashboard combining all routine outputs |
| `prod-memory-management` | Persistent memory system for context across sessions |
| `prod-review-todoist` | Task organization and cleanup in Todoist |
| `prod-trends` | Weekly trends analysis across metrics |
| `gog-calendar` | Google Calendar — view agenda and schedule events |
| `gog-email-triage` | Inbox prioritization and classification |
| `gog-email-draft` | Draft email replies with tone matching |
| `gog-email-send` | Send emails with explicit confirmation |
| `gog-followups` | Track stale email threads and generate reminders |
| `gog-tasks` | Task management and prioritization |
| `int-fathom` | Meeting recordings, transcripts, and summaries |
| `int-sync-meetings` | Sync Fathom recordings to local workspace |
| `int-todoist` | Todoist integration for task management |
| `schedule-task` | One-off task scheduling |
| `ops-process-doc` | Business process documentation with flowcharts and RACI |
| `ops-runbook` | Operational runbooks for recurring procedures |
| `ops-process-optimization` | Workflow improvement and bottleneck analysis |

## Example Interactions

```
/clawdia-assistant good morning
/clawdia-assistant check my emails and flag anything urgent
/clawdia-assistant summarize yesterday's meeting with HostGator
```

## Routines

| Time | Routine |
|------|---------|
| Daily 06:50 | Review Todoist |
| Daily 07:00 | Good Morning |
| Daily 07:15 | Email Triage |
| Every 30min | Sync Meetings |
| Daily 21:00 | End of Day |
| Daily 21:15 | Memory Sync |
| Daily 21:30 | Consolidated Dashboard |
| Friday 08:00 | Weekly Review |
| Sunday 09:00 | Memory Lint |

## Memory

Persistent memory at `.claude/agent-memory/clawdia/`. Stores project context, user preferences, and feedback across sessions.
