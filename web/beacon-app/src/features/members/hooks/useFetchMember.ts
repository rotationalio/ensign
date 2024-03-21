import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { memberRequest } from '../api/memberApiService';
import { MemberQuery } from '../types/memberServices';
export const memberDetailQuery = (memberID: string) => ({
  queryKey: [RQK.MEMBER_DETAIL, memberID],
  queryFn: () => memberRequest(axiosInstance)(memberID),
  enabled: !!memberID,
});

export function useFetchMember(memberID: string): MemberQuery {
  const query = useQuery({
    ...memberDetailQuery(memberID),
    onError: (error: any) => {
      // stop logging 401 & 403 errors to sentry
      if (error?.response?.status !== 401 && error?.response?.status !== 403) {
        Sentry.captureException(error);
      }
    },
  });

  return {
    getMember: query.refetch,
    hasMemberFailed: query.isError,
    isFetchingMember: query.isLoading,
    member: query.data,
    wasMemberFetched: query.isSuccess,
    error: query.error,
  };
}
