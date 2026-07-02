import axios, { InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/store/authStore';
import { applySetupInterceptor } from '@/services/core/setupInterceptor';
import { refreshAccessToken, terminateSession, AUTH_INVALIDATION_ERROR_CODES } from '@/services/core/tokenRefresh';

const evoFlowApi = axios.create({
  baseURL: `${import.meta.env.VITE_EVOFLOW_API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

evoFlowApi.interceptors.request.use((config) => {
  const authHeader = useAuthStore.getState().getAuthHeader();
  if (authHeader) {
    config.headers.Authorization = authHeader.Authorization;
  }
  return config;
});

evoFlowApi.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const newToken = await refreshAccessToken();
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
        }
        return evoFlowApi(originalRequest);
      } catch {
        // refresh failed
      }
    }

    if (error.response?.status === 401) {
      const errorCode = (
        error.response?.data as { error?: { code?: string } } | undefined
      )?.error?.code;

      const isAuthInvalidationCode =
        !errorCode ||
        AUTH_INVALIDATION_ERROR_CODES.has(errorCode);

      if (isAuthInvalidationCode) {
        terminateSession();
      }
    }

    return Promise.reject(error);
  }
);

applySetupInterceptor(evoFlowApi);

export default evoFlowApi;
