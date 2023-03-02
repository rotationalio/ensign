import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { memberRequest } from '../api/memberListAPI';
import { MembersQuery } from '../types/memberServices';

export function useFetchMembers(): MembersQuery {
  const query = useQuery([RQK.MEMBER_LIST], memberRequest(axiosInstance), {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set stale time to 15 minutes
    staleTime: 1000 * 60 * 15,
  });

  return {
    getMembers: query.refetch,
    hasMembersFailed: query.isError,
    isFetchingMembers: query.isLoading,
    members: query.data,
    wasMembersFetched: query.isSuccess,
    error: query.error,
  };
}
