import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';
import { api } from '@/services/core';
import { setupService } from '@/services/setup/setupService';
import { initClarity } from '@/utils/clarityUtils';

export interface GlobalConfig {
  fbAppId?: string;
  fbApiVersion?: string;
  wpAppId?: string;
  wpApiVersion?: string;
  wpWhatsappConfigId?: string;
  instagramAppId?: string;
  googleOAuthClientId?: string;
  azureAppId?: string;
  // 🔒 SECURITY: Don't expose sensitive API URLs to frontend
  // Only boolean indicators to check if config exists
  hasFacebookConfig?: boolean;
  hasWhatsappConfig?: boolean;
  hasInstagramConfig?: boolean;
  hasEvolutionConfig?: boolean;
  hasEvolutionGoConfig?: boolean;
  hasEvolutionHubConfig?: boolean;
  evolutionHubEnabled?: boolean;
  hasTwitterConfig?: boolean;
  openaiConfigured?: boolean;
  enableAccountSignup?: boolean;
  recaptchaSiteKey?: string;
  clarityProjectId?: string;
  appLogoUrl?: string;
  appPrimaryColor?: string;
  companyName?: string;
  appLogoWidth?: string;
  appLogoHeight?: string;
  sidebarCopyrightText?: string;
  supportLink?: string;
  docsLink?: string;
}

interface GlobalConfigContextValue extends GlobalConfig {
  setupRequired: boolean;
  setupLoading: boolean;
}

const GlobalConfigContext = createContext<GlobalConfigContextValue>({
  setupRequired: false,
  setupLoading: true,
});

// Cache global para evitar múltiplas chamadas
let globalConfigCache: GlobalConfig | null = null;
let globalConfigPromise: Promise<GlobalConfig> | null = null;
let setupRequiredCache: boolean | null = null;

// Exportar função para reutilização (com cache)
export const fetchGlobalConfig = async (): Promise<GlobalConfig> => {
  // Se já tem cache, retorna
  if (globalConfigCache) {
    return globalConfigCache;
  }

  // Se já está carregando, retorna a promise existente
  if (globalConfigPromise) {
    return globalConfigPromise;
  }

  // Cria nova promise de carregamento
  globalConfigPromise = (async () => {
    try {
      const res = await api.get('/global_config');
      const data = (res?.data || {}) as GlobalConfig;
      globalConfigCache = data;
      return data;
    } catch (e) {
      console.error('[GlobalConfig] Failed to load from /api/v1/global_config', e);
      globalConfigCache = {};
      return {};
    } finally {
      globalConfigPromise = null;
    }
  })();

  return globalConfigPromise;
};

export const fetchSetupStatus = async (): Promise<boolean> => {
  if (setupRequiredCache !== null) {
    return setupRequiredCache;
  }

  try {
    const status = await setupService.getStatus();
    setupRequiredCache = status.status === 'inactive';
    return setupRequiredCache;
  } catch (e) {
    // Fail closed: when /setup/status errors transiently (network hiccup,
    // auth-service restart) we do NOT push the user to the bootstrap wizard.
    // Fresh installs are already handled correctly by the auth-service
    // (commit 2c7e6b8, EVO-971): /setup/status now reports `inactive`
    // whenever User.exists? is false, so a wiped install reaches /setup
    // through the success path, not through this fallback. The previous
    // fail-open default sent every logged-in user to the wizard on any
    // transient outage, where an operator could re-run bootstrap by
    // mistake.
    // Do not cache the failure — the next render retries.
    console.warn('[GlobalConfig] /setup/status failed, defaulting setupRequired=false', e);
    return false;
  }
};

// Listeners para notificar componentes React quando o cache é limpo
const setupCacheListeners: Set<() => void> = new Set();
const globalConfigListeners: Set<() => void> = new Set();

export const clearGlobalConfigCache = () => {
  globalConfigCache = null;
  globalConfigPromise = null;
  // Notificar listeners para re-fetch (ex.: após admin salvar uma config, os
  // booleans `hasXxxConfig` mudam e o fluxo de criação de canal precisa refletir)
  globalConfigListeners.forEach(listener => listener());
};

// Versão await-able: invalida o cache e aguarda o próximo fetch terminar.
// Usado pelo admin após salvar uma config para garantir que `hasXxxConfig`
// já está atualizado antes da UI prosseguir (sem janela de race entre o
// success toast e a próxima ação do usuário).
export const refreshGlobalConfig = async (): Promise<GlobalConfig> => {
  globalConfigCache = null;
  globalConfigPromise = null;
  const data = await fetchGlobalConfig();
  globalConfigListeners.forEach(listener => listener());
  return data;
};

export const clearSetupCache = () => {
  setupRequiredCache = null;
  // Notificar todos os listeners para re-fetch
  setupCacheListeners.forEach(listener => listener());
};

export const GlobalConfigProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [config, setConfig] = useState<GlobalConfig>(globalConfigCache || {});
  const [setupRequired, setSetupRequired] = useState<boolean>(setupRequiredCache ?? false);
  const [setupLoading, setSetupLoading] = useState<boolean>(setupRequiredCache === null);

  useEffect(() => {
    let mounted = true;

    const loadConfig = () => {
      Promise.all([fetchGlobalConfig(), fetchSetupStatus()]).then(([configData, isSetupRequired]) => {
        if (mounted) {
          setConfig(configData);
          setSetupRequired(isSetupRequired);
          setSetupLoading(false);
          // Initialize Clarity with backend-provided project ID
          if (configData.clarityProjectId) {
            initClarity(configData.clarityProjectId);
          }

          // Injetar White-Label (Nome e Cor Primária)
          if (configData.companyName) {
            document.title = configData.companyName;
          }
          if (configData.appPrimaryColor) {
            const hexToHsl = (hex: string) => {
              hex = hex.replace(/^#/, '');
              if (hex.length === 3) hex = hex.split('').map(x => x + x).join('');
              let r = parseInt(hex.substring(0, 2), 16) / 255;
              let g = parseInt(hex.substring(2, 4), 16) / 255;
              let b = parseInt(hex.substring(4, 6), 16) / 255;
              let max = Math.max(r, g, b), min = Math.min(r, g, b);
              let h = 0, s = 0, l = (max + min) / 2;
              if (max !== min) {
                let d = max - min;
                s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
                switch (max) {
                  case r: h = (g - b) / d + (g < b ? 6 : 0); break;
                  case g: h = (b - r) / d + 2; break;
                  case b: h = (r - g) / d + 4; break;
                }
                h /= 6;
              }
              return `${Math.round(h * 360)} ${Math.round(s * 100)}% ${Math.round(l * 100)}%`;
            };

            const hslValue = hexToHsl(configData.appPrimaryColor);
            
            // Set the raw components for standard Shadcn UI/Tailwind setups (e.g. `210 100% 50%`)
            document.documentElement.style.setProperty('--primary', hslValue);
            document.documentElement.style.setProperty('--color-primary', `hsl(${hslValue})`);
          }
        }
      });
    };

    loadConfig();

    // Re-fetch when setup cache is cleared (after bootstrap)
    const onCacheCleared = () => {
      if (mounted) {
        setSetupLoading(true);
        loadConfig();
      }
    };
    setupCacheListeners.add(onCacheCleared);

    // Re-fetch when global config cache is cleared (admin save invalidates it)
    const onGlobalConfigCleared = () => {
      if (mounted) loadConfig();
    };
    globalConfigListeners.add(onGlobalConfigCleared);

    return () => {
      mounted = false;
      setupCacheListeners.delete(onCacheCleared);
      globalConfigListeners.delete(onGlobalConfigCleared);
    };
  }, []);

  const value = useMemo(
    () => ({ ...config, setupRequired, setupLoading }),
    [config, setupRequired, setupLoading],
  );

  return <GlobalConfigContext.Provider value={value}>{children}</GlobalConfigContext.Provider>;
};

export const useGlobalConfig = (): GlobalConfigContextValue => useContext(GlobalConfigContext);
