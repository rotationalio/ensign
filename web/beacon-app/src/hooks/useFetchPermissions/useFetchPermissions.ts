import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import permissionsRequest from './permissionsApiService';

function useFetchPermissions() {
  const query = useQuery([RQK.PERMISSIONS], permissionsRequest(axiosInstance), {
    onError: (error: any) => {
      // stop logging 401 & 403 errors to sentry
      if (error.response.status !== 401 && error.response.status !== 403) {
        Sentry.captureException(error);
      }
    },
  });
  return {
    getPermissions: query.refetch,
    hasPermissionsFailed: query.isError,
    isFetchingPermissions: query.isLoading,
    permissions: query.data,
    wasPermissionsFetched: query.isSuccess,
    error: query.error,
  };
}

export default useFetchPermissions;
