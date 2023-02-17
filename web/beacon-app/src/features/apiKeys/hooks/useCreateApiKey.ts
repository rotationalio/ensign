import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
// import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { createAPIKey } from '../api/createApiKey';
import { APIKeyMutation } from '../types/createApiKeyService';

export function useCreateAPIKey(): APIKeyMutation {
  const mutation = useMutation([RQK.CREATE_KEY], createAPIKey(axiosInstance));
  return {
    createNewKey: mutation.mutate,
    reset: mutation.reset,
    key: mutation.data as APIKeyMutation['key'],
    hasKeyFailed: mutation.isError,
    wasKeyCreated: mutation.isSuccess,
    isCreatingKey: mutation.isLoading,
    error: mutation.error,
  };
}
