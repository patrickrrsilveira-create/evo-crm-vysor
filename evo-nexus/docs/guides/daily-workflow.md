# Daily Workflow

A practical guide to running your business day with EvoNexus — from morning briefing to end-of-day wrap-up.

## Morning (7:00 - 8:00)

### 1. Good Morning Briefing

```
good morning
```

Clawdia generates your daily briefing:

- Today's calendar events
- Unread emails (prioritized)
- Open tasks from Todoist
- Yesterday's unfinished items
- Key metrics snapshot

### 2. Email Triage

```
triage my inbox
```

Clawdia reads unread emails, classifies by urgency, and proposes actions:

- **Reply now** — urgent items needing immediate response
- **Schedule** — items that need a response but not immediately
- **Archive** — informational emails, no action needed
- **Create task** — emails that represent work to be done

### 3. Task Review

```
review todoist
```

Clawdia organizes your tasks: uncategorized items get sorted, priorities are set, and you get a clear list of what to tackle today.

## During the Day (9:00 - 17:00)

### Project Updates

```
/atlas check linear
/atlas check github
```

Atlas shows sprint progress, issues in review, blockers, open PRs, and anything that needs your attention.

### Community Monitoring

```
/pulse community pulse
```

Pulse analyzes Discord and WhatsApp activity: trending topics, unanswered questions, sentiment shifts, and contributor highlights.

### Content Creation

```
/pixel write a post about [topic]
```

Pixel drafts platform-specific content following your brand voice. Supports LinkedIn, Twitter/X, Threads, and Bluesky.

### Meetings

Fathom syncs automatically every 30 minutes. After any meeting:

```
sync meetings
```

Get transcripts, summaries, and action items from your latest calls.

### Ad-hoc Questions

Just ask naturally:

- "How's the MRR trending?" (Flux answers)
- "Any blockers on EVO-589?" (Atlas checks Linear)
- "What's the community sentiment today?" (Pulse analyzes)
- "Draft a reply to that email from Marcelo" (Clawdia helps)

## Evening (18:00 - 21:00)

### Social Analytics (18:00)

```
social analytics
```

Pixel generates a cross-platform report: YouTube, Instagram, LinkedIn engagement and growth.

### Licensing Report (18:30)

```
licensing daily
```

Atlas reports on open source telemetry: active instances, geographic distribution, version adoption.

### Financial Pulse (19:00)

```
financial pulse
```

Flux queries Stripe and Omie to produce a financial health snapshot: MRR, charges, churn, receivables, payables.

### Community Pulse (20:00)

```
community pulse
```

Pulse generates the full daily community report with Discord and WhatsApp analysis.

### End of Day (21:00)

```
end of day
```

Clawdia wraps up: summarizes what was accomplished, logs decisions, notes pending items, and prepares context for tomorrow.

### Dashboard (21:30)

```
dashboard
```

Clawdia consolidates all daily outputs into a unified 360-degree dashboard.

## Weekly Cadence

| Day | Routine | Agent |
|---|---|---|
| Monday | Linear review, GitHub review | Atlas |
| Wednesday | Linear review, GitHub review | Atlas |
| Friday | Financial weekly | Flux |
| Friday | Licensing weekly | Atlas |
| Friday | Weekly review | Clawdia |
| Friday | Social analytics weekly | Pixel |
| Friday | Strategy digest | Sage |
| Sunday | Health check-in | Kai |

### Friday Deep Dives

Friday is the weekly synthesis day:

```
weekly review
```

Clawdia generates a comprehensive week-in-review. Follow up with:

```
strategy digest
```

Sage consolidates everything into strategic insights and recommendations.

## Monthly Close (1st of each month)

```
monthly close
```

Flux kicks off the month-end closing process:

- Simplified P&L
- Pending reconciliations
- Receivables and payables summary
- Action items for the finance team

Also runs:

- Community monthly report (Pulse)
- Licensing monthly report (Atlas)
- Social analytics monthly (Pixel)

## Automation

All these routines can run automatically via the scheduler:

```bash
make scheduler
```

See [Creating Routines](creating-routines.md) for details on the schedule and customization.
