import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { createTenantRequest } from '../api/createTenantApiService';
import { TenantMutation } from '../types/createTenantService';

export function useCreateTenant(): TenantMutation {
  const mutation = useMutation([RQK.TENANTS], createTenantRequest(axiosInstance), {
    retry: 0,
  });

  return {
    createTenant: mutation.mutate,
    tenant: mutation.data,
    hasTenantFailed: mutation.isError,
    isFetchingTenant: mutation.isLoading,
    wasTenantFetched: mutation.isSuccess,
    error: mutation.error,
  };
}
