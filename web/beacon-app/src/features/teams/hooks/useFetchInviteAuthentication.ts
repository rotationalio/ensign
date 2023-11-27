import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { invitationAuthenticationRequest } from '../api/invitationAuthenticationRequest';

export function useFetchInviteAuthentication(token: string): any {
  const mutation = useMutation([RQK.INVITATION_AUTHENTICATION, token], () =>
    invitationAuthenticationRequest(axiosInstance)(token)
  );

  return {
    invitationRequest: mutation.mutate,
    reset: mutation.reset,
    authData: mutation.data,
    hasInvitationAuthenticationFailed: mutation.isError,
    isFetchingInvitationAuthentication: mutation.isLoading,
    wasInvitationAuthenticated: mutation.isSuccess,
    error: mutation.error,
  };
}

export default useFetchInviteAuthentication;
