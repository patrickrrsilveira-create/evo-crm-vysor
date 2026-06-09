# EvoNexus Roadmap

> Unofficial toolkit for Claude Code — AI-powered business operating system.
>
> This roadmap is updated regularly. Want to vote or suggest? [Open a discussion](https://github.com/EvolutionAPI/evo-nexus/discussions) or [create an issue](https://github.com/EvolutionAPI/evo-nexus/issues).

---

## Legend

| Symbol | Meaning |
|--------|---------|
| `[ ]` | Not started |
| `[x]` | Done |
| `⚠️` | Breaking change |
| `🔥` | High priority |
| `💡` | Needs design discussion first |

---

## v0.4 — Foundation & Stability

> Fix, secure, and improve what already exists before growing.

### Skills

- [x] 🔥 **Evolution product skills** — `int-evolution-api` (33 commands), `int-evolution-go` (24 commands), `int-evo-crm` (48 commands) for managing instances, messages, contacts, conversations, pipelines via REST API.
- [x] **Version indicator & update alerts** — show current version in dashboard sidebar, alert when new GitHub releases are available.

### Developer Experience

- [x] 🔥 **CLI installer** — `npx @evoapi/evo-nexus` — clones repo, installs deps, runs interactive setup wizard.
- [x] **Full Docker install** — `docker compose up dashboard` with multi-stage Dockerfile + GitHub Actions CI pushing to GHCR.
- [x] **Update checker** — dashboard checks GitHub releases and shows upgrade notification.
- [x] **settings.json** — project-level permissions (allow/deny), hooks configuration, thinking mode enabled.
- [x] **CLAUDE.md split** — reduced from 263 to 128 lines. Detailed config moved to `.claude/rules/` (agents, integrations, routines, skills).
- [x] **Inner-loop commands** — `/status` (workspace status) and `/review` (recent changes + next actions).

### Dashboard UX

- [x] **Sidebar reorganization** — 5 collapsible groups (Main, Operations, Data, System, Admin) with localStorage persistence.
- [x] **Active agent visualization** — Claude Code hooks track agent launches via `PreToolUse` events, writing to `agent-status.json`. Dashboard polls `/api/agents/active` and shows "RUNNING" badges with pulse animation on agent cards and overview.
- [x] **Agents page redesign** — unique icons and accent colors per agent, status dots, slash command badges, memory count pills, hover glow effects.
- [x] **Overview page redesign** — stat cards with icons and trend indicators, active agents bar, quick actions row, improved reports and routines tables.

### Agent Generalization

- [x] **Agent prompts generalized** — all 9 agent prompts cleaned of hardcoded personal references. User-specific context preserved in `_improvements.md` per agent memory folder.

---

## v0.5 — Extensibility & Ecosystem

> Make EvoNexus composable and self-extending.

### Agent System

- [x] 🔥 **Generalize existing agents** — all 9 agent prompts generalized. User-specific context preserved in `_improvements.md` per agent memory folder. Adapter patterns documented as future work.
- [x] 🔥 **New business agents** — expand functional coverage:
  - [x] **Marketing Agent** — orchestrate existing `mkt-*` skills, attribution, budget, full funnel
  - [x] **HR / People Agent** — onboarding, 1:1s, performance reviews, hiring pipeline
  - [x] **Customer Success Agent** — health score, churn prediction, NPS/CSAT, client onboarding
  - [x] **Legal / Compliance Agent** — contracts, renewals, GDPR/LGPD, compliance checklists
  - [x] **Product Agent** — discovery, feature prioritization (RICE/ICE), PLG metrics, feedback loop
  - [x] **Data / BI Agent** — cross-area consolidated dashboard, unified KPIs, alerts, trend analysis
- [x] **Custom agents** — `custom-` prefix, gitignored, auto-discovered by dashboard (core/custom badges), `create-agent` skill, `create-command` skill
- [x] **Help agent (Oracle)** — `/oracle` answers questions about the workspace by reading the actual docs. No RAG needed — reads `docs/llms-full.txt` and source files directly

### Routines & Scheduling

- [x] 🔥 **Trigger registry** — define and manage named triggers (webhook, cron, event-based) that invoke skills or routines
- [x] **Non-recurrent scheduled actions** — one-off scheduled tasks (e.g., "post this on LinkedIn Friday at 10am") without creating a full routine
- [x] **Systematic routines** — pure Python routines via `run_script()` — no AI, no tokens, no cost. `create-routine` skill generates the code

### Integrations

- [ ] **Complete Obsidian integration** — finish `obs-*` skills: bidirectional sync, canvas, bases, CLI

### Import / Export

- [x] **Backup system** — export workspace state as ZIP (agents, skills, routines, memory, config); import to restore. Support local, git, and cloud bucket targets.

---

## v0.12 — Engineering Layer

> Add a complete software development team alongside the business agents.

### Engineering Layer (delivered)

- [x] 🔥 **19 engineering agents** — complete dev team derived from [oh-my-claudecode](https://github.com/yeachan-heo/oh-my-claudecode) (MIT, Yeachan Heo). Reasoning tier (opus): apex-architect, echo-analyst, compass-planner, raven-critic, lens-reviewer, zen-simplifier, vault-security. Execution tier (sonnet): bolt-executor, hawk-debugger, grid-tester, probe-qa, oath-verifier, trail-tracer, flow-git, scroll-docs, canvas-designer, prism-scientist. Speed tier (haiku): scout-explorer, quill-writer.
- [x] 🔥 **25 `dev-*` skills** — Tier 1 orchestration (15): `dev-autopilot`, `dev-plan`, `dev-ralplan`, `dev-deep-interview`, `dev-deep-dive`, `dev-external-context`, `dev-trace`, `dev-verify`, `dev-ultraqa`, `dev-visual-verdict`, `dev-ai-slop-cleaner`, `dev-sciomc`, `dev-team`, `dev-ccg`, `dev-ralph`. Tier 2 setup (5): `dev-mcp-setup`, `dev-deepinit`, `dev-project-session-manager`, `dev-configure-notifications`, `dev-release`. Tier 3 utilities (5): `dev-cancel`, `dev-remember`, `dev-ask`, `dev-learner`, `dev-skillify`.
- [x] **15 dev templates** — `.claude/templates/dev-*.md` for each agent's primary output (architecture-decision, work-plan, code-review, bug-report, verification-report, deep-interview-spec, security-audit, test-strategy, trace-report, explore-report, design-spec, analysis-report, research-brief, critique, simplification-report).
- [x] **`workspace/development/` folder** — 7 subfolders (architecture, plans, specs, reviews, debug, verifications, research) for engineering layer artifacts. Distinct from `workspace/projects/` (active git repos).
- [x] **Two-layer architecture documented** — `.claude/rules/agents.md`, `CLAUDE.md`, `docs/agents/overview.md`, `docs/agents/engineering-layer.md`, `docs/architecture.md`, `docs/introduction.md`, site `Home.tsx`.
- [x] **Open source attribution** — `NOTICE.md` at repo root with full MIT license, version pinned (v4.11.4), modifications listed. Credits in `README.md`.
- [x] **Pattern compliance** — all 19 engineering agents follow the EvoNexus standard: rich frontmatter (Examples + commentary), Workspace Context, Shared Knowledge Base, Working Folder, Identity, Anti-patterns, Domain, How You Work, Skills You Can Use, Handoffs, Output Format, Continuity. Verified by `@lens-reviewer` (3 fixes applied: oath-verifier disallowedTools, raven-critic and trail-tracer Skills section).

### Cross-layer pipelines (now possible)

- [x] **End-to-end implementation** — `dev-autopilot` orchestrates spec → plan → code → QA → validation across multiple engineering agents.
- [x] **High-stakes consensus planning** — `dev-ralplan` runs Planner/Architect/Critic consensus loop with RALPLAN-DR structured deliberation.
- [x] **Bug investigation** — `@trail-tracer` (multi-hypothesis) → `@hawk-debugger` (root cause + minimal fix) → `@oath-verifier` (regression check).
- [x] **Pre-merge gate** — `@lens-reviewer` (code quality) → `@vault-security` (OWASP audit) → `dev-ultraqa` (build/test/fix loop) → `@oath-verifier` (acceptance criteria).

---

## v1.0 — Community & Growth

> Community adoption, discoverability, and self-sustaining ecosystem.

### Community & Docs

- [x] 🔥 **Public roadmap** — this file. Community input welcome via [discussions](https://github.com/EvolutionAPI/evo-nexus/discussions).
- [x] **Telegram & Discord channels** — activate community channels, document in README and docs site.
- [ ] **In-app tutorials** — contextual tutorials surfaced inside the dashboard, not just external docs.
- [x] **Resume Claude sessions in chat** — list active/resumable Claude sessions in dashboard chat with `--resume` support.

### Development

- [ ] **Testing framework** — define and implement test strategy for skills, routines, and agent behaviors; prevent regressions.

---

## Contributing

Want to help? Pick any `[ ]` item and:

1. Check [open issues](https://github.com/EvolutionAPI/evo-nexus/issues)
2. Read [CONTRIBUTING.md](CONTRIBUTING.md)
3. For `💡` items, open a [discussion](https://github.com/EvolutionAPI/evo-nexus/discussions) first — design is still open

---

*Last updated: 2026-04-10 — [Evolution Foundation](https://evolutionfoundation.com.br)*
