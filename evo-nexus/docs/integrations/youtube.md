# YouTube Integration

YouTube Data API v3 provides channel statistics, video metrics, and engagement data. The @pixel agent uses it for social media analytics and content performance tracking. Supports multiple accounts via OAuth.

## Setup

### 1. Create OAuth Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Enable the **YouTube Data API v3**
3. Create an **OAuth 2.0 Client ID** (Web application)
4. Under **Authorized redirect URIs**, add: `http://localhost:8765/callback/youtube`

### 2. Configure .env

```env
YOUTUBE_OAUTH_CLIENT_ID=your_client_id_here
YOUTUBE_OAUTH_CLIENT_SECRET=your_client_secret_here
```

### 3. Connect Your Account

```bash
make social-auth
```

Open the dashboard Integrations page and click **Connect** on YouTube. Complete the OAuth flow. Tokens are saved to `.env` automatically.

### 4. Test the Connection

```bash
make social
```

## Available Commands

The `int-youtube` skill queries the YouTube Data API:

| Command | What it does |
|---|---|
| `accounts` | List configured YouTube accounts |
| `channel_stats [account]` | Subscribers, total views, video count |
| `recent_videos [account] [N]` | Last N videos with metrics |
| `top_videos [account] [N]` | Top N videos by views |
| `video_stats VIDEO_ID` | Stats for specific videos |
| `comments VIDEO_ID [N]` | Comments on a video |
| `summary` | Summary across all accounts |

## Key Metrics

- Subscribers (with daily/weekly/monthly delta)
- Total views and per-video views
- Engagement rate: (likes + comments) / views
- Best performing video of the period
- Publishing frequency
- Recent comment sentiment

## Skills That Use YouTube

| Skill | What it does |
|---|---|
| `int-youtube` | Direct YouTube API queries |
| `social-youtube-report` | YouTube analytics HTML report |
| `social-analytics-report` | Cross-platform analytics (includes YouTube) |

## Automated Routines

| Routine | Schedule | Make command |
|---|---|---|
| Social Analytics | 18:00 BRT daily | `make social` |
| Social Weekly | Friday 08:15 BRT | `make social` |
| Social Monthly | 1st of month | `make social` |
