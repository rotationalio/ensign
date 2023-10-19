import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { profileRequest } from '../api/getProfileAPI';
import { ProfileQuery } from '../types/profileService';
export const profileQuery = () => ({
  queryKey: [RQK.PROFILE],
  queryFn: () => profileRequest(axiosInstance)(),
  cacheTime: 0,
});

export function useFetchProfile(): ProfileQuery {
  const query = useQuery({
    ...profileQuery(),
    onError: (error: any) => {
      // stop logging 401 & 403 errors to sentry
      if (error.response.status !== 401 && error.response.status !== 403) {
        Sentry.captureException(error);
      }
    },
  });

  return {
    getProfile: query.refetch,
    hasProfileFailed: query.isError,

    isFetchingProfile: query.isLoading,
    profile: query.data,
    wasProfileFetched: query.isSuccess,
    error: query.error,
  };
}
