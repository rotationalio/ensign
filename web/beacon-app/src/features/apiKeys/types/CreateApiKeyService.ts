import { APIKey, NewAPIKey } from './ApiKeyServices';

export interface APIKeyMutation {
  createNewKey(key: NewAPIKey): void;
  reset(): void;
  key: APIKey;
  hasKeyFailed: boolean;
  wasKeyCreated: boolean;
  isCreatingKey: boolean;
}

export const isKeyCreated = (mutation: APIKeyMutation): mutation is Required<APIKeyMutation> =>
  mutation.wasKeyCreated && mutation.key != undefined;
