import { t } from '@lingui/macro';
import * as Sentry from '@sentry/react';
import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';

import { forgotPasswordRequest } from '../api/ForgotPasswordApiService';
import { ForgotPasswordMutation } from '../types/ForgotPasswordService';

export function useForgotPassword(): ForgotPasswordMutation {
  const mutation = useMutation(forgotPasswordRequest(axiosInstance), {
    cacheTime: 0,
    onError(error: any) {
      Sentry.captureException(error, {
        extra: {
          message: t`Forgot password request failed.`,
        },
      });
    },
  });
  return {
    forgotPassword: mutation.mutate,
    reset: mutation.reset,
    hasForgotPasswordFailed: mutation.isError,
    wasForgotPasswordSuccessful: mutation.isSuccess,
    isLoading: mutation.isLoading,
    error: mutation.error,
  };
}
