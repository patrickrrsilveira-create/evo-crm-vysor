# N8N Webhook Update — Campaign Engine Integration

## Overview

The n8n webhook for WhatsApp video dispatch needs to accept an optional `instance` field from the campaign engine to support multi-number rotation. This change is **retrocompatible** — existing payloads without `instance` will continue to work.

## Current Payload (Legacy)

```json
{
  "telefone": "+5511999999999",
  "video_url": "https://minio-api.vysortech.app.br/files/video.mp4",
  "media": "https://minio-api.vysortech.app.br/files/video.mp4"
}
```

## New Payload (Campaign Engine)

```json
{
  "telefone": "+5511999999999",
  "video_url": "https://minio-api.vysortech.app.br/files/video.mp4",
  "media": "https://minio-api.vysortech.app.br/files/video.mp4",
  "instance": "sender-instance-id"  // NEW: identifies which WhatsApp number sent this
}
```

## Implementation Steps in N8N

### 1. Receive Webhook Trigger
Your existing **Webhook** node remains unchanged — it already receives the full payload.

### 2. Extract Instance (Optional Routing)

Add a **JavaScript** or **Set** node to handle the instance:

**Option A: If using evolution-api instances**
```javascript
// Extract instance to route to correct number
const instance = $input.first().json.instance || 'default';
const telefone = $input.first().json.telefone;
const videoUrl = $input.first().json.video_url;

return {
  instance,
  telefone,
  videoUrl,
  // ... pass to your evolution-api endpoint
};
```

**Option B: If using multiple WhatsApp numbers**
```javascript
const instance = $input.first().json.instance;

// Map instance ID to WhatsApp API credential
const instanceMap = {
  'sender-a': 'whatsapp_cred_1',
  'sender-b': 'whatsapp_cred_2',
  'default': 'whatsapp_cred_default'
};

const credential = instanceMap[instance] || instanceMap['default'];
return {
  credential,
  telefone: $input.first().json.telefone,
  videoUrl: $input.first().json.video_url
};
```

### 3. Forward to WhatsApp Provider

Modify your existing **HTTP Request** or **evolution-api** node:

**If using HTTP to evolution-api:**
```
POST /message/sendMedia/{instance}
Headers:
  - Content-Type: application/json
Body:
{
  "number": {{ $input.first().json.telefone }},
  "mediaUrl": {{ $input.first().json.video_url }},
  "mediaType": "video/mp4"
}
```

**If using n8n's WhatsApp node:**
- Pass the `instance` as a parameter to select the correct WhatsApp connection
- Keep existing logic for fallback to default instance if `instance` is empty

### 4. Error Handling (Same as Before)

```javascript
// If video send fails, return error
if (error) {
  return {
    success: false,
    error: error.message,
    telefone: $input.first().json.telefone,
    instance: $input.first().json.instance || 'default'
  };
}
```

## Testing Checklist

- [ ] Legacy payload (no `instance`) still works → sends from default number
- [ ] New payload with `instance` → sends from specified number
- [ ] Invalid instance ID → fallback to default, not error
- [ ] Video URL validation (must be HTTPS from MinIO)
- [ ] Phone number validation (must have country code)
- [ ] Error responses logged for monitoring

## Rollback Plan

If issues arise, revert to ignoring the `instance` field:
- Remove instance-based routing
- Always use default WhatsApp number
- Existing single-number flow continues to work

## Questions?

Refer to campaign engine specs: `docs/campaign-engine.md` → Section 8 (n8n contract changes).
