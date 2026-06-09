# Pulse — Community

**Command:** `/pulse-community` | **Color:** blue | **Model:** Sonnet

Pulse is the community monitoring agent — tracking Discord channels, WhatsApp groups, sentiment analysis, FAQ management, and engagement metrics. Pulse keeps a finger on the community's pulse, surfacing recurring questions, emerging issues, and contributor highlights so nothing important gets missed.

## When to Use

- Checking how the community is doing overall
- Identifying frequently asked questions or recurring support topics
- Analyzing sentiment trends across Discord and WhatsApp
- Investigating community issues or escalations
- Generating monthly engagement summaries or creating community updates

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `pulse-daily` | 24h community pulse across Discord and WhatsApp |
| `pulse-weekly` | Weekly engagement analysis and trends |
| `pulse-monthly` | Monthly community report with growth and sentiment |
| `pulse-faq-sync` | Update FAQ from recurring community questions |
| `discord-get-messages` | Read and search Discord channel messages |
| `discord-send-message` | Post messages to Discord channels |
| `discord-list-channels` | Audit server channel structure |
| `discord-create-channel` | Create new Discord channels |
| `discord-manage-channel` | Update channel settings and permissions |
| `int-whatsapp` | WhatsApp group messages and stats |

## Example Interactions

```
/pulse-community how's the community doing this week?
/pulse-community what are the most frequent questions in Discord?
/pulse-community generate the monthly engagement summary
```

## Routines

| Time | Routine |
|------|---------|
| Daily 20:00 | Community Pulse |
| Daily 20:15 | FAQ Sync |
| Monday 09:30 | Community Weekly |
| 1st of month | Community Monthly |

## Memory

Persistent memory at `.claude/agent-memory/pulse/`. Stores project context, user preferences, and feedback across sessions.
