import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { memberRequest } from '../api/memberListAPI';
import { MemberQuery } from '../types/memberServices';

export function useFetchMembers(): MemberQuery {
  const query = useQuery([RQK.MEMBER_LIST], memberRequest(axiosInstance), {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set stale time to 15 minutes
    staleTime: 1000 * 60 * 15,
  });

  return {
    getMembers: query.refetch,
    hasMemberFailed: query.isError,
    isFetchingMember: query.isLoading,
    member: query.data,
    wasMemberFetched: query.isSuccess,
    error: query.error,
  };
}
