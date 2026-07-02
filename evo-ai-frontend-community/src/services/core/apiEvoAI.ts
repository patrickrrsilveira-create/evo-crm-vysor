import axios, { InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/store/authStore';
import { applySetupInterceptor } from '@/services/core/setupInterceptor';
import { refreshAccessToken, terminateSession, AUTH_INVALIDATION_ERROR_CODES } from '@/services/core/tokenRefresh';

const evoaiApi = axios.create({
  baseURL: `${import.meta.env.VITE_EVOAI_API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

evoaiApi.interceptors.request.use((config) => {
  const authHeader = useAuthStore.getState().getAuthHeader();
  if (authHeader) {
    config.headers.Authorization = authHeader.Authorization;
  }

  return config;
});

evoaiApi.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const newToken = await refreshAccessToken();
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
        }
        return evoaiApi(originalRequest);
      } catch {
        // refresh failed, fall through to session termination check
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
