import { UseMutateFunction } from '@tanstack/react-query';

import { NewUserAccount, NewUserResponseData } from './RegisterService';

export interface RegistrationMutation {
  createNewAccount: UseMutateFunction<NewUserResponseData, unknown, NewUserAccount, unknown>;
  reset(): void;
  user: NewUserResponseData;
  hasAccountFailed: boolean;
  wasAccountCreated: boolean;
  isCreatingAccount: boolean;
  error: unknown;
}

export const isAccountCreated = (
  mutation: RegistrationMutation
): mutation is Required<RegistrationMutation> =>
  mutation.wasAccountCreated && mutation.user !== undefined;
