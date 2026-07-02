import axios, { InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/store/authStore';
import { applySetupInterceptor } from '@/services/core/setupInterceptor';
import { refreshAccessToken, terminateSession, AUTH_INVALIDATION_ERROR_CODES } from '@/services/core/tokenRefresh';

const agentProcessorApi = axios.create({
  baseURL: `${import.meta.env.VITE_AGENT_PROCESSOR_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

agentProcessorApi.interceptors.request.use(config => {
  const authHeader = useAuthStore.getState().getAuthHeader();
  if (authHeader) {
    config.headers.Authorization = authHeader.Authorization;
  }
  return config;
});

agentProcessorApi.interceptors.response.use(
  response => {
    if (response.data && (response.data.status === 'success' || response.data.success === true) && response.data.data !== undefined) {
      response.data = response.data.data;
    }
    return response;
  },
  async error => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const newToken = await refreshAccessToken();
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
        }
        return agentProcessorApi(originalRequest);
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

      if (isAuthInvalidationCode) {
        terminateSession();
      }
    }

    const detail =
      error?.response?.data?.error?.message ||
      error?.response?.data?.detail ||
      error?.response?.data?.message ||
      error?.message ||
      'Unknown error';
    console.error('Agent Processor API Error:', detail, { status: error?.response?.status });
    return Promise.reject(error);
  },
);

applySetupInterceptor(agentProcessorApi);

export { agentProcessorApi };
export default agentProcessorApi;
