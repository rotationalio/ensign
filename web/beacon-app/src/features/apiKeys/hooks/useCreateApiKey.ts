import { createAPIKey } from "../api/CreateApiKey";
import { RQK } from "@/constants";
import { useMutation } from "@tanstack/react-query";
import axiosInstance from "@/application/api/ApiService";
import { APIKeyMutation } from "../types/CreateAPIKeyService";

export function useCreateAPIKey(): APIKeyMutation {
    const mutation = useMutation([RQK.CREATE_KEY], createAPIKey(axiosInstance), {
     retry: 0
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