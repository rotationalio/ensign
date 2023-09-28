import axios, { AxiosError, AxiosResponse } from 'axios';
import toast from 'react-hot-toast';

import { appConfig } from '@/application/config';
import { clearCookies, getCookie } from '@/utils/cookies';
import ErrorMessage from '@/utils/error-message';

const axiosInstance = axios.create({
  baseURL: `${appConfig.tenantApiUrl}`,
  headers: {
    'Content-Type': 'application/json',
  },
});

axiosInstance.defaults.withCredentials = true;
axiosInstance.interceptors.request.use(
  async (config: any) => {
    // As the server stores the token in an HttpOnly cookie,
    // the access token will be included automatically in the Authorization header of each request.

    const csrfToken = getCookie('csrf_token');
    if (csrfToken) {
      config.headers['X-CSRF-Token'] = csrfToken;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    // if error code is ERR_NETWORK toast network error
    // https://github.com/axios/axios#handling-errors
    if (error?.code === 'ERR_NETWORK') {
      toast.error(`${ErrorMessage.NETWORK_ERROR}`);
    }
    // if status is 401 then clear cookies and logout user
    if (error?.response?.status === 401) {
      // logout();
      clearCookies();
      window.location.href = '/';
    }
    return Promise.reject(error);
  }
);

export type Request = (url: string, options?: any) => Promise<Response>;

export const getValidApiResponse = <T>(
  response: Pick<AxiosResponse, 'status' | 'data' | 'statusText'>
): T => {
  if (response?.status === 200 || response?.status === 201) {
    return response?.data as T;
  }
  if (response?.status === 204) {
    return {} as T;
  }
  throw new Error(response?.statusText);
};
export const getValidApiError = (error: AxiosError): Error => {
  // later we can handle error here by catching axios error code
  const errorMessage = error?.response?.data as any;

  switch (error?.response?.status) {
    case 400:
      // handle 400 error
      return new Error(errorMessage && errorMessage.message ? errorMessage.message : 'Bad Request');
      break;
    case 401:
      // handle 401 error
      return new Error(
        errorMessage && errorMessage.message ? errorMessage.message : 'Unauthorized'
      );
      break;
    case 403:
      // handle 403 error
      return new Error('Forbidden');
      break;
    case 404:
      // handle 404 error
      return new Error('Not Found');
      break;

    default:
      return new Error('Something went wrong');
      break;
  }
};

export default axiosInstance;
