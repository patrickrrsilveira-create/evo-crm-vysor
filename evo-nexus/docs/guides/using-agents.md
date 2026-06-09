# Using Agents

EvoNexus ships with 16 specialized agents, each owning a specific domain. You don't need to memorize commands — just describe what you need and Claude routes to the right agent.

## Getting Started

1. Open Claude Code in the workspace directory
2. Start talking — agents activate automatically based on context

## Quick Examples

### Morning Briefing

Just say:

```
good morning
```

Clawdia checks your calendar, emails, and tasks, then gives you a prioritized briefing for the day.

### Financial Pulse

```
/flux financial pulse
```

Flux queries Stripe and Omie, then generates an HTML report with MRR, charges, churn, receivables, and payables.

### GitHub Review

```
/atlas check github
```

Atlas reviews open PRs, community issues, stars/forks, and recent releases across all Evolution repositories.

### Community Pulse

```
/pulse community pulse
```

Pulse reads Discord and WhatsApp messages from the last 24 hours, analyzes sentiment, identifies top topics, and flags unanswered support questions.

### Social Media Content

```
/pixel write a post about our new release
```

Pixel drafts platform-specific posts (LinkedIn, Twitter/X, Threads) following your brand voice, content pillars, and audience context.

### Strategy Digest

```
/sage strategy digest
```

Sage consolidates financial, product, community, and market data into a weekly strategic summary with insights and recommended actions.

## How Agent Routing Works

You don't have to use slash commands. Claude reads your intent and picks the right agent:

| You say... | Agent activated |
|---|---|
| "check my inbox" | Clawdia |
| "how's the MRR this month?" | Flux |
| "any open PRs?" | Atlas |
| "what's the community talking about?" | Pulse |
| "draft a LinkedIn post" | Pixel |
| "review our OKRs" | Sage |
| "schedule a meeting" | Clawdia |
| "who are our competitors?" | Sage |
| "check my health data" | Kai |
| "create a course module" | Mentor |

## All 16 Agents

| Agent | Command | Domain |
|---|---|---|
| **Clawdia** | `/clawdia` | Hub — calendar, emails, tasks, decisions |
| **Flux** | `/flux` | Finance — cash flow, Stripe, Omie |
| **Atlas** | `/atlas` | Projects — GitHub, Linear, milestones |
| **Pulse** | `/pulse` | Community — Discord, WhatsApp, sentiment |
| **Pixel** | `/pixel` | Social media — content, analytics |
| **Sage** | `/sage` | Strategy — OKRs, roadmap, competitive |
| **Kai** | `/kai` | Personal — health, habits, routine |
| **Nex** | `/nex` | Sales — pipeline, proposals |
| **Mentor** | `/mentor` | Courses — Evo Academy, modules |
| **Oracle** | `/oracle` | Knowledge — workspace docs, how-to, config |
| **Mako** | `/mako` | Marketing — campaigns, SEO, brand |
| **Aria** | `/aria` | HR / People — recruiting, onboarding, performance |
| **Zara** | `/zara` | Customer Success — triage, escalation, health |
| **Lex** | `/lex` | Legal — contracts, compliance, NDA, risk |
| **Nova** | `/nova` | Product — specs, roadmaps, metrics |
| **Dex** | `/dex` | Data / BI — analysis, SQL, dashboards |

## Tips

- **Agents share memory.** Context from one session is available in the next via the `memory/` directory. If you told Clawdia about a deadline, Sage can reference it when generating strategy.
- **Combine agents naturally.** Say "check github and then write a post about the latest release" — Atlas reviews repos, then Pixel drafts the post.
- **Agents respect your preferences.** All responses are in pt-BR, use Brasilia timezone, and follow your professional tone.
- **Check agent definitions** in `.claude/agents/` to see exactly what each agent can do.
