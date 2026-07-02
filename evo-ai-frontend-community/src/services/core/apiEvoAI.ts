import axios from 'axios';
import { useAuthStore } from '@/store/authStore';
import { applySetupInterceptor } from '@/services/core/setupInterceptor';

// Create a separate axios instance for evo-ai-core-service
const evoaiApi = axios.create({
  baseURL: `${import.meta.env.VITE_EVOAI_API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptador para incluir o access_token nos headers das requisições
evoaiApi.interceptors.request.use((config) => {
  // Incluir o access_token nos headers se estiver disponível
  const authHeader = useAuthStore.getState().getAuthHeader();
  if (authHeader) {
    // Usar forma compatível com o tipo AxiosHeaders
    config.headers.Authorization = authHeader.Authorization;
  }

  return config;
});

import apiAuth from '@/services/core/apiAuth';
import { InternalAxiosRequestConfig, AxiosError, AxiosRequestConfig } from 'axios';
import { requestMonitor } from '@/utils/requestMonitor';

let isRefreshing = false;
let isTerminatingSession = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: Error | null, token: string | null = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });

  failedQueue = [];
};

const terminateSession = () => {
  if (isTerminatingSession) return;
  isTerminatingSession = true;

  try {
    window.dispatchEvent(new CustomEvent('evolution:auth-lost'));
  } catch {
    // noop
  }

  useAuthStore.getState().clearUser();
};

const AUTH_INVALIDATION_ERROR_CODES = new Set<string>([
  'UNAUTHORIZED',
  'INVALID_TOKEN',
  'TOKEN_EXPIRED',
  'MISSING_TOKEN',
  'INVALID_CREDENTIALS',
  'SESSION_EXPIRED',
]);

evoaiApi.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(() => {
            const authHeader = useAuthStore.getState().getAuthHeader();
            if (authHeader && originalRequest.headers) {
              originalRequest.headers.Authorization = authHeader.Authorization;
            }
            return evoaiApi(originalRequest);
          })
          .catch(err => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const refreshResponse = await apiAuth.post('/auth/refresh');
        const refreshData = refreshResponse.data?.data || refreshResponse.data;
        const newAccessToken = refreshData?.access_token || refreshData?.token?.access_token;

        if (!newAccessToken) {
          throw new Error('New token not received');
        }

        useAuthStore.getState().setAccessToken(newAccessToken);
        processQueue(null, newAccessToken);

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        }

        isRefreshing = false;
        return evoaiApi(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError as Error, null);
        isRefreshing = false;
      }
    }

    if (error.response?.status === 401) {
      const errorCode = (
        error.response?.data as { error?: { code?: string } } | undefined
      )?.error?.code;
      
      const isAuthInvalidationCode =
        !errorCode ||
        AUTH_INVALIDATION_ERROR_CODES.has(errorCode);

      if (!isAuthInvalidationCode) {
        return Promise.reject(error);
      }

      terminateSession();
    }

    return Promise.reject(error);
  }
);

applySetupInterceptor(evoaiApi);

export default evoaiApi;
