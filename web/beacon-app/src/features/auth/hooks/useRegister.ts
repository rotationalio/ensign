import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';

import { createAccountRequest } from '../api/RegisterApiService';
import type { RegistrationMutation } from '../types/CreateAccountService';

export function useRegister(): RegistrationMutation {
  const mutation = useMutation(createAccountRequest(axiosInstance), {
    retry: 0,
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
