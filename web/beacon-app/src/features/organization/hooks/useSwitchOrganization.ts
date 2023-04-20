import { useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';

import axiosInstance from '@/application/api/ApiService';
import type { SwitchMutation } from '@/features/organization/types/switchService';
import { setCookie } from '@/utils/cookies';

import { switchOrganizationRequest } from '../api/switchOrganizationApi';
export function useSwitchOrganization(): SwitchMutation {
  const mutation = useMutation(switchOrganizationRequest(axiosInstance), {
    onSuccess: (data) => {
      setCookie('bc_rtk', data.refresh_token);
      setCookie('bc_atk', data.access_token);
    },
    onError(error: any) {
      toast.error(error?.response?.data?.error);
    },
  });

  return {
    switch: mutation.mutate,
    reset: mutation.reset,
    hasSwitchFailed: mutation.isError,
    isSwitching: mutation.isLoading,
    auth: mutation.data as SwitchMutation['auth'],
    wasSwitchFetched: mutation.isSuccess,
    error: mutation.error,
  };
}
