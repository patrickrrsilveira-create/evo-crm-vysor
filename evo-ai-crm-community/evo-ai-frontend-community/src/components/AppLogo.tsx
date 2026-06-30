import type { CSSProperties } from 'react';
import { useDarkMode } from '../hooks/useDarkMode';
import { useGlobalConfig } from '@/contexts/GlobalConfigContext';
import logoTextDark from '../assets/vysor_logo_text.png';
import logoTextLight from '../assets/vysor_logo_text.png';
import logoFullDark from '../assets/vysor_logo_full.png';
import logoFullLight from '../assets/vysor_logo_full.png';

interface AppLogoProps {
  className?: string;
  alt?: string;
  style?: CSSProperties;
  forceTheme?: 'dark' | 'light';
  variant?: 'default' | 'login';
}

export function AppLogo({ className, alt = 'VYSOR CRM', style, forceTheme, variant = 'default' }: AppLogoProps) {
  const { theme } = useDarkMode();
  const config = useGlobalConfig();
  
  const effectiveTheme = forceTheme ?? theme;
  
  const isLogin = variant === 'login';
  const logoUrl = isLogin && config.appLoginLogoUrl ? config.appLoginLogoUrl : config.appLogoUrl;
  const logoWidth = isLogin && config.appLoginLogoWidth ? config.appLoginLogoWidth : config.appLogoWidth;
  const logoHeight = isLogin && config.appLoginLogoHeight ? config.appLoginLogoHeight : config.appLogoHeight;

  const defaultLogo = isLogin ? (effectiveTheme === 'dark' ? logoFullDark : logoFullLight) : (effectiveTheme === 'dark' ? logoTextDark : logoTextLight);
  const src = logoUrl || defaultLogo;
  const effectiveAlt = config.companyName || alt;

  const finalStyle = {
    ...style,
    ...(logoWidth ? { width: `${logoWidth}px` } : {}),
    ...(logoHeight ? { height: `${logoHeight}px` } : {}),
  };

  return <img src={src} alt={effectiveAlt} className={className} style={finalStyle} />;
}
