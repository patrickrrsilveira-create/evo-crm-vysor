# Flux — Finance & CFO

**Command:** `/flux-finance` | **Color:** orange | **Model:** Sonnet

Flux is the virtual CFO — managing all financial operations including cash flow analysis, income statements, Stripe MRR tracking, Omie ERP integration, monthly closing processes, journal entries, reconciliation, and variance analysis. Flux ensures financial health is always visible and actionable.

## When to Use

- Running the monthly closing process or checking close status
- Reviewing financial health, cash flow, or MRR trends
- Checking pending invoices, accounts payable, or receivable
- Preparing journal entries or reconciling accounts
- Analyzing budget vs. actual variances

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `fin-daily-pulse` | Daily financial snapshot from Stripe and Omie |
| `fin-weekly-report` | Weekly consolidation of revenue, expenses, and cash flow |
| `fin-monthly-close-kickoff` | Initiates month-end closing with checklist and P&L |
| `fin-close-management` | Track close tasks, dependencies, and status |
| `fin-financial-statements` | Generate P&L, balance sheet, and cash flow statements |
| `fin-journal-entry` | Prepare journal entries with debits and credits |
| `fin-journal-entry-prep` | Month-end entry preparation with supporting docs |
| `fin-reconciliation` | GL vs subledger and bank reconciliation |
| `fin-variance-analysis` | Budget vs actual decomposition with narratives |
| `fin-audit-support` | SOX compliance documentation and support |
| `fin-sox-testing` | Control testing sample selection and workpapers |
| `int-stripe` | Stripe API — charges, subscriptions, MRR, churn |
| `int-omie` | Omie ERP — clients, invoices, financials |

## Example Interactions

```
/flux-finance how's our cash flow this month?
/flux-finance start the monthly closing for March
/flux-finance what's our current MRR and churn rate?
```

## Routines

| Time | Routine |
|------|---------|
| Daily 19:00 | Financial Pulse |
| Friday 07:30 | Financial Weekly |
| 1st of month | Monthly Close Kickoff |

## Memory

Persistent memory at `.claude/agent-memory/flux/`. Stores project context, user preferences, and feedback across sessions.
