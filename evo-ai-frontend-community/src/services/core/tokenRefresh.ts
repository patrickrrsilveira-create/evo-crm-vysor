import apiAuth from './apiAuth';
import { useAuthStore } from '@/store/authStore';

let activeRefreshPromise: Promise<string> | null = null;
let isTerminating = false;

export function refreshAccessToken(): Promise<string> {
  if (activeRefreshPromise) return activeRefreshPromise;

  activeRefreshPromise = (async () => {
    const response = await apiAuth.post('/auth/refresh');
    const data = response.data?.data || response.data;
    const token = data?.access_token || data?.token?.access_token;
    if (!token) throw new Error('New token not received');
    useAuthStore.getState().setAccessToken(token);
    return token;
  })().finally(() => {
    activeRefreshPromise = null;
  });

  return activeRefreshPromise;
}

export function terminateSession(): void {
  if (isTerminating) return;
  isTerminating = true;

  try {
    window.dispatchEvent(new CustomEvent('evolution:auth-lost'));
  } catch {
    // noop
  }

  useAuthStore.getState().clearUser();
}

export const AUTH_INVALIDATION_ERROR_CODES = new Set<string>([
  'UNAUTHORIZED',
  'INVALID_TOKEN',
  'TOKEN_EXPIRED',
  'MISSING_TOKEN',
  'INVALID_CREDENTIALS',
  'SESSION_EXPIRED',
]);
