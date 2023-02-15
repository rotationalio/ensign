import { APIKey } from './apiKeyService';

export interface APIKeyMutation {
  createNewKey(): void;
  reset(): void;
  key: APIKey;
  hasKeyFailed: boolean;
  wasKeyCreated: boolean;
  isCreatingKey: boolean;
  error: any;
}

export const isKeyCreated = (mutation: APIKeyMutation): mutation is Required<APIKeyMutation> =>
  mutation.wasKeyCreated && mutation.key != undefined;
