import axios, { AxiosError, AxiosResponse } from 'axios';

// import QuarterDeckAuth from '@/lib/quaterdeck-auth';
import { appConfig } from '@/application/config';
import { getCookie, setCookie } from '@/utils/cookies';
import { decodeToken } from '@/utils/decodeToken';

const axiosInstance = axios.create({
  baseURL: `${appConfig.tenantApiUrl}`,
  headers: {
    'Content-Type': 'application/json',
  },
});

axiosInstance.defaults.withCredentials = true;
// intercept request and check if token has expired or not
axiosInstance.interceptors.request.use(
  async (config: any) => {
    const token = getCookie('bc_atk');
    const decodedToken = token && decodeToken(token);
    if (decodedToken) {
      const { exp } = decodedToken;
      const now = new Date().getTime() / 1000;
      if (exp < now) {
        // refresh token
      }
    }
    config.headers.Authorization = `Bearer ${token}`;
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
  if (response?.status === 200 || response?.status === 201) {
    return response?.data as T;
  }
  if (response?.status === 204) {
    return {} as T;
  }
  throw new Error(response?.data);
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

export const setAuthorization = () => {
  axiosInstance.defaults.headers.common.Authorization = `Bearer token`;
  // axiosInstance.defaults.headers.common['X-CSRF-Token'] = csrfToken;
};

export const refreshToken = async () => {
  const refreshToken = getCookie('bc_rtk');
  const accessToken = getCookie('bc_atk');
  if (refreshToken) {
    const d = decodeToken(accessToken) as any;
    const exp = d?.exp;
    const now = new Date().getTime() / 1000;
    if (exp < now) {
      const response = await axiosInstance.post('/refresh', {
        data: JSON.stringify({
          refresh_token: refreshToken,
        }),
      });
      if (response.status === 200) {
        const { access_token, refresh_token } = response.data;
        setCookie('bc_atk', access_token);
        setCookie('bc_rtk', refresh_token);
        axiosInstance.defaults.headers.common.Authorization = `Bearer ${access_token}`;
      }
    }
  }
};

export default axiosInstance;
