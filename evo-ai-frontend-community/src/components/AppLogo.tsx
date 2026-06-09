import type { CSSProperties } from 'react';
import { useDarkMode } from '../hooks/useDarkMode';
import { useGlobalConfig } from '@/contexts/GlobalConfigContext';
import logoDark from '../assets/EVO_CRM.svg';
import logoLight from '../assets/EVO_CRM_light.svg';

interface AppLogoProps {
  className?: string;
  alt?: string;
  style?: CSSProperties;
  forceTheme?: 'dark' | 'light';
}

export function AppLogo({ className, alt = 'EVO CRM', style, forceTheme }: AppLogoProps) {
  const { theme } = useDarkMode();
  const { appLogoUrl, companyName } = useGlobalConfig();
  
  const effectiveTheme = forceTheme ?? theme;
  const src = appLogoUrl || (effectiveTheme === 'dark' ? logoDark : logoLight);
  const effectiveAlt = companyName || alt;

  return <img src={src} alt={effectiveAlt} className={className} style={style} />;
}
