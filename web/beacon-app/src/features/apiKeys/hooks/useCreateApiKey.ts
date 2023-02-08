import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { createAPIKey } from '../api/CreateApiKey';
import { APIKeyMutation } from '../types/CreateApiKeyService';

export function useCreateAPIKey(): APIKeyMutation {
  const mutation = useMutation([RQK.CREATE_KEY], createAPIKey(axiosInstance), {
    retry: 0,
  });

  return {
    createNewKey: mutation.mutate,
    reset: mutation.reset,
    key: mutation.data as APIKeyMutation['key'],
    hasKeyFailed: mutation.isError,
    wasKeyCreated: mutation.isSuccess,
    isCreatingKey: mutation.isLoading,
  };
}
