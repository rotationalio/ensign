import { useMutation } from '@tanstack/react-query';
import { RQK } from '@/constants/queryKeys';
import { createAccountRequest } from './api/AccountApiService';
import axiosInstance from '@/application/api/ApiService';
import type { RegistrationMutation } from './CreateAccountService';


export function useCreateAccount(): RegistrationMutation {
    const mutation = useMutation([RQK.CREATE_ACCOUNT], createAccountRequest(axiosInstance));

    return {
        createNewAccount: mutation.mutate,
        reset: mutation.reset,
        hasAccountFailed: mutation.isError,
        isCreatingAccount: mutation.isLoading,
        user: mutation.data as RegistrationMutation['user'],
        wasAccountCreated: mutation.isSuccess,
    };
}

