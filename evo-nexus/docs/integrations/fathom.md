# Fathom Integration

Fathom records and transcribes meetings, providing summaries, transcripts, and action items. The @clawdia agent uses Fathom to keep track of discussions, decisions, and follow-ups across all calls.

## Setup

### 1. Get Your API Key

1. Go to [fathom.video](https://fathom.video)
2. Navigate to **Settings > API**
3. Copy your API key (requires Fathom Team or Enterprise plan)

### 2. Configure .env

```env
FATHOM_API_KEY=your_api_key_here
```

### 3. Test the Connection

```bash
make sync
```

This fetches recent meetings from Fathom, saves transcripts and summaries locally, and updates the meeting index.

## Available Commands

The `int-fathom` skill queries the Fathom API via a shell script:

| Command | What it does |
|---|---|
| `meetings` | List recent meetings |
| `meetings --include-summary` | List meetings with AI summaries |
| `meetings --include-actions` | List meetings with action items |
| `meetings --include-transcript` | List meetings with full transcripts |
| `meetings --after DATE --before DATE` | Filter by date range |
| `meetings --recorded-by EMAIL` | Filter by recorder |
| `meetings --team NAME` | Filter by team |
| `summary RECORDING_ID` | Get summary for a specific recording |
| `transcript RECORDING_ID` | Get transcript for a specific recording |
| `teams` | List teams |
| `members` | List team members |

## Key Data Points

- Meeting title, date, duration, and participants
- AI-generated summaries (markdown formatted)
- Action items with assignees
- Full transcripts with speaker labels
- Team and participant filtering

## Skills That Use Fathom

| Skill | What it does |
|---|---|
| `int-fathom` | Direct Fathom API queries -- meetings, summaries, transcripts |
| `int-sync-meetings` | Full sync pipeline -- fetch, deduplicate, save, and index meetings |

## Automated Routines

| Routine | Schedule | Make command |
|---|---|---|
| Sync Meetings | Every 30 min | `make sync` |
