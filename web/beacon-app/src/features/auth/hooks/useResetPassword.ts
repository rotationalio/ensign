import { t } from '@lingui/macro';
import * as Sentry from '@sentry/react';
import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';

import { resetPasswordRequest } from '../api/ResetPasswordApi';
import { ResetPasswordMutation } from '../types/ResetPasswordService';

export function useResetPassword(): ResetPasswordMutation {
  const mutation = useMutation(resetPasswordRequest(axiosInstance), {
    onError(error: any) {
      Sentry.captureException(error, {
        extra: {
          message: t`Reset password request failed.`,
        },
      });
    },
  });
  return {
    resetPassword: mutation.mutate,
    reset: mutation.reset,
    hasResetPasswordFailed: mutation.isError,
    wasResetPasswordSuccessful: mutation.isSuccess,
    isLoading: mutation.isLoading,
    error: mutation.error,
  };
}
