import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';

import { projectQueryAPI } from '../api/projectQueryApiService';
import { ProjectQueryDTO, ProjectQueryMutation } from '../types/projectQueryService';

export function useProjectQuery(): ProjectQueryMutation {
  const query = useMutation((payload: ProjectQueryDTO) => projectQueryAPI(axiosInstance)(payload), {
    onSuccess: () => {
      query.reset();
    },
  });

  return {
    error: query.error,
    getProjectQuery: query.mutate,
    hasProjectQueryFailed: query.isError,
    isCreatingProjectQuery: query.isLoading,
    projectQuery: query.data,
    reset: query.reset,
    wasProjectQueryCreated: query.isSuccess,
  };
}
