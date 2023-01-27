import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { createAccountRequest } from '../api/RegisterApiService';
import type { RegistrationMutation } from '../types/CreateAccountService';
export function useRegister(): RegistrationMutation {
  const mutation = useMutation([RQK.CREATE_ACCOUNT], createAccountRequest(axiosInstance), {
    retry: 0,
  });

  return {
    createNewAccount: mutation.mutate,
    reset: mutation.reset,
    hasAccountFailed: mutation.isError,
    isCreatingAccount: mutation.isLoading,
    user: mutation.data as RegistrationMutation['user'],
    wasAccountCreated: mutation.isSuccess,
  };
}
