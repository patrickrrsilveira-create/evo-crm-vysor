import { useState, useEffect, useCallback } from 'react';
import { agentIntegrationsService } from '@/services/agents/agentIntegrationsService';

interface ElevenLabsConfig {
  provider?: string;
  connected?: boolean;
  respondInAudio?: 'when_client_asks' | 'always' | 'never';
  voice?: string;
  stability?: number;
  similarity?: number;
}

interface GoogleCalendarConfig {
  provider?: string;
  connected?: boolean;
  calendar_id?: string;
}

interface GoogleSheetsConfig {
  provider?: string;
  connected?: boolean;
  spreadsheet_id?: string;
}

export interface KnowledgeNexusConfig {
  provider?: string;
  connected?: boolean;
  nexus_base_url?: string;
  // Sanitized away when reading back from the backend, but the dialog populates
  // it on save when the user enters a new key. Optional so a re-edit doesn't
  // need to retype the key (omit the field to keep the stored value).
  nexus_api_key?: string;
  space_id?: string;
  default_top_k?: number;
  default_filters?: Record<string, unknown>;
  timeout_seconds?: number;
}

interface UseIntegrationsReturn {
  // Configs
  elevenLabsConfig: ElevenLabsConfig | null;
  googleCalendarConfig: GoogleCalendarConfig | null;
  googleSheetsConfig: GoogleSheetsConfig | null;
  knowledgeNexusConfig: KnowledgeNexusConfig | null;
  microsoftTeamsConfig: any | null;

  // Status
  credentialsConfigured: Record<string, boolean>;
  isCheckingIntegrations: boolean;

  // Actions
  reloadConfigs: () => Promise<void>;
  isConnected: (integrationId: string) => boolean;
}

/**
 * Sanitize integration config to remove sensitive fields.
 * Security: This is a defense-in-depth measure to prevent sensitive data
 * from being stored in frontend state, even if backend accidentally sends it.
 */
function sanitizeConfig(config: Record<string, unknown>): Record<string, unknown> {
  if (!config) return config;

  const sanitized = { ...config };

  // Remove sensitive fields
  const sensitiveFields = [
    'apiKey',
    'api_key',
    'nexus_api_key',
    'access_token',
    'client_secret',
    'refresh_token',
    'token',
    'code_verifier',
    'pkce_verifiers',
  ];

  sensitiveFields.forEach(field => {
    delete sanitized[field];
  });

  // Remove any token-like values (REST API keys: sk_, rk_, pk_)
  Object.keys(sanitized).forEach(key => {
    const value = sanitized[key];
    if (typeof value === 'string' && value.match(/^(sk_|rk_|pk_)/)) {
      delete sanitized[key];
    }
  });

  return sanitized;
}

export function useIntegrations(agentId: string): UseIntegrationsReturn {
  const [elevenLabsConfig, setElevenLabsConfig] = useState<ElevenLabsConfig | null>(null);
  const [googleCalendarConfig, setGoogleCalendarConfig] = useState<GoogleCalendarConfig | null>(null);
  const [googleSheetsConfig, setGoogleSheetsConfig] = useState<GoogleSheetsConfig | null>(null);
  const [knowledgeNexusConfig, setKnowledgeNexusConfig] = useState<KnowledgeNexusConfig | null>(null);
  const [microsoftTeamsConfig, setMicrosoftTeamsConfig] = useState<any | null>(null);

  const [isCheckingIntegrations, setIsCheckingIntegrations] = useState(true);
  const [credentialsConfigured, setCredentialsConfigured] = useState<Record<string, boolean>>({
    elevenlabs: false,
    tts: false,
    'google-calendar': false,
    'google-sheets': false,
    'knowledge-nexus': false,
    'microsoft-teams': true,
  });

  const loadConfigs = useCallback(async () => {
    if (!agentId) return;

    setIsCheckingIntegrations(true);

    try {
      // Backend returns an array of { provider, config, ... } items.
      const items = await agentIntegrationsService.getAgentIntegrations(agentId);

      // Build a provider→config map normalized to hyphen-case keys
      // (backend stores underscored provider names like "google_calendar").
      const configsByProvider: Record<string, Record<string, unknown>> = {};
      const credentialsConfiguredNext: Record<string, boolean> = {
        elevenlabs: false,
        'google-calendar': false,
        'google-sheets': false,
        'knowledge-nexus': false,
        'microsoft-teams': true,
      };

      items.forEach(item => {
        const key = (item.provider || '').replace(/_/g, '-');
        if (!key) return;
        configsByProvider[key] = item.config || {};
        credentialsConfiguredNext[key] = true;
      });

      // Ensure microsoft-teams is always true
      credentialsConfiguredNext['microsoft-teams'] = true;

      setCredentialsConfigured(credentialsConfiguredNext);

      // Sanitize configs before storing (defense-in-depth security measure)
      // TTS can be stored under either 'tts' or 'elevenlabs' provider
      setElevenLabsConfig(
        (configsByProvider.tts || configsByProvider.elevenlabs)
          ? (sanitizeConfig(configsByProvider.tts || configsByProvider.elevenlabs) as unknown as ElevenLabsConfig)
          : null
      );
      const calendarConfig = (configsByProvider['google-calendar'] || {}) as Record<string, unknown>;

      setGoogleCalendarConfig(
        (Object.keys(calendarConfig).length > 0 || calendarConfig.connected)
          ? (sanitizeConfig(calendarConfig) as unknown as GoogleCalendarConfig)
          : null
      );
      const sheetsConfig = (configsByProvider['google-sheets'] || {}) as Record<string, unknown>;

      setGoogleSheetsConfig(
        (Object.keys(sheetsConfig).length > 0 || sheetsConfig.connected)
          ? (sanitizeConfig(sheetsConfig) as unknown as GoogleSheetsConfig)
          : null
      );
      setKnowledgeNexusConfig(
        configsByProvider['knowledge-nexus']
          ? (sanitizeConfig(
              configsByProvider['knowledge-nexus']
            ) as unknown as KnowledgeNexusConfig)
          : null
      );
      const teamsConfig = (configsByProvider['microsoft-teams'] || {}) as Record<string, unknown>;
      setMicrosoftTeamsConfig(
        (Object.keys(teamsConfig).length > 0 || teamsConfig.connected)
          ? (sanitizeConfig(teamsConfig) as unknown as any)
          : null
      );
    } catch (error) {
      console.error('Error loading integrations:', error);
      // Reset both credentials flags and per-integration configs so the UI
      // never renders stale "connected" state after a network failure.
      setElevenLabsConfig(null);
      setGoogleCalendarConfig(null);
      setGoogleSheetsConfig(null);
      setKnowledgeNexusConfig(null);
      setMicrosoftTeamsConfig(null);
      setCredentialsConfigured({
        elevenlabs: false,
        tts: false,
        'google-calendar': false,
        'google-sheets': false,
        'knowledge-nexus': false,
        'microsoft-teams': false,
      });
    } finally {
      setIsCheckingIntegrations(false);
    }
  }, [agentId]);

  useEffect(() => {
    loadConfigs();
  }, [loadConfigs]);

  const isConnected = useCallback(
    (integrationId: string): boolean => {
      const configMap: Record<
        string,
        any
      > = {
        elevenlabs: elevenLabsConfig,
        tts: elevenLabsConfig,
        'google-calendar': googleCalendarConfig,
        'google-sheets': googleSheetsConfig,
        'knowledge-nexus': knowledgeNexusConfig,
        'microsoft-teams': microsoftTeamsConfig,
      };

      const config = configMap[integrationId];
      return config?.connected === true;
    },
    [elevenLabsConfig, googleCalendarConfig, googleSheetsConfig, knowledgeNexusConfig, microsoftTeamsConfig]
  );

  return {
    elevenLabsConfig,
    googleCalendarConfig,
    googleSheetsConfig,
    knowledgeNexusConfig,
    microsoftTeamsConfig,
    credentialsConfigured,
    isCheckingIntegrations,
    reloadConfigs: loadConfigs,
    isConnected,
  };
}
