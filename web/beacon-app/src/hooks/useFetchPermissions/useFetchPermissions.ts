import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import permissionsRequest from './permissionsApiService';

function useFetchPermissions() {
  const query = useQuery([RQK.PERMISSIONS], permissionsRequest(axiosInstance), {
    onError: (error) => {
      Sentry.captureException(error);
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
