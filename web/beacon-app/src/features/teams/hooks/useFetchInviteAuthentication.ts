import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { getInviteAuthenticationRequest } from '../api/getInviteAuthenticationRequest';

export function useFetchInviteAuthentication(token: string): any {
  const mutation = useMutation([RQK.INVITE_AUTHENTICATION, token], () =>
    getInviteAuthenticationRequest(axiosInstance)(token)
  );

  return {
    invite: mutation.mutate,
    reset: mutation.reset,
    auth: mutation.data,
    hasInviteAuthenticationFailed: mutation.isError,
    isFetchingInviteAuthentication: mutation.isLoading,
    wasInviteAuthenticated: mutation.isSuccess,
    error: mutation.error,
  };
}

export default useFetchInviteAuthentication;
