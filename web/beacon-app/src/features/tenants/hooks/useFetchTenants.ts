import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { tenantsRequest } from '../api/tenantListAPI';
import { TenantsQuery } from '../types/tenantServices';

export function useFetchTenants(): TenantsQuery {
  const query = useQuery([RQK.TENANTS], tenantsRequest(axiosInstance), {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set stale time to 15 minutes
    // TODO: Change stale time
    staleTime: 1000 * 60 * 15,
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
