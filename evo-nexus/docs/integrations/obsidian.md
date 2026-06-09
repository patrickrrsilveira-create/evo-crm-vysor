# Obsidian Integration

Obsidian is used as the local vault for notes, knowledge base, and workspace management. Claude interacts with Obsidian via a CLI that communicates with the running Obsidian instance.

## Setup

### 1. Install Obsidian

Download and install from [obsidian.md](https://obsidian.md). Open your vault.

### 2. Ensure the CLI Is Available

The `obsidian` CLI must be installed and accessible from the terminal. Run:

```bash
obsidian help
```

If it works, the integration is ready. No `.env` variables are required.

## Available Commands

The `obs-obsidian-cli` skill provides vault interaction:

| Command | What it does |
|---|---|
| `obsidian read file="Note"` | Read a note by name |
| `obsidian create name="Note" content="..."` | Create a new note |
| `obsidian append file="Note" content="..."` | Append content to a note |
| `obsidian search query="term" limit=10` | Search vault content |
| `obsidian daily:read` | Read today's daily note |
| `obsidian daily:append content="..."` | Append to today's daily note |
| `obsidian property:set name="key" value="val"` | Set frontmatter properties |
| `obsidian tasks daily todo` | List daily tasks |
| `obsidian tags sort=count counts` | List tags with counts |
| `obsidian backlinks file="Note"` | Show backlinks for a note |

## Related Skills

EvoNexus includes several Obsidian-specific skills:

| Skill | What it does |
|---|---|
| `obs-obsidian-cli` | Core vault operations -- read, create, search, manage notes |
| `obs-obsidian-markdown` | Obsidian Flavored Markdown -- wikilinks, callouts, embeds, properties |
| `obs-json-canvas` | Create and edit `.canvas` files (nodes, edges, groups) |
| `obs-obsidian-bases` | Create and edit `.base` files (database views, filters, formulas) |
| `obs-defuddle` | Extract clean markdown from web pages using Defuddle CLI |

## Key Features

- **File targeting:** Use `file=<name>` (wikilink-style) or `path=<path>` (exact path from vault root)
- **Vault targeting:** Use `vault="Name"` to target a specific vault
- **Silent mode:** Add `silent` flag to prevent files from opening in the UI
- **Clipboard:** Add `--copy` to copy output to clipboard
- **Plugin development:** Supports reloading plugins, running JavaScript, capturing errors, and taking screenshots
