import { t } from '@lingui/macro';
import * as Sentry from '@sentry/react';
import { useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';

import axiosInstance from '@/application/api/ApiService';

import { loginRequest } from '../api/LoginApiService';
import type { LoginMutation } from '../types/LoginService';
export function useLogin(): LoginMutation {
  const mutation = useMutation(loginRequest(axiosInstance), {
    onError(error: any) {
      Sentry.captureException(error, {
        extra: {
          message: 'Login failed',
        },
      });
      toast.error(error?.response?.data?.error || t`Something went wrong, please try again later.`);
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
