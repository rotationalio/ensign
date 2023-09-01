import Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { tenantsRequest } from '../api/tenantListAPI';
import { TenantsQuery } from '../types/tenantServices';

export function useFetchTenants(): TenantsQuery {
  const query = useQuery([RQK.TENANTS], tenantsRequest(axiosInstance), {
    onError: (error) => {
      Sentry.captureException(error);
    },
  });

  return {
    getTenants: query.refetch,
    hasTenantsFailed: query.isError,
    isFetchingTenants: query.isLoading,
    tenants: query.data,
    wasTenantsFetched: query.isSuccess,
    error: query.error,
  };
}
