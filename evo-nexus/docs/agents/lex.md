# Lex — Legal & Compliance

**Command:** `/lex-legal` | **Color:** purple | **Model:** Sonnet

Lex is the Legal and Compliance agent, handling contract review, NDA triage, LGPD compliance, risk assessment, legal briefs, vendor checks, and signature requests. He provides clause-by-clause analysis, flags regulatory concerns, and ensures legal workflows move efficiently from review to execution.

## When to Use

- You need to review a contract and flag risky or non-standard clauses
- An NDA arrived from a vendor and needs quick green/yellow/red classification
- You want to run a LGPD or regulatory compliance check on a product feature or initiative
- You need to assess the legal risk of a vendor agreement or partnership
- You want to draft a legal brief or prepare for a meeting with legal context

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `legal-review-contract` | Clause-by-clause contract review against negotiation playbook |
| `legal-compliance-check` | Regulatory assessment for proposed actions or features |
| `legal-triage-nda` | Quick NDA classification as green, yellow, or red |
| `legal-brief` | Contextual legal briefings for daily work or incident response |
| `legal-response` | Templated legal responses with escalation checks |
| `legal-risk-assessment` | Severity x likelihood risk matrix with escalation criteria |
| `legal-meeting-briefing` | Meeting prep with legal context and action item tracking |
| `legal-signature-request` | E-signature workflow with pre-signature checklist |
| `legal-vendor-check` | Vendor agreement audit across connected systems |
| `ops-risk-assessment` | Operational risk evaluation and mitigation planning |
| `ops-compliance-tracking` | Audit readiness tracking and compliance requirement monitoring |

## Example Interactions

```
/lex-legal review this contract and flag risky clauses
/lex-legal triage the NDA we received from vendor X
/lex-legal run a LGPD compliance check on the new analytics feature
```

## Routines

This agent is not used in any scheduled routines.

## Memory

Persistent memory at `.claude/agent-memory/lex/`. Stores project context, user preferences, and feedback across sessions.
