import { t } from '@lingui/macro';
import { useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { updateProfileAPI } from '../api/updateProfileAPI';
import { ProfileUpdateMutation } from '../types/profileService';

export function useUpdateProfile(): ProfileUpdateMutation {
  const mutation = useMutation(updateProfileAPI(axiosInstance), {
    retry: 0,
    onSuccess: (data: any) => {
      queryClient.invalidateQueries({ queryKey: [RQK.MEMBER_LIST] });
      queryClient.invalidateQueries({ queryKey: [RQK.MEMBER_DETAIL] });
      queryClient.invalidateQueries({ queryKey: [RQK.PROFILE] });
      queryClient.setQueriesData([RQK.PROFILE], (oldData: any) => {
        return { ...oldData, ...data };
      });
    },
    onError: (error: any) => {
      if (error?.response?.status !== 400 && error?.response?.status !== 409) {
        toast.error(
          error?.response?.data?.error ||
            t`Something went wrong. Please try again or contact support.`
        );
      }
    },
  });
  return {
    updateProfile: mutation.mutate,
    reset: mutation.reset,
    profile: mutation.data as ProfileUpdateMutation['profile'],
    hasProfileFailed: mutation.isError,
    wasProfileUpdated: mutation.isSuccess,
    isUpdatingProfile: mutation.isLoading,
    error: mutation.error,
  };
}
