# Instagram Integration

Instagram Graph API provides profile statistics, post engagement, reach, and impressions. The @pixel agent uses it for social media analytics. Supports multiple accounts via OAuth through Facebook.

## Setup

### 1. Create a Facebook App

1. Go to [developers.facebook.com](https://developers.facebook.com)
2. Create a new app (type: Business)
3. Navigate to **Settings > Basic**
4. Copy your **App ID** and **App Secret**
5. Go to **Facebook Login > Settings** and add the redirect URI:
   - `http://localhost:8765/callback/instagram`
   - If using ngrok, also add your ngrok URL: `https://your-id.ngrok.io/callback/instagram`

### 2. Configure .env

```env
META_APP_ID=your_app_id_here
META_APP_SECRET=your_app_secret_here
```

### 3. Connect Your Account

```bash
make social-auth
```

Open the dashboard Integrations page and click **Connect** on Instagram. Complete the Facebook OAuth flow. The Instagram Business/Creator account linked to your Facebook Page will be connected.

### 4. Test the Connection

```bash
make social
```

## Available Commands

The `int-instagram` skill queries the Instagram Graph API:

| Command | What it does |
|---|---|
| `accounts` | List configured Instagram accounts |
| `profile [account]` | Followers, bio, media count |
| `recent_posts [account] [N]` | Last N posts with engagement |
| `top_posts [account] [N]` | Top N posts by engagement |
| `post_insights POST_ID` | Detailed insights for a specific post |
| `account_insights [account]` | Impressions, reach, profile views (30 days) |
| `summary` | Summary across all accounts |

## Key Metrics

- Followers (delta via daily snapshots)
- Engagement rate: (likes + comments) / followers
- Reach and impressions
- Profile views
- Best post of the period
- Reels vs static post performance
- Publishing frequency

## Rate Limits

- Platform endpoints: 4800 calls per 24 hours
- Business Discovery / Hashtag: 200 calls per hour per user

## Skills That Use Instagram

| Skill | What it does |
|---|---|
| `int-instagram` | Direct Instagram Graph API queries |
| `social-instagram-report` | Instagram analytics HTML report |
| `social-analytics-report` | Cross-platform analytics (includes Instagram) |

## Automated Routines

| Routine | Schedule | Make command |
|---|---|---|
| Social Analytics | 18:00 BRT daily | `make social` |
| Social Weekly | Friday 08:15 BRT | `make social` |
| Social Monthly | 1st of month | `make social` |
