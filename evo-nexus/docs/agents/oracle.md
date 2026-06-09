# Oracle — Workspace Knowledge

**Command:** `/oracle` | **Color:** amber | **Model:** Sonnet

Oracle is the workspace knowledge agent — it answers questions about EvoNexus itself. Agents, skills, routines, integrations, the dashboard, configuration, and architecture are all in its domain. Oracle reads documentation before answering and never guesses. If it does not find the answer in the docs, it says so.

## When to Use

- Understand how to create a new routine, agent, or skill
- Find out what agents are available and what each one does
- Learn how the scheduler works or how routines are triggered
- Check what skills a specific agent has access to
- Discover how to add a new integration or configure an existing one

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `initial-setup` | Onboard new users — introduce agents, skills, routines, and the dashboard |
| `create-agent` | Guided wizard to create a new custom agent |
| `create-command` | Create a new slash command for Claude Code |
| `create-routine` | Guided wizard to create a new automated routine |

## Example Interactions

```
/oracle how do I create a new routine?
/oracle what skills does Flux have?
/oracle how do I add a new MCP integration?
```

## Routines

This agent is not used in any scheduled routines.

## Memory

Oracle has no persistent agent memory. It reads documentation on every invocation to provide accurate, up-to-date answers rather than building institutional knowledge.
