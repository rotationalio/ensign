import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { tenantRequest } from '../api/tenantListAPI';
import { TenantQuery } from '../types/tenantServices';

export function useFetchTenants(): TenantQuery {
  const query = useQuery([RQK.TENANTS], tenantRequest(axiosInstance), {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set stale time to 15 minutes
    // TODO: Change stale time
    staleTime: 1000 * 60 * 15,
  });

  return {
    getTenant: query.refetch,
    hasTenantFailed: query.isError,
    isFetchingTenant: query.isLoading,
    tenants: query.data,
    wasTenantFetched: query.isSuccess,
    error: query.error,
  };
}
a;
