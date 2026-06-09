# Telegram Integration

Telegram provides real-time messaging and notifications via a bot connected through MCP. The @clawdia agent uses it to send alerts, receive commands, and communicate with the user on the go.

> **Channel mode:** Telegram also works as a [Channel](../guides/channels.md#setup-telegram) — a bidirectional chat bridge that pushes messages into your Claude Code session. Start with `make telegram`. See the [Channels Guide](../guides/channels.md) for full setup including pairing and security.

## Setup

### 1. Create a Bot

1. Open Telegram and message [@BotFather](https://t.me/BotFather)
2. Send `/newbot` and follow the prompts
3. Copy the **Bot Token** provided

### 2. Get Your Chat ID

1. Send `/start` to your new bot
2. Visit `https://api.telegram.org/bot<TOKEN>/getUpdates`
3. Find `chat.id` in the response

### 3. Configure .env

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here
```

### 4. Start the Bot

```bash
make telegram
```

This starts the Telegram bot in a background screen session.

## Available Tools

Telegram is connected via MCP plugin. The following tools are available:

| Tool | What it does |
|---|---|
| `reply` | Send a message or reply to a specific message |
| `edit_message` | Edit an already sent message |
| `react` | Add an emoji reaction to a message |
| `download_attachment` | Download photo, file, or audio from a message |

### Sending Messages

```
reply(chat_id="...", text="Your message here")
```

### Sending Attachments

```
reply(chat_id="...", text="Here is the report", files=["/path/to/file.pdf"])
```

### Replying to a Specific Message

```
reply(chat_id="...", text="Your reply", reply_to="message_id")
```

## Key Behaviors

- Messages sent in the Claude session are **not visible** to Telegram -- everything must go through the `reply` tool
- Edits do not trigger push notifications -- send a new message for important updates
- The Bot API has no history or search -- only incoming messages are visible
- Access control is managed via the `/telegram:access` skill

## Skills That Use Telegram

| Skill | What it does |
|---|---|
| `int-telegram` | Send, reply, react, edit messages via MCP |

## Automated Routines

The Telegram bot runs continuously (`make telegram`) and receives messages in real time. Many routines send notification summaries to Telegram as part of their output.
