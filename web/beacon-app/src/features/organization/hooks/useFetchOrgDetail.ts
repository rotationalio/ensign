import Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { orgRequest } from '../api/orgDetailApi';
import { OrgDetailQuery } from '../types/organizationService';

export function useFetchOrg(orgID: string): OrgDetailQuery {
  const query = useQuery([RQK.ORG_DETAIL, orgID], () => orgRequest(axiosInstance)(orgID), {
    enabled: !!orgID,
    onError: (error) => {
      Sentry.captureException(error);
    },
  });

  return {
    getOrgDetail: query.refetch,
    hasOrgFailed: query.isError,
    isFetchingOrg: query.isLoading,
    org: query.data as OrgDetailQuery['org'],
    wasOrgFetched: query.isSuccess,
    error: query.error,
  };
}
