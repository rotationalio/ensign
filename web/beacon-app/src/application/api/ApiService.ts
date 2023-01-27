import axios, { AxiosError, AxiosResponse } from 'axios';

import { appConfig } from '@/application/config';

const axiosInstance = axios.create({
  baseURL: `${appConfig.apiUrl}/${appConfig.apiVersion}`,
  headers: {
    'Content-Type': 'application/json',
  },
});

axiosInstance.defaults.withCredentials = true;
// intercept request and check if token has expired or not
axiosInstance.interceptors.request.use(
  async (config: any) => {
    // handle refresh token here at every request
    return config;
  },
  (error) => {
    Promise.reject(error);
  }
);
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    return Promise.reject(error);

    // }
  }
);

export type Request = (url: string, options?: any) => Promise<Response>;

export const getValidApiResponse = <T>(
  response: Pick<AxiosResponse, 'status' | 'data' | 'statusText'>
): T => {
  if (response?.status === 200 || response?.status === 204) {
    return response?.data as T;
  }
  throw new Error(response.statusText);
};
export const getValidApiError = (error: AxiosError): Error => {
  // later we can handle error here by catching axios error code
  return new Error(error.message);
};

export const setAuthorization = () => {
  axiosInstance.defaults.headers.common.Authorization = `Bearer token`;
  // axiosInstance.defaults.headers.common['X-CSRF-Token'] = csrfToken;
};

export default axiosInstance;
