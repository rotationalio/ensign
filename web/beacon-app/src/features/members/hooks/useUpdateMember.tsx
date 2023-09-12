import { t } from '@lingui/macro';
import { useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { updateMemberAPI } from '../api/updateMemberAPI';
import { MemberUpdateMutation } from '../types/memberServices';

export function useUpdateMember(): MemberUpdateMutation {
  const mutation = useMutation(updateMemberAPI(axiosInstance), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.MEMBER_LIST] });
      queryClient.invalidateQueries({ queryKey: [RQK.MEMBER_DETAIL] });
    },
    onError: (error: any) => {
      if (error?.response?.status !== 400) {
        toast.error(
          error?.response?.data?.error ||
            t`Something went wrong. Please try again or contact support.`
        );
      }
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
