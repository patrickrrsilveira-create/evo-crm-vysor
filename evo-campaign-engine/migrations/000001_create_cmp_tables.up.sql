BEGIN;

-- Throttle profiles: reusable pacing presets per channel type
CREATE TABLE IF NOT EXISTS cmp_throttle_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(120) NOT NULL,
    channel_type    VARCHAR(30)  NOT NULL,

    daily_cap_start     INT     NOT NULL DEFAULT 20,
    daily_cap_max       INT     NOT NULL DEFAULT 500,
    warmup_multiplier   NUMERIC NOT NULL DEFAULT 1.3,
    warmup_step_days    INT     NOT NULL DEFAULT 1,
    hourly_cap          INT     NOT NULL DEFAULT 60,

    min_delay_sec           INT NOT NULL DEFAULT 45,
    max_delay_sec           INT NOT NULL DEFAULT 120,
    coffee_break_every_n    INT NOT NULL DEFAULT 25,
    coffee_break_min_sec    INT NOT NULL DEFAULT 120,
    coffee_break_max_sec    INT NOT NULL DEFAULT 300,

    quiet_hours_start   TIME    NOT NULL DEFAULT '20:00',
    quiet_hours_end     TIME    NOT NULL DEFAULT '08:00',
    respect_timezone    BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Sender instances: each inbox allocated for campaigns
CREATE TABLE IF NOT EXISTS cmp_sender_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inbox_id        INT          NOT NULL,
    channel_type    VARCHAR(30)  NOT NULL,
    identifier      VARCHAR(120) NOT NULL,

    warmup_started_at   TIMESTAMPTZ,
    warmup_stage_day    INT     NOT NULL DEFAULT 0,
    health_score        INT     NOT NULL DEFAULT 100,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    last_sent_at        TIMESTAMPTZ,
    sent_today          INT     NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_cmp_sender_instances_status ON cmp_sender_instances(status);

-- Campaigns
CREATE TABLE IF NOT EXISTS cmp_campaigns (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id          INT          NOT NULL,
    name                VARCHAR(255) NOT NULL,
    trigger_type        VARCHAR(20)  NOT NULL DEFAULT 'manual',
    trigger_config      JSONB        NOT NULL DEFAULT '{}',
    throttle_profile_id UUID         REFERENCES cmp_throttle_profiles(id),
    status              VARCHAR(20)  NOT NULL DEFAULT 'draft',

    total_recipients    INT NOT NULL DEFAULT 0,
    sent_count          INT NOT NULL DEFAULT 0,
    delivered_count     INT NOT NULL DEFAULT 0,
    failed_count        INT NOT NULL DEFAULT 0,

    created_by  INT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at  TIMESTAMPTZ,
    finished_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_cmp_campaigns_account ON cmp_campaigns(account_id);
CREATE INDEX IF NOT EXISTS idx_cmp_campaigns_status  ON cmp_campaigns(status);

-- Campaign channels (pool): N inboxes per campaign for rotation
CREATE TABLE IF NOT EXISTS cmp_campaign_channels (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES cmp_campaigns(id) ON DELETE CASCADE,
    inbox_id    INT  NOT NULL,
    UNIQUE(campaign_id, inbox_id)
);

-- Campaign audience
CREATE TABLE IF NOT EXISTS cmp_campaign_audience (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id     UUID        NOT NULL REFERENCES cmp_campaigns(id) ON DELETE CASCADE,
    mode            VARCHAR(20) NOT NULL DEFAULT 'static',
    segment_filter  JSONB       NOT NULL DEFAULT '{}',
    total_contacts  INT         NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Audience members: snapshot of contacts resolved at dispatch time
CREATE TABLE IF NOT EXISTS cmp_audience_members (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID        NOT NULL REFERENCES cmp_campaigns(id) ON DELETE CASCADE,
    contact_id  INT         NOT NULL,
    recipient   VARCHAR(120) NOT NULL,
    timezone    VARCHAR(60),
    state       VARCHAR(20) NOT NULL DEFAULT 'pending'
);
CREATE INDEX IF NOT EXISTS idx_cmp_audience_members_campaign ON cmp_audience_members(campaign_id, state);

-- Message variants (spintax / A-B testing)
CREATE TABLE IF NOT EXISTS cmp_message_variants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES cmp_campaigns(id) ON DELETE CASCADE,
    subject     TEXT,
    body        TEXT NOT NULL,
    media_url   TEXT,
    media_type  VARCHAR(30),
    weight      INT  NOT NULL DEFAULT 1,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Send jobs: one per contact per campaign (unit of work)
CREATE TABLE IF NOT EXISTS cmp_send_jobs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id         UUID        NOT NULL REFERENCES cmp_campaigns(id) ON DELETE CASCADE,
    contact_id          INT         NOT NULL,
    sender_instance_id  UUID        REFERENCES cmp_sender_instances(id),
    recipient           VARCHAR(120) NOT NULL,
    variant_id          UUID        REFERENCES cmp_message_variants(id),
    state               VARCHAR(20) NOT NULL DEFAULT 'queued',

    scheduled_at    TIMESTAMPTZ,
    sent_at         TIMESTAMPTZ,
    delivered_at    TIMESTAMPTZ,
    read_at         TIMESTAMPTZ,

    attempts    INT  NOT NULL DEFAULT 0,
    last_error  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_cmp_send_jobs_campaign_state ON cmp_send_jobs(campaign_id, state);
CREATE INDEX IF NOT EXISTS idx_cmp_send_jobs_scheduled      ON cmp_send_jobs(scheduled_at) WHERE state IN ('queued','scheduled');

-- Suppression list (opt-out, bounce, manual block)
CREATE TABLE IF NOT EXISTS cmp_suppression (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id   INT          NOT NULL,
    contact_id   INT          NOT NULL,
    channel_type VARCHAR(30)  NOT NULL,
    reason       VARCHAR(60)  NOT NULL DEFAULT 'opt_out',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    UNIQUE(account_id, contact_id, channel_type)
);

-- Delivery events: raw webhook feedback for metrics + ban detection
CREATE TABLE IF NOT EXISTS cmp_delivery_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    send_job_id     UUID        NOT NULL REFERENCES cmp_send_jobs(id) ON DELETE CASCADE,
    status          VARCHAR(30) NOT NULL,
    provider_status VARCHAR(60),
    ts              TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_cmp_delivery_events_job ON cmp_delivery_events(send_job_id);

COMMIT;
