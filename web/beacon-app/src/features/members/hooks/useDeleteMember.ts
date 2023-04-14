import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { deleteMemberRequest } from '../api/deleteMemberApi';
import { MemberDeleteMutation } from '../types/memberServices';

export function useDeleteMember(memberId: string): MemberDeleteMutation {
  const mutation = useMutation(() => deleteMemberRequest(axiosInstance)(memberId), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.MEMBER_LIST] });
    },
  });

  return {
    deleteMember: mutation.mutate,
    reset: mutation.reset,
    member: mutation.data,
    hasMemberFailed: mutation.isError,
    wasMemberDeleted: mutation.isSuccess,
    isDeletingMember: mutation.isLoading,
    error: mutation.error,
  };
}
