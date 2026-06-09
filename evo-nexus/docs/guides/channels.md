# Channels — Bidirectional Chat with Claude Code

Channels let you push messages into a running Claude Code session and receive replies back — like a chat bridge between external platforms and your terminal.

> **Channels vs Integrations:** Integrations (API clients, MCP servers) are tools Claude calls on demand during a task. Channels are the reverse: external systems that **push** events into the session. You can have both — for example, Discord as an integration (community monitoring via @pulse) AND as a channel (bidirectional chat with Claude).

## Supported Channels

| Channel | Platform | Type | Make Command |
|---------|----------|------|-------------|
| **Telegram** | Mobile/Desktop | Bot (polling) | `make telegram` |
| **Discord** | Desktop/Mobile | Bot (gateway) | `make discord-channel` |
| **iMessage** | macOS only | Native (Messages.app) | `make imessage` |

## Prerequisites

- **Claude Code v2.1.80+** (check with `claude --version`)
- **Bun** installed (check with `bun --version`, install at [bun.sh](https://bun.sh))
- **claude.ai login** (Console and API key authentication are not supported)
- **Plugin marketplace** configured: `/plugin marketplace add anthropics/claude-plugins-official`

## Setup: Telegram

Full integration docs at [docs/integrations/telegram.md](../integrations/telegram.md).

### 1. Create a Bot

1. Open Telegram and message [@BotFather](https://t.me/BotFather)
2. Send `/newbot` and follow the prompts
3. Copy the **Bot Token** provided

### 2. Install the Plugin

In Claude Code:

```
/plugin install telegram@claude-plugins-official
```

If the plugin is not found, update the marketplace:

```
/plugin marketplace update claude-plugins-official
```

Then reload:

```
/reload-plugins
```

### 3. Configure the Token

```
/telegram:configure <your_token>
```

This saves to `~/.claude/channels/telegram/.env`. Alternatively, set `TELEGRAM_BOT_TOKEN` in your shell environment.

### 4. Start the Channel

**Manual (foreground):**

```bash
claude --channels plugin:telegram@claude-plugins-official
```

**Background (recommended):**

```bash
make telegram
```

This starts the bot in a background `screen` session with `--dangerously-skip-permissions`.

### 5. Pair Your Account

1. In Telegram, send any message to your bot
2. The bot replies with a **pairing code**
3. In Claude Code, run:

```
/telegram:access pair <code>
```

4. Lock down access to your account only:

```
/telegram:access policy allowlist
```

### Manage

| Command | What it does |
|---------|-------------|
| `make telegram` | Start bot in background |
| `make telegram-stop` | Stop the bot |
| `make telegram-attach` | Attach to terminal (Ctrl+A D to detach) |

---

## Setup: Discord Channel

> **Note:** This is the Discord **Channel** (bidirectional chat with Claude Code), separate from the Discord API integration used by @pulse for community monitoring. They are independent and can coexist.

### 1. Create a Discord Bot

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Click **New Application** and give it a name
3. In the **Bot** section, click **Reset Token** and copy the token

### 2. Enable Message Content Intent

In your bot settings, under **Privileged Gateway Intents**, enable:

- **Message Content Intent** (required)

### 3. Invite the Bot to Your Server

Go to **OAuth2 > URL Generator**. Select the `bot` scope and enable these permissions:

- View Channels
- Send Messages
- Send Messages in Threads
- Read Message History
- Attach Files
- Add Reactions

Open the generated URL to add the bot to your server.

### 4. Install the Plugin

In Claude Code:

```
/plugin install discord@claude-plugins-official
```

If not found, update the marketplace:

```
/plugin marketplace update claude-plugins-official
```

Then reload:

```
/reload-plugins
```

### 5. Configure the Token

```
/discord:configure <your_token>
```

Saves to `~/.claude/channels/discord/.env`. Alternatively, set `DISCORD_BOT_TOKEN` in your environment.

### 6. Start the Channel

**Manual (foreground):**

```bash
claude --channels plugin:discord@claude-plugins-official
```

**Background (recommended):**

```bash
make discord-channel
```

### 7. Pair Your Account

1. DM your bot on Discord
2. The bot replies with a **pairing code**
3. In Claude Code:

```
/discord:access pair <code>
```

4. Lock down access:

```
/discord:access policy allowlist
```

### Manage

| Command | What it does |
|---------|-------------|
| `make discord-channel` | Start channel in background |
| `make discord-channel-stop` | Stop the channel |
| `make discord-channel-attach` | Attach to terminal (Ctrl+A D to detach) |

---

## Setup: iMessage

> **Requirement:** macOS only. Reads the Messages.app database directly and sends replies via AppleScript. No bot token or external service needed.

### 1. Grant Full Disk Access

The Messages database (`~/Library/Messages/chat.db`) is protected by macOS. On the first read, macOS prompts for permission — click **Allow**.

If the prompt doesn't appear or you clicked "Don't Allow", grant access manually:

**System Settings > Privacy & Security > Full Disk Access** — add your terminal app (Terminal, iTerm, etc.)

### 2. Install the Plugin

In Claude Code:

```
/plugin install imessage@claude-plugins-official
```

If not found:

```
/plugin marketplace update claude-plugins-official
```

### 3. Start the Channel

**Manual (foreground):**

```bash
claude --channels plugin:imessage@claude-plugins-official
```

**Background (recommended):**

```bash
make imessage
```

### 4. Test — Send a Message to Yourself

Open Messages on any device signed into your Apple ID and send a message **to yourself**. Self-chat works automatically with no access configuration needed.

> On Claude's first reply, macOS prompts for Automation permission to control Messages. Click **OK**.

### 5. Allow Other Senders

By default, only your own messages pass through. To allow another contact:

```
/imessage:access allow +15551234567
```

Accepted handles: phone numbers with country code (`+1...`) or Apple ID emails.

### Manage

| Command | What it does |
|---------|-------------|
| `make imessage` | Start channel in background |
| `make imessage-stop` | Stop the channel |
| `make imessage-attach` | Attach to terminal (Ctrl+A D to detach) |

---

## Multiple Channels

You can run multiple channels simultaneously. Each runs in its own `screen` session:

```bash
make telegram
make discord-channel
make imessage
```

Or in a single Claude Code session (foreground):

```bash
claude --channels plugin:telegram@claude-plugins-official plugin:discord@claude-plugins-official plugin:imessage@claude-plugins-official
```

## Security

- Each channel maintains a **sender allowlist** — only approved IDs can push messages
- Telegram and Discord use **pairing**: the sender messages the bot, receives a code, and you approve it in Claude Code
- iMessage allows self-chat automatically; other contacts are added by handle
- **Permission relay**: bidirectional channels can forward permission prompts so you can approve/deny remotely (e.g., approve a `Bash` command from Telegram)
- Being in `.mcp.json` is not enough — the server must also be named in `--channels`

## Enterprise / Team

On Team and Enterprise plans, channels are **disabled by default**. The admin controls access via:

| Setting | Purpose |
|---------|---------|
| `channelsEnabled` | Master switch. Must be `true` for any channel to work |
| `allowedChannelPlugins` | Which plugins can register as channels (replaces Anthropic's default list) |

Configure at: **claude.ai > Admin settings > Claude Code > Channels**

## Custom Channels

You can build your own channels for systems that don't have a plugin yet (CI/CD, monitoring, webhooks). See the [Channels Reference](channels-reference.md) for the full technical documentation.

## Troubleshooting

**Bot doesn't respond on Telegram/Discord:**
- Make sure Claude Code is running with `--channels` (or via `make telegram` / `make discord-channel`)
- Use `make telegram-attach` or `make discord-channel-attach` to view the terminal

**"Plugin not found":**
- Run `/plugin marketplace update claude-plugins-official`
- Then `/reload-plugins`

**"Blocked by org policy":**
- Admin needs to enable `channelsEnabled` in managed settings

**iMessage — "authorization denied":**
- Grant Full Disk Access to your terminal in System Settings > Privacy & Security

**Messages not arriving:**
- Confirm you completed pairing (`/telegram:access pair`, `/discord:access pair`)
- Confirm the policy is set to allowlist (`/telegram:access policy allowlist`)

**Permission prompts block the session:**
- Channels with permission relay forward the prompts. Reply from chat or the local terminal
- For unattended use: `--dangerously-skip-permissions` (only in trusted environments)
