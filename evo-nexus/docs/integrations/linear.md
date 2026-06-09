# Linear Integration

Linear is used for issue tracking, sprint management, and project planning. The @atlas agent monitors issues, blockers, sprint progress, and review queues via the Linear MCP server.

## Setup

### 1. Configure MCP

Linear connects through an MCP server (`linear-server`) configured in `.claude/settings.json`. Claude Code starts it automatically when Linear tools are needed.

No `.env` variables are required -- the MCP server handles authentication through its own configuration.

### 2. Test the Connection

```bash
make linear
```

This runs the Linear review routine, which checks issues in review, blockers, stale items, and sprint progress.

## Available Tools

The Linear MCP server exposes these tools:

| Tool | What it does |
|---|---|
| `list_issues` | List issues with filters (state, assignee, priority) |
| `get_issue` | Get details of a specific issue |
| `save_issue` | Create or update an issue |
| `list_projects` | List all projects |
| `get_project` | Get project details |
| `list_cycles` | List sprints/cycles (current, past, future) |
| `list_teams` | List teams |
| `list_users` | List team members |
| `list_comments` | List comments on an issue |
| `save_comment` | Add a comment to an issue |
| `list_milestones` | List milestones |
| `research` | Search across Linear data |

## Key Metrics

The `int-linear-review` skill analyzes:

- **In Review:** Issues awaiting review, who needs to review, days pending
- **Blockers:** Urgent and high-priority issues with responsible assignees
- **Stale issues:** In Progress items with no updates in 3+ days
- **My issues:** User's assigned issues sorted by priority
- **Sprint progress:** Completion percentage and deadline

## Skills That Use Linear

| Skill | What it does |
|---|---|
| `int-linear-review` | Review issues, blockers, stale items, sprint progress |
| `evo-sprint-planning` | Generate sprint plans from epics |
| `evo-sprint-status` | Sprint status summary with risks |

## Automated Routines

| Routine | Schedule | Make command |
|---|---|---|
| Linear Review | Mon/Wed/Fri 09:00 BRT | `make linear` |
