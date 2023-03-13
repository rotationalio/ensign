import { UseMutateFunction } from '@tanstack/react-query';

import { APIKey } from './apiKeyService';

export interface APIKeyMutation {
  createProjectNewKey: UseMutateFunction<APIKey, unknown, APIKeyDTO, unknown>;
  reset(): void;
  key: APIKey;
  hasKeyFailed: boolean;
  wasKeyCreated: boolean;
  isCreatingKey: boolean;
  error: any;
}

export interface NewAPIKey {
  name: string;
  permissions: string[];
}

export type APIKeyDTO = {
  projectID: string;
} & NewAPIKey;

export const isKeyCreated = (mutation: APIKeyMutation): mutation is Required<APIKeyMutation> =>
  mutation.wasKeyCreated && mutation.key != undefined;
