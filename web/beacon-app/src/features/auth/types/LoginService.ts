// import { UseMutateFunction } from '@tanstack/react-query';

import { User } from './RegisterService';
export interface UserAuthResponse {
  access_token: string;
  refresh_token: string;
}

export interface LoginMutation {
  authenticate: (user: AuthUser) => void;
  reset: () => void;
  auth: UserAuthResponse;
  error: any;
  isAuthenticating: boolean;
  authenticated: boolean;
  hasAuthFailed: boolean;
}

export type AuthUser = Pick<User, 'email' | 'password'>;

export const isAuthenticated = (mutation: LoginMutation): mutation is Required<LoginMutation> =>
  mutation.authenticated && mutation.auth !== undefined;
