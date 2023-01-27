import { useMutation } from '@tanstack/react-query';
import { loginRequest } from '../api/LoginApiService';
import type { LoginMutation } from '../types/LoginService';
import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';
export function useLogin(): LoginMutation {
    const mutation = useMutation([RQK.LOGIN], loginRequest(axiosInstance), {
        retry: 0,
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
