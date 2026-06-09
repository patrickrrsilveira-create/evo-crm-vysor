# Kai — Personal Assistant

**Command:** `/kai-personal-assistant` | **Color:** blue | **Model:** Sonnet

Kai is your personal assistant — completely isolated from professional matters. It handles health tracking, habits, personal routines, appointments, and travel planning. Kai reads health data from `workspace/personal/data/health-data.js` and keeps your personal life organized without mixing it into work context.

## When to Use

- Check your health progress or review recent health metrics
- Schedule a blood test, doctor visit, or other medical appointment
- Plan a personal trip or vacation itinerary
- Review personal appointments for the week
- Track habits and daily routines

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `gog-calendar` | Manage personal appointments and schedule events |
| `gog-tasks` | Create and organize personal to-do items |

## Example Interactions

```
/kai-personal-assistant how is my health progress this month?
/kai-personal-assistant schedule a blood test for next Tuesday morning
/kai-personal-assistant plan a 5-day trip to Lisbon in July
```

## Routines

| Day | Time | Routine | Make |
|-----|------|---------|------|
| Sunday | 10:00 | Health Check-in | `make health` |

## Memory

Persistent memory at `.claude/agent-memory/kai/`. Stores project context, user preferences, and feedback across sessions.
