import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';

import { createTenantRequest } from '../api/createTenantApiService';
import { TenantMutation } from '../types/createTenantService';

export function useCreateTenant(): TenantMutation {
  const mutation = useMutation(createTenantRequest(axiosInstance), {
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
