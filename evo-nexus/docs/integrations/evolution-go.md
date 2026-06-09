# Evolution Go Integration

Evolution Go is a high-performance Go implementation of the WhatsApp messaging platform. It provides the same capabilities as Evolution API with improved performance and resource efficiency.

## Setup

### 1. Get Your API Credentials

1. Deploy your own Evolution Go instance — see [github.com/EvolutionAPI/evolution-go](https://github.com/EvolutionAPI/evolution-go)
2. After deployment, copy your **API URL** and **API Key** from the instance settings

### 2. Configure .env

```env
EVOLUTION_GO_URL=https://your-evolution-go-instance.com
EVOLUTION_GO_KEY=your_api_key_here
```

Set this via terminal (`nano .env`) or the dashboard .env editor.

### 3. Test the Connection

```
> List my Evolution Go instances
```

Claude uses the `int-evolution-go` skill to call the API and return your active instances.

## What It Does

The Evolution Go integration provides:

| Feature | Description |
|---|---|
| **Instances** | List, create, connect, and manage WhatsApp instances |
| **Messages** | Send text, media, location, contact, link, sticker, and poll messages |
| **Reactions** | React to, edit, and delete messages |
| **Connection** | QR code generation and connection status monitoring |
| **Performance** | Lower memory footprint and faster response times than the Node.js version |

## Skills That Use Evolution Go

| Skill | What it does |
|---|---|
| `int-evolution-go` | Direct API access — manage instances, send messages, react/edit/delete messages, check connection status |
| `pulse-daily` | Community pulse includes WhatsApp group activity |
| `pulse-weekly` | Weekly community report aggregates WhatsApp data |

## Automated Routines

| Routine | Schedule | Make command |
|---|---|---|
| Community Pulse | 20:00 BRT daily | `make community` |
| Community Weekly | Monday 09:30 BRT | `make community-week` |
| Community Monthly | 1st of month | `make community-month` |

## Evolution Go vs Evolution API

| Aspect | Evolution API | Evolution Go |
|---|---|---|
| Language | Node.js / TypeScript | Go |
| Performance | Standard | Higher throughput, lower memory |
| API contract | REST | REST (compatible) |
| Maturity | Established, widely deployed | Newer, actively developed |

Both share the same API patterns. You can use either or both depending on your needs.

## Security Notes

- The API key has full access to your Evolution Go instance. Use instance-scoped keys when possible.
- Never commit `.env` to version control. The `.gitignore` already excludes it.
- The dashboard masks the key in the .env editor (displayed as `*****`).
