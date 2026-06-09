# LinkedIn Integration

LinkedIn API provides profile information and organization statistics. The @pixel agent uses it for social media analytics. Supports multiple accounts via OAuth.

## Setup

### 1. Create a LinkedIn App

1. Go to [linkedin.com/developers](https://www.linkedin.com/developers/)
2. Create a new app
3. Navigate to **Auth** tab
4. Copy your **Client ID** and **Client Secret**
5. Under **Authorized redirect URLs**, add:
   - `http://localhost:8765/callback/linkedin`
   - If using ngrok, also add your ngrok URL: `https://your-id.ngrok.io/callback/linkedin`

### 2. Configure .env

```env
LINKEDIN_CLIENT_ID=your_client_id_here
LINKEDIN_CLIENT_SECRET=your_client_secret_here
```

### 3. Connect Your Account

```bash
make social-auth
```

Open the dashboard Integrations page and click **Connect** on LinkedIn. Complete the OAuth flow.

### 4. Test the Connection

```bash
make social
```

## Available Commands

The `int-linkedin` skill queries the LinkedIn API:

| Command | What it does |
|---|---|
| `accounts` | List configured LinkedIn accounts |
| `profile [account]` | Name, email, photo |
| `my_posts [account] [N]` | Recent posts (requires `w_member_social` scope) |
| `post_stats POST_URN` | Reactions and comments on a post |
| `org_followers [account]` | Organization followers (requires Advertising API) |
| `summary` | Summary across all accounts |

## Available Scopes

| Scope | Status | LinkedIn Product |
|---|---|---|
| `openid profile email` | Active | Sign In with OpenID Connect |
| `w_member_social` | Active | Share on LinkedIn |
| `r_organization_social` | Pending | Advertising API (requires approval) |
| `r_organization_admin` | Pending | Advertising API (requires approval) |

## Current Limitations

- **Posts:** Reading posts requires the `w_member_social` scope
- **Company Page analytics:** Requires Advertising API approval (pending)
- **Workaround:** Export CSV from LinkedIn Analytics for Company Page data

## Skills That Use LinkedIn

| Skill | What it does |
|---|---|
| `int-linkedin` | Direct LinkedIn API queries |
| `social-linkedin-report` | LinkedIn analytics HTML report |
| `social-analytics-report` | Cross-platform analytics (includes LinkedIn) |

## Automated Routines

| Routine | Schedule | Make command |
|---|---|---|
| Social Analytics | 18:00 BRT daily | `make social` |
| Social Weekly | Friday 08:15 BRT | `make social` |
| Social Monthly | 1st of month | `make social` |
