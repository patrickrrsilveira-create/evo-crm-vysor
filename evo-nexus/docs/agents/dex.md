# Dex — Data & BI

**Command:** `/dex-data` | **Color:** yellow | **Model:** Sonnet

Dex is the Data and BI agent, specializing in data analysis, SQL queries, interactive dashboards (Chart.js), statistical analysis, data profiling, and validation. He connects to Stripe, Omie, and Licensing APIs to turn raw data into actionable insights, visualizations, and quality-checked deliverables.

## When to Use

- You need to analyze a financial trend like MRR over the last 3 months
- You want an optimized SQL query for a specific data question (e.g., churned customers)
- You need an interactive HTML dashboard for licensing growth or revenue metrics
- You want to run statistical analysis on conversion rates or experiment results
- A dataset needs profiling or validation before publishing or sharing

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `data-analyze` | Answers ad-hoc data questions from quick lookups to full reports |
| `data-build-dashboard` | Builds interactive HTML dashboards with charts, filters, and tables |
| `data-write-query` | Writes optimized SQL for PostgreSQL with best practices |
| `data-explore` | Profiles datasets for shape, quality, nulls, and distributions |
| `data-create-viz` | Creates professional Python visualizations in Evolution theme |
| `data-statistical-analysis` | Runs descriptive stats, trend analysis, and hypothesis testing |
| `data-validate` | Pre-delivery QA for methodology, accuracy, and bias checks |
| `int-stripe` | Accesses Stripe financial data (charges, subscriptions, MRR) |
| `int-omie` | Accesses Omie ERP data (clients, invoices, financials) |
| `int-licensing` | Queries open source telemetry (instances, geo, versions) |

## Example Interactions

```
/dex-data analyze MRR trend for the last 3 months
/dex-data build a licensing growth dashboard with geo breakdown
/dex-data run statistical analysis on trial-to-paid conversion rates
```

## Routines

This agent is not used in any scheduled routines.

## Memory

Persistent memory at `.claude/agent-memory/dex/`. Stores project context, user preferences, and feedback across sessions.
