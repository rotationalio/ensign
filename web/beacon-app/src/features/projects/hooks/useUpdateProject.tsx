import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { updateProjectAPI } from '../api/updateProjectApiService';
import { ProjectUpdateMutation } from '../types/updateProjectService';

export function useUpdateProject(): ProjectUpdateMutation {
  const mutation = useMutation(updateProjectAPI(axiosInstance), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECTS] });
      queryClient.invalidateQueries({ queryKey: [RQK.QUICK_VIEW] });
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECT_QUICK_VIEW] });
    },
  });
  return {
    updateProject: mutation.mutate,
    reset: mutation.reset,
    project: mutation.data as ProjectUpdateMutation['project'],
    hasProjectFailed: mutation.isError,
    wasProjectCreated: mutation.isSuccess,
    isCreatingProject: mutation.isLoading,
    error: mutation.error,
  };
}
