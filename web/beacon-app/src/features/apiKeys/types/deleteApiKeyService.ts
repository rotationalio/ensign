import { UseMutateFunction } from '@tanstack/react-query';

import { APIKey } from './apiKeyService';

export interface DeleteAPIKeyMutation {
  deleteApiKey: UseMutateFunction<unknown, unknown, void, unknown>;
  reset(): void;
  key: APIKey;
  hasKeyDeletedFailed: boolean;
  wasKeyDeleted: boolean;
  isDeletingKey: boolean;
  error: any;
}

export type APIKeyDTO = {
  topicID: string;
};

export const isKeyDeleted = (
  mutation: DeleteAPIKeyMutation
): mutation is Required<DeleteAPIKeyMutation> =>
  mutation.wasKeyDeleted && mutation.key != undefined;
