# Settings

The **Settings** page (`/settings`) is the central configuration hub for the workspace. It replaces the old Config page and provides three tabs.

## Workspace Tab

Edit `config/workspace.yaml` fields:

| Field | Description |
|-------|-------------|
| **Workspace Name** | Display name for the workspace |
| **Owner** | Workspace owner name |
| **Company** | Organization name |
| **Language** | Response language (20 options: pt-BR, en-US, es, fr, de, ja, ko, zh-CN, etc.) |
| **Timezone** | Schedule timezone (default: America/Sao_Paulo) |
| **Dashboard Port** | Port the dashboard runs on (default: 8080) |

Changes are saved to `config/workspace.yaml` using a read-merge-write pattern that preserves other YAML keys.

## Routines Tab

View and manage all scheduled routines from `config/routines.yaml`:

- **Toggle switch** — enable or disable any routine instantly
- **Inline schedule edit** — click the time/interval to edit, press Enter to save
- **Grouped by frequency** — daily, weekly, and monthly sections
- **Agent badges** — shows which agent owns each routine
- **Reload scheduler** — after changes, the scheduler is automatically notified via sentinel file

Routine creation and deletion are handled by agents (via the `create-routine` skill), not the UI.

## Reference Tab

Read-only views of key configuration files:

- **CLAUDE.md** — rendered workspace instructions
- **Makefile targets** — all available `make` commands
- **Commands** — slash command definitions

## API Reference

### Workspace
```
GET  /api/settings/workspace                          — read workspace config
PUT  /api/settings/workspace                          — update workspace fields
```

### Routines
```
GET  /api/settings/routines                           — list all routines
PATCH /api/settings/routines/<freq>/<slug>/toggle     — toggle enabled
PUT  /api/settings/routines/<freq>/<slug>             — update schedule/fields
```

### Scheduler
```
POST /api/settings/scheduler/reload                   — trigger scheduler reload
```

All write endpoints require authentication and `config:manage` permission.
