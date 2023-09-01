import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { memberRequest } from '../api/memberListAPI';
import { MembersQuery } from '../types/memberServices';

export function useFetchMembers(): MembersQuery {
  const query = useQuery([RQK.MEMBER_LIST], memberRequest(axiosInstance), {
    onError(error: any) {
      Sentry.captureException(error);
      toast.error(error?.response?.data?.error || 'Something went wrong');
    },
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
