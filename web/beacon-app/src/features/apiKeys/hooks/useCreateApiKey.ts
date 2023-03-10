import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { createProjectAPIKey } from '../api/createApiKey';
import { APIKeyMutation } from '../types/createApiKeyService';

export function useCreateProjectAPIKey(projectID: string): APIKeyMutation {
  const mutation = useMutation(() => createProjectAPIKey(axiosInstance)(projectID), {
    onSuccess: () => {
      // Invalidate the projects list query so that the new key is reflected in the UI
      queryClient.invalidateQueries([RQK.PROJECTS]);
    },
  });
  return {
    createProjectNewKey: mutation.mutate,
    reset: mutation.reset,
    key: mutation.data as APIKeyMutation['key'],
    hasKeyFailed: mutation.isError,
    wasKeyCreated: mutation.isSuccess,
    isCreatingKey: mutation.isLoading,
    error: mutation.error,
  };
}
