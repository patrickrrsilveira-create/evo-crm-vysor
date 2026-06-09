# Omie ERP Integration

Omie is a Brazilian cloud ERP used for managing clients, products, sales orders, invoices (NF-e), financials (accounts receivable/payable), and stock. The @flux agent uses Omie data alongside Stripe for a complete financial picture.

## Setup

### 1. Get Your API Credentials

1. Go to [omie.com.br](https://omie.com.br)
2. Navigate to **Settings > Integrations > API**
3. Copy your **App Key** and **App Secret**

### 2. Configure .env

```env
OMIE_APP_KEY=your_app_key_here
OMIE_APP_SECRET=your_app_secret_here
```

### 3. Test the Connection

```bash
make fin-pulse
```

This runs the daily financial pulse, which queries both Stripe and Omie. If Omie credentials are valid, you will see accounts payable/receivable data in the report.

## Available Commands

The `int-omie` skill provides these operations via a Python client:

| Category | Commands |
|---|---|
| **Clients** | `clientes_listar`, `clientes_buscar` (by CNPJ/CPF or code), `clientes_detalhar` |
| **Products** | `produtos_listar`, `produtos_detalhar` |
| **Sales Orders** | `pedidos_listar`, `pedidos_detalhar`, `pedidos_status` |
| **Financials** | `contas_receber`, `contas_pagar`, `resumo_financeiro` |
| **Invoices (NF-e)** | `nfe_listar`, `nfe_detalhar` |
| **Stock** | `estoque_posicao`, `estoque_produto` |

All commands support pagination with `[pagina] [por_pagina]` arguments.

## Webhook Support

Omie can send real-time events via webhooks. Supported events include order creation/update, invoice issuance/cancellation, payment received/sent, and client creation/update.

## Skills That Use Omie

| Skill | What it does |
|---|---|
| `int-omie` | Direct Omie API queries -- clients, orders, invoices, financials, stock |
| `fin-daily-pulse` | Daily snapshot combining Stripe revenue + Omie expenses |
| `fin-monthly-close-kickoff` | Month-end close with P&L, reconciliations, action items |
| `fin-financial-statements` | Full financial statements with variance analysis |

## Automated Routines

| Routine | Schedule | Make command |
|---|---|---|
| Financial Pulse | 19:00 BRT daily | `make fin-pulse` |
| Financial Weekly | Friday 07:30 BRT | `make fin-weekly` |
| Monthly Close Kickoff | 1st of month | `make fin-close` |
