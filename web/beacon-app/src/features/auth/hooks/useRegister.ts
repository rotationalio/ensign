import { t } from '@lingui/macro';
import { useMutation } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';

import axiosInstance from '@/application/api/ApiService';

import { createAccountRequest } from '../api/RegisterApiService';
import type { RegistrationMutation } from '../types/CreateAccountService';

export function useRegister(): RegistrationMutation {
  const mutation = useMutation(createAccountRequest(axiosInstance), {
    retry: 0,
    onError(error: any) {
      if (error.response.status === 409) {
        toast.error(t`User already exists.`);
      } else {
        toast.error(error?.response?.data?.error);
      }
    },
  });

  return {
    createNewAccount: mutation.mutate,
    reset: mutation.reset,
    hasAccountFailed: mutation.isError,
    error: mutation.error,
    isCreatingAccount: mutation.isLoading,
    user: mutation.data as RegistrationMutation['user'],
    wasAccountCreated: mutation.isSuccess,
  };
}
