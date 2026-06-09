# workspace/learning — Knowledge Retention Layer

Spaced repetition system for EvoNexus, built on SM-2 algorithm over plain markdown files.

## Purpose

Closes the retention gap: EvoNexus captures knowledge (`memory/`, `agent-memory/`, `meetings/`) but had no mechanism to bring a fact back before it was forgotten. This layer adds **deliberate capture + spaced review + active recall**.

## Directory Layout

```
workspace/learning/
├── README.md          ← this file (versioned)
├── facts/             ← individual fact files in SM-2 frontmatter format (gitignored)
├── decks/             ← deck metadata and optional deck-level notes (gitignored)
└── .state/
    └── review-log.jsonl  ← append-only review history (gitignored)
```

**`facts/` is gitignored** — fact content may include text from third-party articles,
books, or meetings that cannot be committed to the MIT-licensed repo.

## Fact File Format

Each file in `facts/` follows the pattern `{YYYY-MM-DD}-{slug}.md`:

```markdown
---
id: 2026-04-11-claude-skills-markdown
source: https://... (or "manual")
deck: claude-skills
created: 2026-04-11
next_review: 2026-04-12
interval: 1
ease: 2.5
reps: 0
lapses: 0
---

**Fact:** [The atomic fact in one or two sentences.]

**Why it matters:** [Why Davidson should remember this.]

**Retrieval Q:** [A question that can be answered from memory, used during review.]
```

## SM-2 Algorithm (implemented in skill prompts)

| Grade | Effect |
|-------|--------|
| Again (0) | `reps=0`, `interval=1`, `ease=max(1.3, ease-0.2)`, `lapses++` |
| Hard  (3) | `interval=round(interval*1.2)`, `ease=max(1.3, ease-0.15)`, `reps++` |
| Good  (4) | reps=0→interval=1; reps=1→interval=6; else `interval=round(interval*ease)`. `reps++` |
| Easy  (5) | Like Good + `interval=round(interval*1.3)` extra, `ease+=0.15`, `reps++` |

`next_review = review_date + interval days`

## Skills

| Skill | What it does |
|-------|-------------|
| `learn-capture` | Extracts 1-5 atomic facts from pasted text or URL (paste-text mode) |
| `learn-review`  | Reviews facts with `next_review <= today`, applies SM-2, updates frontmatter |
| `learn-quiz`    | Generates open questions from a deck — read-only, does NOT update SM-2 state |
| `learn-stats`   | Reports totals, overdue count, retention rate, active decks, weekly adds |

## Language Policy

- Skills and this README are in **English** (EvoNexus repo policy — open source).
- Fact content is in **pt-BR** by default (`workspace.language`).
- All user content (`facts/`, `decks/`, `.state/`) is gitignored — privacy and licensing.

## Privacy & Licensing

Facts may contain summaries of third-party content (articles, books). These are **never committed** to the repo. The `.gitignore` at the repo root explicitly excludes `workspace/learning/facts/`, `workspace/learning/decks/`, and `workspace/learning/.state/`.

Only this `README.md` is tracked in git.
