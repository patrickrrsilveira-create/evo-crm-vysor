# Nex — Sales

**Command:** `/nex-sales` | **Color:** red | **Model:** Sonnet

Nex is the sales agent responsible for pipeline management, lead qualification, proposal drafting, and commercial KPIs. It operates at L1 Observer level — meaning it generates reports and drafts for approval but never sends anything directly. Use Nex for all commercial and sales-related tasks.

## When to Use

- You need a snapshot of the current sales pipeline status
- You want to prepare or review a commercial proposal
- You need to qualify a new lead or score existing ones
- You want to track commercial KPIs (conversion rate, deal velocity, win rate)
- You need to evaluate a vendor or compare vendor proposals

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `int-evo-crm` | Access CRM data — contacts, conversations, pipelines, stages, and labels |
| `ops-vendor-review` | Evaluate vendors with cost analysis, risk assessment, and recommendation |

## Example Interactions

```
/nex-sales what's the current state of our sales pipeline?
/nex-sales prepare a proposal for the HostGator partnership renewal
/nex-sales qualify this lead: mid-size SaaS company, 50 employees, looking for WhatsApp automation
```

## Routines

This agent is not used in any scheduled routines.

## Memory

Persistent memory at `.claude/agent-memory/nex/`. Stores project context, user preferences, and feedback across sessions.
