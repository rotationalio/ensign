import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { updateMemberAPI } from '../api/updateMemberAPI';
import { MemberUpdateMutation } from '../types/onboardingServices';

export function useUpdateMember(): MemberUpdateMutation {
  const mutation = useMutation(updateMemberAPI(axiosInstance), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.MEMBER_LIST] });
    },
  });
  return {
    updateMember: mutation.mutate,
    reset: mutation.reset,
    member: mutation.data as MemberUpdateMutation['member'],
    hasMemberFailed: mutation.isError,
    wasMemberUpdated: mutation.isSuccess,
    isUpdatingMember: mutation.isLoading,
    error: mutation.error,
  };
}
