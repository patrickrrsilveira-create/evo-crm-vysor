# Core Routines

Core routines are the backbone of the EvoNexus daily loop. They ship with the repo, are hardcoded in `scheduler.py`, and run automatically. Unlike custom routines (which live in `ADWs/routines/custom/` and are configured via `config/routines.yaml`), core routines require zero configuration.

There are 6 core routines. Together they form a closed loop: **orient in the morning, work during the day, consolidate at night, backup data, sync memory, and review weekly**.

| Routine | Schedule | Agent | Make |
|---------|----------|-------|------|
| Good Morning | Daily 07:00 | @clawdia | `make morning` |
| End of Day | Daily 21:00 | @clawdia | `make eod` |
| Memory Sync | Daily 21:15 | @clawdia | `make memory` |
| Memory Lint | Sunday 09:00 | @clawdia | `make memory-lint` |
| Weekly Review | Friday 08:00 | @clawdia | `make weekly` |
| Daily Backup | Daily 21:00 | systematic | `make backup-daily` |

---

## Daily Backup

**Why it matters:** Gitignored data (memory, configs, logs, customizations, dashboard DB) is lost on clean installs or machine changes. This routine exports all user data as a ZIP, enabling portability without breaking `git pull` updates. If S3 is configured, it also uploads the backup to the cloud.

**Execution:** Pure Python via `run_script()` (no AI, no tokens)
**Timeout:** 5 minutes

### What it does

1. Collects all gitignored files using `git ls-files --others --ignored --exclude-standard`
2. Filters out heavy reconstructible dirs (`node_modules/`, `.venv/`, `__pycache__/`, `dist/`)
3. Creates a ZIP with `manifest.json` (version, timestamp, hostname, file list)
4. Saves to `backups/evonexus-backup-YYYYMMDD-HHMMSS.zip`
5. If `BACKUP_S3_BUCKET` env var is set, uploads the ZIP to S3

### Output

A ZIP file in `backups/` containing all workspace user data. Restorable via `make restore FILE=<path>`.

### CLI

```bash
make backup           # Manual backup (local)
make backup-s3        # Manual backup (local + S3)
make restore FILE=<path> MODE=merge    # Restore (merge: skip existing)
make restore FILE=<path> MODE=replace  # Restore (replace: overwrite all)
make backup-list      # List available backups
make backup-daily     # Run the daily backup routine
```

---

## Good Morning

**Why it matters:** Without a structured start, the day drifts. This routine scans your calendar, emails, tasks, and recent work to give you a clear picture of where you left off and what needs attention today.

**Skill:** `prod-good-morning`
**Timeout:** 10 minutes

### What it does

1. Reads `CLAUDE.md` (master context) and the last 3 session logs from `workspace/daily-logs/`
2. Reads active project overviews
3. Gathers live data:
   - Today's calendar via Google Calendar MCP
   - Unread emails via Gmail MCP
   - Today's tasks via Todoist
4. Delivers a brief recap (pt-BR):
   - Recent work (2-4 bullets)
   - Open/mid-flight items
   - Today's meetings and times
   - Attention-needed emails
   - Priority tasks
5. Gives one clear recommendation on what to work on first
6. Asks what you want to tackle

### Output

A structured morning briefing displayed in the terminal. No file saved -- this is an interactive session starter.

---

## End of Day

**Why it matters:** Without consolidation, context is lost between sessions. This routine captures everything that happened -- decisions, learnings, people updates, pending items -- so tomorrow's Good Morning has fresh context to work with.

**Skill:** `prod-end-of-day`
**Timeout:** 10 minutes

### What it does

1. Collects data from the entire day:
   - Agent memory updates from `.claude/agent-memory/` (all agents)
   - ADW logs from `ADWs/logs/YYYY-MM-DD.jsonl`
   - Meeting summaries from `workspace/meetings/summaries/`
   - Tasks completed/created via Todoist
   - Git changes (`git diff --stat`, `git log --oneline --since="today 00:00"`)
   - Current session conversation
2. Consolidates into a daily log:
   - Decisions made and why
   - Patterns, feedback, and corrections
   - New context about people
   - Real pending items (not aspirational)

### Output

A daily log saved to `workspace/daily-logs/[C] YYYY-MM-DD-eod.md`.

---

## Memory Sync

**Why it matters:** Memory files drift out of date if not actively maintained. This routine reads recent daily logs, meetings, and git history to extract new facts and propagate them across the memory system, keeping it consistent and current.

**Execution:** Direct prompt via `run_claude()` (not a skill)
**Timeout:** 10 minutes

### What it does

1. Reads last 3 daily logs from `workspace/daily-logs/`
2. Reads meeting summaries from the last 3 days
3. Analyzes `git log --oneline --since="3 days ago"`
4. Extracts and classifies:
   - Decisions --> `project` type memories
   - People context --> `user` type memories
   - Feedback/corrections --> `feedback` type memories
   - New terms/references --> `reference` type memories
5. De-duplicates: updates existing memories instead of creating duplicates
6. **Ingest propagation:** When updating a memory, finds and updates other memories that reference the same entity
7. Updates `MEMORY.md` index and `memory/index.md`
8. Appends to `memory/log.md`: `[DATE] SYNC -- summary`

### Output

A report showing count of memories created/updated by type and cross-references propagated.

---

## Memory Lint

**Why it matters:** Over time, memory accumulates contradictions, orphan files, stale claims, and coverage gaps. This weekly health check catches problems before they compound, keeping the knowledge base reliable.

**Execution:** Direct prompt via `run_claude()` (not a skill)
**Timeout:** 10 minutes
**Schedule:** Every Sunday at 09:00

### What it does

1. Reads `memory/index.md` to get the full catalog
2. Verifies every listed file exists and is not empty
3. Detects orphan files (exist in `memory/` but not in the index)
4. Cross-reference checks:
   - People mentioned in projects should have a file in `memory/people/`
   - Projects mentioned in people files should have a file in `memory/projects/`
   - Glossary terms should point to valid entries
5. Stale data check:
   - Flags Linear issues that may be closed
   - Flags people roles/companies that may have changed
   - Flags any claim with a date older than 60 days
6. Contradiction check:
   - Compares `CLAUDE.md` hot cache with `memory/glossary.md`
   - Compares people info in `CLAUDE.md` vs `memory/people/*.md`
7. Coverage gaps:
   - Active projects in `CLAUDE.md` without a `memory/projects/` file
   - People in `CLAUDE.md` without a `memory/people/` file
8. Auto-fixes what it can:
   - Updates `memory/index.md` with orphan files
   - Fixes obvious contradictions (prefers the more recent source)
   - Adds missing cross-references
9. Appends to `memory/log.md`: `[DATE] LINT -- summary`

### Output

A report listing issues found, fixes applied automatically, and items needing manual review.

---

## Weekly Review

**Why it matters:** Daily logs capture details, but the weekly review steps back to see the big picture -- what actually got done, what's carrying over, and what's coming next week. It produces a formatted report you can reference or share.

**Execution:** Direct prompt via `run_claude()` (not a skill)
**Timeout:** 15 minutes

### What it does

1. Gathers the full week's data:
   - Meetings synced via `/int-sync-meetings`
   - Tasks via `/prod-review-todoist` (completed, overdue, next week)
   - Next week's agenda via Google Calendar
   - Memory consolidation from the week's daily logs
2. Generates reports in two formats:
   - **HTML** using template `.claude/templates/html/weekly-review.html`
   - **Markdown** using template `.claude/templates/weekly-review.md`

### Output

Dual-format report saved to `workspace/daily-logs/`:
- `[C] YYYY-WXX-weekly-review.html`
- `[C] YYYY-WXX-weekly-review.md`

---

## The Daily Loop

The core routines form a continuous feedback loop:

```
07:00  Good Morning ──> orients the day (reads yesterday's EOD + memory)
         │
      [you work, agents run, meetings happen]
         │
21:00  End of Day ────> captures everything that happened
21:00  Daily Backup ──> exports all user data as ZIP (+ S3 if configured)
21:15  Memory Sync ───> extracts facts, updates memory system
         │
      [overnight]
         │
07:00  Good Morning ──> reads updated memory, starts fresh
         ...
08:00  Weekly Review ─> (Fridays) big picture summary
09:00  Memory Lint ───> (Sundays) health check on memory base
```

Each routine feeds the next. Good Morning reads what End of Day wrote. Memory Sync propagates what End of Day captured. Weekly Review consolidates what Memory Sync organized. Memory Lint ensures the whole system stays healthy.

## Running Manually

```bash
make morning         # Good Morning
make eod             # End of Day
make backup-daily    # Daily Backup
make memory          # Memory Sync
make memory-lint     # Memory Lint
make weekly          # Weekly Review
```

All core routines can be run at any time -- they are not restricted to their scheduled times.
