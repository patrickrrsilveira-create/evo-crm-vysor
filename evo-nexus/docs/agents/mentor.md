# Mentor — Courses & Education

**Command:** `/mentor-courses` | **Color:** purple | **Model:** Sonnet

Mentor is the education agent focused on course creation, learning paths, study plans, and didactic material for the Evo Academy platform. It helps structure educational content from scratch, build progressive learning journeys, and produce lesson material that is clear, engaging, and well-organized.

## When to Use

- You need to structure a new course from objectives to modules and lessons
- You want to create a learning path that guides students through a topic progressively
- You need to build a specific module or lesson about a given subject
- You want to explain a technical concept in a way suitable for a course lesson
- You need a study plan for learning a new technology or framework

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `obs-obsidian-markdown` | Create and edit Obsidian-flavored markdown content with wikilinks, callouts, and properties |
| `obs-obsidian-cli` | Interact with Obsidian vaults — read, create, search, and manage notes |

## Example Interactions

```
/mentor-courses structure a course on building AI agents with the Evolution API
/mentor-courses create a learning path for someone new to WhatsApp automation
/mentor-courses build a module about webhook configuration for intermediate users
```

## Routines

This agent is not used in any scheduled routines.

## Memory

Persistent memory at `.claude/agent-memory/mentor/`. Stores project context, user preferences, and feedback across sessions.
