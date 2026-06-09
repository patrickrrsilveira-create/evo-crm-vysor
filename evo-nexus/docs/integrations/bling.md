# Bling ERP Integration

Bling is a Brazilian cloud ERP for managing products, sales orders, contacts (customers/suppliers), fiscal invoices (NF-e), and stock. The @flux agent uses Bling alongside Omie, Stripe, and Asaas to get a complete financial and operational picture for Brazilian businesses.

The `int-bling` skill talks to **Bling API v3** via OAuth2 with **automatic token refresh** — you authenticate once with `make bling-auth` and the skill keeps the session alive forever.

## Setup

### 1. Create a Bling App

1. Log in at [developer.bling.com.br](https://developer.bling.com.br)
2. Go to **Meus Apps** and create a new application
3. Set the **redirect URI** to exactly:
   ```
   http://localhost:8787/callback
   ```
4. Copy the **Client ID** and **Client Secret**

### 2. Configure .env

```env
BLING_CLIENT_ID=your_client_id_here
BLING_CLIENT_SECRET=your_client_secret_here
```

Leave `BLING_ACCESS_TOKEN` and `BLING_REFRESH_TOKEN` empty — the next step fills them in.

### 3. Run the OAuth Login

```bash
make bling-auth
```

What it does:

1. Starts a local callback server on `http://localhost:8787/callback`
2. Opens your browser at Bling's authorization page
3. You log in and approve the app
4. Bling redirects back to the local server with an authorization code
5. The helper exchanges the code for `access_token` + `refresh_token`
6. Both tokens are persisted to `.env`

After a successful run, `.env` contains:

```env
BLING_CLIENT_ID=...
BLING_CLIENT_SECRET=...
BLING_ACCESS_TOKEN=...
BLING_REFRESH_TOKEN=...
```

You only do this once per Bling account. From now on, the skill refreshes tokens automatically whenever they expire.

### 4. Test the Connection

Run any read operation through the Python client:

```bash
python3 .claude/skills/int-bling/scripts/bling_client.py GET /produtos --params page=1 limit=10
```

If the JSON response includes products, you are connected. If you get a `401`, the client will automatically refresh the token and retry once.

## How Automatic Refresh Works

The Python client (`.claude/skills/int-bling/scripts/bling_client.py`) is the single entry point for all API calls. On any `HTTP 401`:

1. It exchanges `BLING_REFRESH_TOKEN` for a new `access_token` + `refresh_token` at `https://www.bling.com.br/Api/v3/oauth/token` using HTTP Basic auth with Client ID and Client Secret.
2. Both new tokens are persisted back to `.env` and to `os.environ`.
3. The original request is retried once with the fresh access token.

> **Important:** Bling **rotates the refresh token on every refresh**, so always use the provided client. Never call the API directly with `curl` unless you are debugging — otherwise the next refresh will use a stale refresh token and the session will break. If that happens, run `make bling-auth` again to start over.

## Calling the Client

```bash
# List products (paginated)
python3 .claude/skills/int-bling/scripts/bling_client.py GET /produtos --params page=1 limit=50

# List orders for a date range
python3 .claude/skills/int-bling/scripts/bling_client.py GET /pedidos/vendas --params dataInicial=2026-01-01 dataFinal=2026-04-10

# Create a contact (PJ)
python3 .claude/skills/int-bling/scripts/bling_client.py POST /contatos --body '{"nome":"Acme Ltda","tipo":"J","numeroDocumento":"12345678000100"}'

# Update a product's price
python3 .claude/skills/int-bling/scripts/bling_client.py PUT /produtos/123 --body '{"preco":99.90}'
```

The client prints the JSON response to stdout and writes errors to stderr.

## Available Operations

The `int-bling` skill documents **10 operations** covering the core of an ERP workflow:

| Category | Operations |
|---|---|
| **Products** | list, create (fields: `nome`, `codigo`, `preco`, `precoCusto`, `tipo` P/S/K, `formato` S/E/V) |
| **Sales Orders** | list, create (fields: `contato.id`, `itens[]` with `produto.id`, `quantidade`, `valor`, `desconto`) |
| **Contacts** | list, create (`F` for people, `J` for companies, with CPF/CNPJ) |
| **Fiscal Invoices (NF-e)** | list, create (issued from an existing order) |
| **Stock** | get stock level, update stock by warehouse (`depositoId`) |

All list endpoints support pagination via `page` and `limit` (default 100). Dates use the `YYYY-MM-DD` format.

## Advanced Endpoints (Reference Only)

The following are **not** wrapped by the skill yet — use the reference implementation in `workspace/projects/mcp-dev-brasil/packages/erp/bling/` and the [official API v3 docs](https://developer.bling.com.br) as guides:

- Product categories and variations
- Purchase orders and production management
- Warehouse/depósito CRUD
- Stock movement history and audit trails
- NF-e cancellations and complementary invoices
- Webhooks

If you need any of these, extend `bling_client.py` and document the new operation in `.claude/skills/int-bling/SKILL.md`.

## Skills That Use Bling

| Skill | What it does |
|---|---|
| `int-bling` | Direct Bling API queries — products, orders, contacts, NF-e, stock |
| `fin-daily-pulse` | Daily snapshot — can include Bling invoices alongside Stripe/Asaas/Omie |
| `fin-weekly-report` | Weekly consolidation across all connected ERPs and payment gateways |
| `fin-monthly-close-kickoff` | Month-end close — Bling invoices feed into the revenue reconciliation |

## Troubleshooting

### `make bling-auth` says "missing BLING_CLIENT_ID / BLING_CLIENT_SECRET"
Fill both values in `.env` first, then run again.

### Browser does not open automatically
Copy the URL the CLI prints and paste it manually. The local callback server will still catch the redirect.

### Redirect URI mismatch
Bling will show an error page if the app's redirect URI is not exactly `http://localhost:8787/callback`. Fix it in the app settings at `developer.bling.com.br`.

### 401 errors keep coming back
Your refresh token is probably stale because some other process bypassed `bling_client.py`. Run `make bling-auth` again to get a fresh pair of tokens.

### State mismatch error
The CLI refused the callback because of a CSRF check failure (very rare — usually means the browser returned an old callback). Re-run `make bling-auth`.
