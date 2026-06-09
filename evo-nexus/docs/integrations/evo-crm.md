# Evo CRM Integration

Evo CRM integration connects EvoNexus to the AI-powered CRM platform for managing contacts, conversations, pipelines, and labels.

## Setup

### 1. Get Your API Credentials

1. Deploy or access your Evo CRM instance — see [github.com/EvolutionAPI/evo-crm-community](https://github.com/EvolutionAPI/evo-crm-community)
2. Generate an API token from the CRM admin settings
3. Copy your **CRM URL** and **API Token**

### 2. Configure .env

```env
EVO_CRM_URL=https://your-evo-crm-instance.com
EVO_CRM_TOKEN=your_api_token_here
```

Set this via terminal (`nano .env`) or the dashboard .env editor.

### 3. Test the Connection

```
> List my CRM contacts
```

Claude uses the `int-evo-crm` skill to query the CRM API and return your data.

## What It Does

The Evo CRM integration provides:

| Feature | Description |
|---|---|
| **Contacts** | List, search, create, and update contacts |
| **Conversations** | View and manage customer conversations across channels |
| **Messages** | Read conversation messages and history |
| **Inboxes** | Manage communication inboxes (WhatsApp, email, web chat) |
| **Pipelines** | View and manage sales pipelines with stages and items |
| **Labels** | Organize contacts and conversations with labels |

## Skills That Use Evo CRM

| Skill | What it does |
|---|---|
| `int-evo-crm` | Direct CRM API access — query contacts, conversations, messages, inboxes, pipelines (stages + items), labels. Supports filtering, pagination, and CRUD operations. |

## Use Cases

### Sales Pipeline Management
```
> Show me all pipeline items in the "Negotiation" stage
> Move deal X to "Closed Won"
```

### Customer Lookup
```
> Find contacts from company Acme
> Show recent conversations with Davidson
```

### Conversation History
```
> Show the last 10 messages in conversation #42
> List all open conversations in the WhatsApp inbox
```

## Evo CRM + Evolution API

Evo CRM works best alongside Evolution API or Evolution Go. The CRM manages customer relationships and pipelines, while the Evolution messaging APIs handle WhatsApp communication. Together, they provide a complete customer engagement stack:

- **Evo CRM** — contacts, pipelines, conversation management
- **Evolution API / Go** — message delivery and WhatsApp instance management

## Security Notes

- The API token has full access to your CRM data. Use scoped tokens if available.
- Never commit `.env` to version control. The `.gitignore` already excludes it.
- The dashboard masks the token in the .env editor (displayed as `*****`).
