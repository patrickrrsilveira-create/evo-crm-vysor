# Evolution API Integration

Evolution API integration provides WhatsApp messaging capabilities: manage instances, send and receive messages, handle chats, contacts, and groups.

## Setup

### 1. Get Your API Credentials

1. Deploy your own Evolution API instance or use a hosted one — see [github.com/EvolutionAPI/evolution-api](https://github.com/EvolutionAPI/evolution-api)
2. After deployment, copy your **API URL** and **API Key** from the instance settings

### 2. Configure .env

```env
EVOLUTION_API_URL=https://your-evolution-instance.com
EVOLUTION_API_KEY=your_api_key_here
```

Set this via terminal (`nano .env`) or the dashboard .env editor.

### 3. Test the Connection

```
> List my WhatsApp instances
```

Claude uses the `int-evolution-api` skill to call the API and return your active instances.

## What It Does

The Evolution API integration provides:

| Feature | Description |
|---|---|
| **Instances** | List, create, connect, and manage WhatsApp instances |
| **Messages** | Send text, media, location, contact, button, list, and poll messages |
| **Chats** | Query chat history and contact information |
| **Groups** | Create, manage, and query WhatsApp groups |
| **Webhooks** | Configure webhooks for real-time event notifications |
| **Connection** | QR code generation and connection status monitoring |

## Skills That Use Evolution API

| Skill | What it does |
|---|---|
| `int-evolution-api` | Direct API access — manage instances, send messages, query chats and groups, configure webhooks |
| `pulse-daily` | Community pulse includes WhatsApp group activity via Evolution API |
| `pulse-weekly` | Weekly community report aggregates WhatsApp data |

## Automated Routines

| Routine | Schedule | Make command |
|---|---|---|
| Community Pulse | 20:00 BRT daily | `make community` |
| Community Weekly | Monday 09:30 BRT | `make community-week` |
| Community Monthly | 1st of month | `make community-month` |

## Evolution API + Evolution Go

Evolution API (Node.js) and Evolution Go are two implementations of the same WhatsApp messaging platform. You can use either or both — they share the same API contract. Evolution Go is a newer, high-performance Go implementation.

If you run both, configure separate env vars for each and specify which instance to use when sending messages.

## Security Notes

- The API key has full access to your Evolution API instance. Use instance-scoped keys when possible.
- Never commit `.env` to version control. The `.gitignore` already excludes it.
- The dashboard masks the key in the .env editor (displayed as `*****`).
