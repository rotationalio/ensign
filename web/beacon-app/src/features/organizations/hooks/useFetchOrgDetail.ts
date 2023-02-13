import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { orgRequest } from '../api/orgDetailApi';
import { OrgDetailQuery } from '../types/organizationService';

export function useFetchOrg(orgID: string): OrgDetailQuery {
  const query = useQuery([RQK.ORG_DETAIL, orgID] as const, () => orgRequest(axiosInstance)(orgID), {
    enabled: !!orgID,
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set stale time to 15 minutes
    // TODO: Change stale time.
    staleTime: 1000 * 60 * 15,
  });

  return {
    hasOrgFailed: query.isError,
    isFetchingOrg: query.isLoading,
    org: query.data as OrgDetailQuery['org'],
    wasOrgFetched: query.isSuccess,
    error: query.error,
  };
}
