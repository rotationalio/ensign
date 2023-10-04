import { UseMutateFunction } from '@tanstack/react-query';

export interface ForgotPasswordMutation {
  forgotPassword: UseMutateFunction<unknown, unknown, ForgotPasswordDTO, unknown>;
  reset: () => void;
  hasForgotPasswordFailed: boolean;
  wasForgotPasswordSuccessful: boolean;
  isLoading: boolean;
  error: any;
}

export interface ForgotPasswordDTO {
  email: string;
}
