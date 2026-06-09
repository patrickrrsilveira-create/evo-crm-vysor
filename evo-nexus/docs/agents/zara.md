# Zara — Customer Success

**Command:** `/zara-cs` | **Color:** cyan | **Model:** Sonnet

Zara is the Customer Success agent, responsible for ticket triage, escalation packaging, customer research, draft responses, knowledge base articles, and CRM integration. She ensures support tickets are prioritized correctly, customers get professional replies, and recurring issues are captured as self-service documentation.

## When to Use

- You need to triage open support tickets and assign priority levels (P1-P4)
- You want to draft a professional response to a customer complaint (e.g., downtime, billing issue)
- A bug or incident needs to be escalated to engineering or product with full context
- A resolved ticket should be turned into a knowledge base article for self-service
- You want to identify which customers are at churn risk based on recent interactions

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `cs-ticket-triage` | Classifies and prioritizes support tickets from P1 (critical) to P4 (low) |
| `cs-customer-escalation` | Packages escalations with full context for devs or product teams |
| `cs-customer-research` | Multi-source lookup to gather customer history and context |
| `cs-draft-response` | Generates professional, situation-aware customer replies |
| `cs-kb-article` | Creates knowledge base articles from resolved issues |
| `int-evo-crm` | Accesses CRM data for customer records, conversations, and pipelines |

## Example Interactions

```
/zara-cs triage the open support tickets from today
/zara-cs draft a response for the customer reporting downtime on instance X
/zara-cs escalate the authentication bug to engineering with full context
```

## Routines

This agent is not used in any scheduled routines.

## Memory

Persistent memory at `.claude/agent-memory/zara/`. Stores project context, user preferences, and feedback across sessions.
