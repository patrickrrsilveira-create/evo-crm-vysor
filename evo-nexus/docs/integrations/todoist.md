# Todoist Integration

Todoist is used for task management across the workspace. The @clawdia agent reviews, organizes, and prioritizes tasks daily.

## Setup

### 1. Get Your API Token

1. Go to [app.todoist.com](https://app.todoist.com)
2. Navigate to **Settings > Integrations > Developer**
3. Copy your **API Token**

### 2. Install the CLI

```bash
npm install -g todoist-ts-cli@^0.2.0
```

### 3. Configure

```bash
todoist auth <your-token>
# or set in .env:
# TODOIST_API_TOKEN=your_token_here
```

### 4. Test the Connection

```bash
make review
```

This runs the Todoist review routine, which lists and organizes tasks in the Evolution project.

## Available Commands

The `int-todoist` skill uses the `todoist` CLI:

| Command | What it does |
|---|---|
| `todoist` / `todoist today` | Show today's tasks |
| `todoist tasks --all` | List all tasks |
| `todoist tasks -p "Project"` | Tasks in a specific project |
| `todoist tasks -f "p1"` | Filter by priority |
| `todoist add "Task" --due DATE` | Create a task with due date |
| `todoist add "Task" --priority 1` | Create a high-priority task |
| `todoist done ID` | Complete a task |
| `todoist reopen ID` | Reopen a completed task |
| `todoist update ID --due DATE` | Reschedule a task |
| `todoist move ID -p "Project"` | Move task to another project |
| `todoist search "query"` | Search tasks |
| `todoist projects` | List projects |

## Skills That Use Todoist

| Skill | What it does |
|---|---|
| `int-todoist` | Direct task management -- create, update, complete, search |
| `prod-review-todoist` | Review and organize tasks: categorize, translate to PT-BR, clean up |
| `prod-good-morning` | Morning briefing includes today's tasks |
| `prod-end-of-day` | End-of-day review checks task completion |

## Automated Routines

| Routine | Schedule | Make command |
|---|---|---|
| Review Todoist | 06:50 BRT daily | `make review` |
| Good Morning | 07:00 BRT daily | `make morning` |
| End of Day | 21:00 BRT daily | `make eod` |
