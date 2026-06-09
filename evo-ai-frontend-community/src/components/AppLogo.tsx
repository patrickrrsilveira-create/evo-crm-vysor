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
  variant?: 'default' | 'login';
}

export function AppLogo({ className, alt = 'EVO CRM', style, forceTheme, variant = 'default' }: AppLogoProps) {
  const { theme } = useDarkMode();
  const config = useGlobalConfig();
  
  const effectiveTheme = forceTheme ?? theme;
  
  const isLogin = variant === 'login';
  const logoUrl = isLogin && config.appLoginLogoUrl ? config.appLoginLogoUrl : config.appLogoUrl;
  const logoWidth = isLogin && config.appLoginLogoWidth ? config.appLoginLogoWidth : config.appLogoWidth;
  const logoHeight = isLogin && config.appLoginLogoHeight ? config.appLoginLogoHeight : config.appLogoHeight;

  const src = logoUrl || (effectiveTheme === 'dark' ? logoDark : logoLight);
  const effectiveAlt = config.companyName || alt;

  const finalStyle = {
    ...style,
    ...(logoWidth ? { width: `${logoWidth}px` } : {}),
    ...(logoHeight ? { height: `${logoHeight}px` } : {}),
  };

  return <img src={src} alt={effectiveAlt} className={className} style={finalStyle} />;
}
