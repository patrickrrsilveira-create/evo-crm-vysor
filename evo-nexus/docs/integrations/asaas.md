# Asaas Integration

Asaas is a Brazilian payment platform that handles billing via Pix, boleto, credit card, and subscriptions, plus marketplace split payments and webhook events. The @flux agent uses Asaas alongside Stripe for a complete picture of incoming revenue from Brazilian customers.

The `int-asaas` skill talks to **Asaas API v3** via a static API key — no OAuth, no token rotation. Setup is a single `.env` edit.

## Setup

### 1. Get Your API Key

1. Log in at [asaas.com](https://www.asaas.com)
2. Go to **Minha Conta > Integrações > API**
3. Generate (or copy) your API key
4. Asaas has separate keys for **sandbox** and **production** — start with sandbox until you are confident the integration is behaving correctly

### 2. Configure .env

```env
ASAAS_API_KEY=your_api_key_here
ASAAS_SANDBOX=true
```

**`ASAAS_SANDBOX=true` is the safe default.** It makes all calls hit `https://sandbox.asaas.com/api/v3` instead of production. Set it to `false` only when you are ready to move charges, customers, and subscriptions into your real account.

### 3. Test the Connection

Use the skill directly through curl (the simplest test):

```bash
curl -H "access_token: $ASAAS_API_KEY" \
  https://sandbox.asaas.com/api/v3/customers?limit=1
```

If you get a JSON list (possibly empty), you are connected. Switch to `https://api.asaas.com/v3` for production.

## How Authentication Works

Asaas uses a static API key sent as a custom HTTP header on every request:

```
access_token: <your ASAAS_API_KEY>
```

No OAuth, no refresh tokens, no expiration. The key stays valid until you rotate it manually in the Asaas dashboard. If you rotate, update `.env` and you are done.

The base URL is chosen dynamically based on `ASAAS_SANDBOX`:

- `ASAAS_SANDBOX=true` → `https://sandbox.asaas.com/api/v3` (default)
- `ASAAS_SANDBOX=false` → `https://api.asaas.com/v3`

## Available Operations

The `int-asaas` skill documents **15 operations** grouped by domain:

| Category | Operations |
|---|---|
| **Payments** | create, get by ID, list (with filters), get Pix QR code, get boleto PDF + barcode |
| **Customers** | create (with CPF/CNPJ validation), list |
| **Subscriptions** | create, list, cancel (recurring billing) |
| **Financial** | get account balance, create transfer (Pix out or TED) |
| **Marketplace** | create subaccount (for split payments) |
| **Utilities** | get installment details, list webhook events |

### Enums You Will Use Often

**`billingType`** (how the customer pays):
- `BOLETO` — Brazilian bank slip
- `CREDIT_CARD` — credit card
- `PIX` — instant transfer
- `UNDEFINED` — let the customer choose at checkout

**Payment `status`** (lifecycle):
- `PENDING` — awaiting payment
- `RECEIVED` — paid
- `CONFIRMED` — paid and reconciled
- `OVERDUE` — past due date
- `REFUNDED` — refunded
- `RECEIVED_IN_CASH` — marked as paid in cash
- `REFUND_REQUESTED`, `CHARGEBACK_REQUESTED`, `AWAITING_CHARGEBACK_REVERSAL`
- `DUNNING_REQUESTED`, `DUNNING_RECEIVED` — collections
- `AWAITING_RISK_ANALYSIS` — under fraud review

### Data Formats

All operations expect Brazilian document formats stripped of punctuation:

| Field | Format |
|---|---|
| CPF | 11 digits, no dashes (e.g. `12345678900`) |
| CNPJ | 14 digits, no dashes (e.g. `12345678000100`) |
| CEP | 8 digits, no dashes (e.g. `01310100`) |
| Dates | ISO `YYYY-MM-DD` |
| Values | Decimal with dot (e.g. `99.90`) |

## Example Calls

```bash
# Create a customer
curl -X POST https://sandbox.asaas.com/api/v3/customers \
  -H "access_token: $ASAAS_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name":"João Silva","cpfCnpj":"12345678900","email":"joao@example.com"}'

# Create a Pix payment
curl -X POST https://sandbox.asaas.com/api/v3/payments \
  -H "access_token: $ASAAS_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"customer":"cus_XXX","billingType":"PIX","value":99.90,"dueDate":"2026-04-15","description":"Plan subscription"}'

# Get the Pix QR code for a payment
curl https://sandbox.asaas.com/api/v3/payments/pay_XXX/pixQrCode \
  -H "access_token: $ASAAS_API_KEY"

# List subscriptions
curl https://sandbox.asaas.com/api/v3/subscriptions?limit=20 \
  -H "access_token: $ASAAS_API_KEY"

# Get account balance
curl https://sandbox.asaas.com/api/v3/finance/balance \
  -H "access_token: $ASAAS_API_KEY"
```

## Webhooks

Asaas can send real-time events for payment status changes (created, received, overdue, refunded, chargeback) and subscription lifecycle. Configure the webhook URL in the Asaas dashboard under **Integrações > Webhooks** — this part is done in the dashboard, not via the API.

The `get_webhook_events` operation is available to list the last events received by Asaas (useful for debugging delivery failures).

## Advanced Endpoints (Reference Only)

The following are documented by the Asaas API but **not** wrapped by the skill yet:

- Credit card tokenization (requires extra fields and PCI considerations)
- Receivables anticipation (`antecipações`)
- Tax/fiscal info attached to payments
- Batch payment creation

If you need any of these, extend the skill using the reference implementation at `workspace/projects/mcp-dev-brasil/packages/payments/asaas/` and the [official docs](https://docs.asaas.com).

## Skills That Use Asaas

| Skill | What it does |
|---|---|
| `int-asaas` | Direct Asaas API queries — payments, customers, subscriptions, balance, transfers |
| `fin-daily-pulse` | Daily snapshot — can include Asaas receipts alongside Stripe |
| `fin-weekly-report` | Weekly consolidation across Stripe + Asaas + Omie + Bling |
| `fin-monthly-close-kickoff` | Month-end close — Asaas subscriptions and payments feed into revenue reconciliation |

## Troubleshooting

### 401 `invalid_token`
Your `ASAAS_API_KEY` is wrong, expired (manually rotated), or you are mixing sandbox and production keys. Re-copy the key from the Asaas dashboard, making sure you are on the correct environment, and confirm `ASAAS_SANDBOX` matches.

### 400 `invalid cpfCnpj`
You passed a CPF or CNPJ with dots, dashes, or slashes. Asaas expects digits only — 11 for CPF, 14 for CNPJ.

### Sandbox calls returning production data (or vice versa)
Double-check `ASAAS_SANDBOX` in `.env`. The two environments are completely isolated — a customer created in sandbox does **not** exist in production.

### Webhook is not firing
Webhooks are configured in the Asaas dashboard, not in `.env`. Go to **Integrações > Webhooks**, make sure the URL is reachable from the public internet, and use `get_webhook_events` to inspect delivery attempts.
