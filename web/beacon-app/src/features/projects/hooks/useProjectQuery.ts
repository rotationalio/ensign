import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';

import { projectQueryAPI } from '../api/projectQueryApiService';
import { ProjectQueryDTO, ProjectQueryMutation } from '../types/projectQueryService';

export function useProjectQuery(): ProjectQueryMutation {
  const mutation = useMutation((payload: ProjectQueryDTO) =>
    projectQueryAPI(axiosInstance)(payload)
  );

  return {
    error: mutation.error,
    getProjectQuery: mutation.mutate,
    hasProjectQueryFailed: mutation.isError,
    isCreatingProjectQuery: mutation.isLoading,
    projectQuery: mutation.data,
    reset: mutation.reset,
    wasProjectQueryCreated: mutation.isSuccess,
  };
}
