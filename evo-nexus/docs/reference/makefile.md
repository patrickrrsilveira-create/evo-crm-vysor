# Makefile Reference

All available `make` commands. Run `make help` to see this list in your terminal.

## Setup

```bash
make setup           # Interactive setup wizard (runs setup.py)
                     # Checks prerequisites (Claude Code, uv, Node.js)
                     # Installs Python deps via `uv sync`
                     # Creates config files, .env, CLAUDE.md, workspace folders
```

## Dashboard

```bash
make dashboard-app   # Start the web dashboard (React + Flask) on localhost:8080
                     # Builds frontend, then starts the Flask backend
```

## Core Routines

These ship with the repo and run the essential daily loop:

```bash
make morning         # Morning briefing -- agenda, emails, tasks (@clawdia)
make eod             # End of day -- memory, logs, learnings (@clawdia)
make backup-daily    # Daily backup -- export gitignored user data as ZIP (systematic)
make memory          # Memory sync -- consolidate agent memory (@clawdia)
make memory-lint     # Memory health check -- contradictions, gaps, stale data (@clawdia)
make weekly          # Full weekly review (@clawdia)
```

## Dynamic Routine Runner

Custom routines (user-specific, in `ADWs/routines/custom/`) are run via the dynamic runner. No hardcoded Makefile targets needed -- routines are discovered automatically from script files.

```bash
make run R=fin-pulse        # Financial pulse (@flux)
make run R=community        # Daily community pulse (@pulse)
make run R=licensing-month  # Monthly licensing report (@atlas)
make run R=social           # Social analytics (@pixel)
make run R=strategy         # Strategy digest (@sage)

make list-routines          # List all available routines (core + custom)
```

Any script in `ADWs/routines/` or `ADWs/routines/custom/` is automatically discoverable. The routine ID is derived from the script name (e.g., `financial_pulse.py` becomes `fin-pulse`).

## Combos

```bash
make daily           # Runs: sync + review (sync meetings then organize tasks)
```

## Servers & Channels

```bash
make scheduler              # Start the routine scheduler (runs in foreground)

# Telegram channel
make telegram               # Start Telegram bot in background (screen session)
make telegram-stop          # Stop the Telegram bot
make telegram-attach        # Attach to Telegram terminal (Ctrl+A D to detach)

# Discord channel (bidirectional chat, separate from the @pulse API integration)
make discord-channel        # Start Discord channel in background (screen session)
make discord-channel-stop   # Stop the Discord channel
make discord-channel-attach # Attach to Discord channel terminal (Ctrl+A D to detach)

# iMessage channel (macOS only)
make imessage               # Start iMessage channel in background (screen session)
make imessage-stop          # Stop the iMessage channel
make imessage-attach        # Attach to iMessage terminal (Ctrl+A D to detach)
```

See [Channels Guide](../guides/channels.md) for full setup instructions.

## Backup & Restore

```bash
make backup          # Export all gitignored user data to local ZIP (backups/)
make backup-s3       # Export to local ZIP + upload to S3
make restore FILE=<path>              # Restore from ZIP (default mode: merge)
make restore FILE=<path> MODE=replace # Restore overwriting existing files
make backup-list     # List local backups
make backup-list TARGET=s3            # List S3 backups
```

Backup captures: memory, agent-memory, workspace data, config, dashboard DB, logs, custom agents/commands/templates/routines, `.env`. Excludes reconstructible dirs (`node_modules/`, `.venv/`, `dist/`).

Restore modes:
- **merge** (default): only restore files that don't exist locally
- **replace**: overwrite everything with backup versions

S3 requires `BACKUP_S3_BUCKET` in `.env` and `boto3` installed (`uv add boto3`).

## Observability

```bash
make logs            # Show latest JSONL log entries
make logs-detail     # List detailed log files
make logs-tail       # Show the latest full detailed log
make metrics         # Per-routine metrics table (runs, cost, tokens, success rate)
make clean-logs      # Remove logs older than 30 days
```

## Docker (VPS Deployment)

```bash
make docker-dashboard  # Start dashboard in Docker (port 8080)
make docker-telegram   # Start Telegram bot in Docker
make docker-down       # Stop all containers
make docker-logs       # Show container logs (follow mode)
make docker-run ADW=good_morning.py  # Run a specific routine in Docker
make docker-build      # Build the Docker image
```

## Agent Teams (Experimental)

Parallel multi-agent versions of consolidation routines. Higher token cost (~3-5x), faster execution.

```bash
make team-strategy   # Strategy digest via agent team (@sage leads, @atlas/@flux/@pulse/@pixel)
make team-dashboard  # Dashboard via agent team (@clawdia leads, @atlas/@flux/@pulse/@dex)
make team-weekly     # Weekly review via agent team (@clawdia leads, @atlas/@flux/@pulse)
```

## Help

```bash
make help            # Show all available commands with descriptions
```
