import { UseMutateFunction } from '@tanstack/react-query';

import { UserAuthResponse } from '@/features/auth';

export interface InviteAuthenticationMutation {
  invite: UseMutateFunction<UserAuthResponse, unknown, InviteAuthenticationDTO, unknown>;
  reset(): void;
  auth: UserAuthResponse;
  hasInviteAuthenticationFailed: boolean;
  wasAuthenticated: boolean;
  isFetchingInviteAuthentication: boolean;
  error: any;
}

export interface InviteAuthenticationDTO {
  token: string;
}
