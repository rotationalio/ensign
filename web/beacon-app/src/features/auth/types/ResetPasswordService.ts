import { UseMutateFunction } from '@tanstack/react-query';

export interface ResetPasswordMutation {
  resetPassword: UseMutateFunction<unknown, unknown, ResetPasswordDTO, unknown>;
  reset: () => void;
  hasResetPasswordFailed: boolean;
  wasResetPasswordSuccessful: boolean;
  isLoading: boolean;
  error: any;
}

export interface ResetPasswordDTO {
  token: string;
  password: string;
  pwcheck: string;
}
