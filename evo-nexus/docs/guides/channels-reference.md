# Channels Reference — Building Custom Channels

Technical reference for building an MCP server that pushes webhooks, alerts, and chat messages into a Claude Code session.

> To use existing channels (Telegram, Discord, iMessage), see the [Channels Guide](channels.md).

## Overview

A channel is an MCP server that runs on the same machine as Claude Code. Claude Code spawns it as a subprocess and communicates over stdio. Your channel server bridges external systems and the session:

- **Chat platforms** (Telegram, Discord): your plugin runs locally and polls the platform's API. When someone DMs the bot, the plugin receives the message and forwards it to Claude.
- **Webhooks** (CI, monitoring): your server listens on a local HTTP port. External systems POST to that port, and your server pushes the payload to Claude.

## Requirements

- [`@modelcontextprotocol/sdk`](https://www.npmjs.com/package/@modelcontextprotocol/sdk) (MCP SDK)
- Node.js-compatible runtime: [Bun](https://bun.sh), [Node](https://nodejs.org), or [Deno](https://deno.com)

Your server must:

1. Declare the `claude/channel` capability so Claude Code registers a notification listener
2. Emit `notifications/claude/channel` events when something happens
3. Connect via [stdio transport](https://modelcontextprotocol.io/docs/concepts/transports#standard-io)

## Example: Webhook Receiver

A single-file server that listens for HTTP requests and forwards them into your Claude Code session.

### 1. Create the Project

```bash
mkdir webhook-channel && cd webhook-channel
bun add @modelcontextprotocol/sdk
```

### 2. Write the Server

```ts
// webhook.ts
#!/usr/bin/env bun
import { Server } from '@modelcontextprotocol/sdk/server/index.js'
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js'

const mcp = new Server(
  { name: 'webhook', version: '0.0.1' },
  {
    capabilities: { experimental: { 'claude/channel': {} } },
    instructions: 'Events from the webhook channel arrive as <channel source="webhook" ...>. Read them and act, no reply expected.',
  },
)

await mcp.connect(new StdioServerTransport())

Bun.serve({
  port: 8788,
  hostname: '127.0.0.1',
  async fetch(req) {
    const body = await req.text()
    await mcp.notification({
      method: 'notifications/claude/channel',
      params: {
        content: body,
        meta: { path: new URL(req.url).pathname, method: req.method },
      },
    })
    return new Response('ok')
  },
})
```

### 3. Register in MCP Config

```json
// .mcp.json
{
  "mcpServers": {
    "webhook": { "command": "bun", "args": ["./webhook.ts"] }
  }
}
```

### 4. Test

```bash
# Terminal 1: start Claude Code with development flag
claude --dangerously-load-development-channels server:webhook

# Terminal 2: simulate a webhook
curl -X POST localhost:8788 -d "build failed on main: https://ci.example.com/run/1234"
```

The payload arrives as:

```
<channel source="webhook" path="/" method="POST">build failed on main: https://ci.example.com/run/1234</channel>
```

## Server Options

| Field | Type | Description |
|-------|------|-------------|
| `capabilities.experimental['claude/channel']` | `object` | **Required.** Always `{}`. Registers the notification listener |
| `capabilities.experimental['claude/channel/permission']` | `object` | Optional. Declares that the channel can receive permission relay |
| `capabilities.tools` | `object` | Two-way only. Always `{}`. Enables tool discovery |
| `instructions` | `string` | Recommended. Added to Claude's system prompt |

## Notification Format

Emit `notifications/claude/channel` with two params:

| Field | Type | Description |
|-------|------|-------------|
| `content` | `string` | The event body. Delivered as the body of the `<channel>` tag |
| `meta` | `Record<string, string>` | Optional. Each entry becomes an attribute on the `<channel>` tag |

```ts
await mcp.notification({
  method: 'notifications/claude/channel',
  params: {
    content: 'build failed on main',
    meta: { severity: 'high', run_id: '1234' },
  },
})
```

Result in Claude's context:

```
<channel source="your-channel" severity="high" run_id="1234">
build failed on main
</channel>
```

**Note on meta keys:** only letters, digits, and underscores. Keys with hyphens are silently dropped.

## Expose a Reply Tool (Two-Way)

For bidirectional channels (chat bridges), expose a standard MCP tool:

```ts
import { ListToolsRequestSchema, CallToolRequestSchema } from '@modelcontextprotocol/sdk/types.js'

// Tool discovery
mcp.setRequestHandler(ListToolsRequestSchema, async () => ({
  tools: [{
    name: 'reply',
    description: 'Send a message back over this channel',
    inputSchema: {
      type: 'object',
      properties: {
        chat_id: { type: 'string', description: 'The conversation to reply in' },
        text: { type: 'string', description: 'The message to send' },
      },
      required: ['chat_id', 'text'],
    },
  }],
}))

// Tool handler
mcp.setRequestHandler(CallToolRequestSchema, async req => {
  if (req.params.name === 'reply') {
    const { chat_id, text } = req.params.arguments as { chat_id: string; text: string }
    await sendToYourPlatform(chat_id, text)
    return { content: [{ type: 'text', text: 'sent' }] }
  }
  throw new Error(`unknown tool: ${req.params.name}`)
})
```

Update `instructions` so Claude knows to use the tool:

```ts
instructions: 'Messages arrive as <channel source="..." chat_id="...">. Reply with the reply tool, passing the chat_id from the tag.'
```

## Gate Inbound Messages (Security)

An ungated channel is a prompt injection vector. Always check the sender against an allowlist **before** emitting:

```ts
const allowed = new Set(loadAllowlist())

// In your message handler:
if (!allowed.has(message.from.id)) {
  return  // drop silently
}
await mcp.notification({ ... })
```

**Important:** Gate on the sender's ID (`message.from.id`), not the chat/room (`message.chat.id`). In group chats, anyone in the room could inject messages.

## Permission Relay

> Requires Claude Code v2.1.81+

Bidirectional channels can forward permission prompts for remote approval.

### How It Works

1. Claude Code generates a request ID and notifies your server
2. Your server forwards the prompt to the chat platform
3. The remote user replies with `yes <id>` or `no <id>`
4. Your handler parses the reply and emits a verdict

The local terminal dialog stays open — the first response (local or remote) is applied.

### Implementation

**1. Declare the capability:**

```ts
capabilities: {
  experimental: {
    'claude/channel': {},
    'claude/channel/permission': {},  // opt in to relay
  },
  tools: {},
},
```

**2. Handle incoming request:**

```ts
import { z } from 'zod'

const PermissionRequestSchema = z.object({
  method: z.literal('notifications/claude/channel/permission_request'),
  params: z.object({
    request_id: z.string(),
    tool_name: z.string(),
    description: z.string(),
    input_preview: z.string(),
  }),
})

mcp.setNotificationHandler(PermissionRequestSchema, async ({ params }) => {
  send(
    `Claude wants to run ${params.tool_name}: ${params.description}\n\n` +
    `Reply "yes ${params.request_id}" or "no ${params.request_id}"`,
  )
})
```

**3. Intercept verdict in inbound handler:**

```ts
const PERMISSION_REPLY_RE = /^\s*(y|yes|n|no)\s+([a-km-z]{5})\s*$/i

// In your inbound message handler:
const m = PERMISSION_REPLY_RE.exec(message.text)
if (m) {
  await mcp.notification({
    method: 'notifications/claude/channel/permission',
    params: {
      request_id: m[2].toLowerCase(),
      behavior: m[1].toLowerCase().startsWith('y') ? 'allow' : 'deny',
    },
  })
  return  // don't forward as chat
}
```

### Permission Request Fields

| Field | Description |
|-------|-------------|
| `request_id` | 5 lowercase letters (no `l`). Include in the prompt to be echoed in the reply |
| `tool_name` | Tool name (e.g., `Bash`, `Write`) |
| `description` | Human-readable summary of what the tool call does |
| `input_preview` | Tool arguments as JSON, truncated to 200 chars |

## Package as a Plugin

To make your channel installable and shareable:

1. Package it as a [plugin](https://code.claude.com/docs/en/plugins)
2. Publish to a [marketplace](https://code.claude.com/docs/en/plugin-marketplaces)
3. Users install with `/plugin install` and enable with `--channels plugin:<name>@<marketplace>`

Channels published to custom marketplaces still require `--dangerously-load-development-channels` until added to the official allowlist or the organization's `allowedChannelPlugins`.

## Testing During the Research Preview

Custom channels are not on the approved allowlist. Use the development flag:

```bash
# Plugin in development
claude --dangerously-load-development-channels plugin:yourplugin@yourmarketplace

# Bare .mcp.json server
claude --dangerously-load-development-channels server:webhook
```

The bypass is per-entry. The `channelsEnabled` organization policy still applies.

## References

- [Channels Guide](channels.md) — setup for official channels (Telegram, Discord, iMessage)
- [Working implementations](https://github.com/anthropics/claude-plugins-official/tree/main/external_plugins) — complete server code with pairing, reply tools, and file attachments
- [MCP](https://modelcontextprotocol.io) — the underlying protocol that channel servers implement
