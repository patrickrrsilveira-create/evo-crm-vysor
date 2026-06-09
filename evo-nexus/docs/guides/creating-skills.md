# Creating Skills

Skills are reusable instruction sets that teach Claude how to perform specific tasks. EvoNexus includes 126 built-in skills — you can create your own.

## Using the Skill Creator

The fastest way:

```
create a skill for weekly team reports
```

The `skill-creator` plugin walks you through:

1. **Name and prefix** — follows the naming convention (e.g., `prod-weekly-team-report`)
2. **Description** — the trigger phrase that activates the skill
3. **Instructions** — what Claude should do when the skill is invoked

## Manual Creation

Skills live in `.claude/skills/` organized by prefix.

### File Structure

```
.claude/skills/
  prefix-name/
    SKILL.md
```

### SKILL.md Format

```markdown
---
name: weekly-team-report
description: Generate a weekly team status report from Linear, GitHub, and meeting notes. Use when the user says 'team report', 'weekly team update', or 'what did the team do this week'.
---

# Weekly Team Report

## Steps

1. Query Linear for completed issues this week
2. Query GitHub for merged PRs
3. Check meeting summaries from the week
4. Generate a structured report with:
   - Completed work by team member
   - In-progress items
   - Blockers
   - Next week priorities

## Output

Generate an HTML report using the Evolution brand template.
Save to `workspace/daily-logs/` with filename `[C] Team Report YYYY-MM-DD.html`.

## Rules

- Always group by team member
- Include links to Linear issues and GitHub PRs
- Highlight blockers in red
```

### Key Points

- **The `description` field is the trigger.** Claude matches user intent against skill descriptions to decide which skill to activate. Make it specific and include common trigger phrases.
- **The `name` field** is the skill identifier.
- **The body** contains the full instructions Claude follows when the skill is activated.

## Naming Convention

Skills use a prefix that indicates their category:

| Prefix | Category | Examples |
|---|---|---|
| `prod-` | Productivity | `prod-morning`, `prod-eod` |
| `fin-` | Finance | `fin-daily-pulse`, `fin-journal-entry` |
| `social-` | Social media | `social-post-writer`, `social-analytics-report` |
| `int-` | Integrations | `int-stripe`, `int-github-review` |
| `pulse-` | Community | `pulse-daily`, `pulse-faq-sync` |
| `sage-` | Strategy | `sage-okr-review`, `sage-strategy-digest` |
| `evo-` | Evo Method | `evo-dev`, `evo-create-prd` |
| `mkt-` | Marketing | `mkt-campaign-plan`, `mkt-seo-audit` |
| `gog-` | Google | `gog-calendar`, `gog-email-triage` |
| `obs-` | Obsidian | `obs-obsidian-cli`, `obs-defuddle` |
| `discord-` | Discord | `discord-get-messages`, `discord-send-message` |

## Example: Creating a "Weekly Report" Skill

1. Create the directory:

```bash
mkdir -p .claude/skills/prod-weekly-report
```

2. Create `SKILL.md`:

```markdown
---
name: weekly-report
description: Generate a consolidated weekly report covering projects, finance, community, and strategy. Use when the user says 'weekly report', 'week summary', 'what happened this week'.
---

# Weekly Report Generator

Consolidate data from all domains into a single weekly summary.

## Data Sources

1. **Projects** — Linear issues completed, GitHub PRs merged
2. **Finance** — Stripe MRR, revenue, churn
3. **Community** — Discord activity, support topics
4. **Strategy** — OKR progress, key decisions

## Output

HTML report saved to `workspace/daily-logs/`.
```

3. Test it — just describe the task:

```
generate a weekly report
```

Claude matches your intent to the skill description and follows the instructions.

## Tips

- **Be specific in descriptions.** "Generate financial reports" is too vague. "Generate financial statements (income statement, balance sheet, cash flow) with period-over-period comparison" tells Claude exactly when to activate.
- **Include trigger phrases.** Add common ways users might ask for the skill: "Use when the user says 'X', 'Y', or 'Z'."
- **Reference other skills.** Skills can call other skills — e.g., a weekly report skill can reference the Stripe integration skill for financial data.
- **Test iteratively.** Create the skill, test it, refine the instructions based on the output quality.
