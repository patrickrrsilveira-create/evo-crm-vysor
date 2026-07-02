import axios, { AxiosRequestConfig, AxiosError, InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/store/authStore';
import { requestMonitor } from '@/utils/requestMonitor';
import { applySetupInterceptor } from '@/services/core/setupInterceptor';
import { refreshAccessToken, terminateSession, AUTH_INVALIDATION_ERROR_CODES } from '@/services/core/tokenRefresh';

const api = axios.create({
  baseURL: `${import.meta.env.VITE_API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(config => {
  const requestId = requestMonitor.logRequest(
    config.method?.toUpperCase() || 'GET',
    config.url || '',
  );

  (config as AxiosRequestConfig & { requestId?: string; requestStartTime?: number }).requestId = requestId;
  (config as AxiosRequestConfig & { requestId?: string; requestStartTime?: number }).requestStartTime = Date.now();

  const authHeader = useAuthStore.getState().getAuthHeader();
  if (authHeader) {
    config.headers.Authorization = authHeader.Authorization;
  }

  if (config.data instanceof FormData) {
    delete config.headers['Content-Type'];
    delete config.headers['content-type'];
  }

  return config;
});

api.interceptors.response.use(
  response => {
    const config = response.config as AxiosRequestConfig & { requestId?: string; requestStartTime?: number };
    if (config.requestId && config.requestStartTime) {
      const duration = Date.now() - config.requestStartTime;
      requestMonitor.logResponse(config.requestId, response.status, duration);
    }

    return response;
  },
  async error => {
    const config = (error as AxiosError).config as (AxiosRequestConfig & { requestId?: string }) | undefined;
    if (config?.requestId) {
      const errorData = (error as AxiosError).response?.data as
        | { error?: { message?: string }; message?: string }
        | undefined;
      const errorMessage =
        errorData?.error?.message ||
        errorData?.message ||
        (error as AxiosError).message ||
        'Unknown error';
      requestMonitor.logError(config.requestId, errorMessage);
    }

    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const newToken = await refreshAccessToken();
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
        }
        return api(originalRequest);
      } catch {
        // refresh failed, fall through to session termination check
      }
    }

    if (error.response?.status === 401) {
      const isUnreadCountEndpoint = error.config?.url?.includes('/unread_count');

      if (isUnreadCountEndpoint) {
        return Promise.reject(error);
      }

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

    if (error.response?.status === 403) {
      return Promise.reject(error);
    }

    return Promise.reject(error);
  },
);

applySetupInterceptor(api);

export default api;
