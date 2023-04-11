import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { createMemberRequest } from '../api/createMemberApiService';
import { MemberMutation } from '../types/memberServices';

export function useCreateMember(): MemberMutation {
  const mutation = useMutation(createMemberRequest(axiosInstance), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.MEMBER_LIST] });
    },
  });
  return {
    createMember: mutation.mutate,
    reset: mutation.reset,
    member: mutation.data as MemberMutation['member'],
    hasMemberFailed: mutation.isError,
    wasMemberCreated: mutation.isSuccess,
    isCreatingMember: mutation.isLoading,
    error: mutation.error,
  };
}
