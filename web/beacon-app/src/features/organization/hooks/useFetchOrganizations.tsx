import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { organizationRequest } from '../api/organizationListApi';
import { OrgListQuery } from '../types/organizationService';

export function useFetchOrganizations(): OrgListQuery {
  const query = useQuery([RQK.ORGANIZATION_LIST], organizationRequest(axiosInstance), {
    onError(error: any) {
      toast.error(error?.response?.data?.error || 'Something went wrong');
    },
  });

  return {
    getOrgList: query.refetch,
    hasOrgListFailed: query.isError,
    isFetchingOrgList: query.isLoading,
    organizations: query.data,
    wasOrgListFetched: query.isSuccess,
    error: query.error,
  };
}
