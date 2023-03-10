import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { setCookie } from '@/utils/cookies';

import { loginRequest } from '../api/LoginApiService';
import type { LoginMutation } from '../types/LoginService';
export function useLogin(): LoginMutation {
  const mutation = useMutation(loginRequest(axiosInstance), {
    onSuccess: (data) => {
      setCookie('bc_rtk', data.refresh_token);
      setCookie('bc_atk', data.access_token);
    },
  });

  return {
    authenticate: mutation.mutate,
    reset: mutation.reset,
    hasAuthFailed: mutation.isError,
    isAuthenticating: mutation.isLoading,
    auth: mutation.data as LoginMutation['auth'],
    authenticated: mutation.isSuccess,
    error: mutation.error,
  };
}
