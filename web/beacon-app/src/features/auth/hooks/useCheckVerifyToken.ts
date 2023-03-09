import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { checkVerifyTokenRequest } from '../api/verifyTokenApiService';

export function useCheckVerifyToken(token: string): any {
  const mutation = useMutation([RQK.VERIFY_EMAIL, token], () =>
    checkVerifyTokenRequest(axiosInstance)(token)
  );
  return {
    verifyUserEmail: mutation.mutate,
    reset: mutation.reset,
    data: mutation.data,
    hasVerificationFailed: mutation.isError,
    wasVerificationChecked: mutation.isSuccess,
    isCheckingToken: mutation.isLoading,
    error: mutation.error,
  };
}
