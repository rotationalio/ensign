import { t } from '@lingui/macro';
import { useEffect } from 'react';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';

import { APP_ROUTE } from '@/constants';

import { ResetPasswordDTO } from '../types/ResetPasswordService';
import { useResetPassword } from './useResetPassword';
export const useSubmitResetPassword = () => {
  const navigate = useNavigate();
  const {
    resetPassword: resetPasswordMutation,
    reset,
    hasResetPasswordFailed,
    isLoading,
    error,
    wasResetPasswordSuccessful,
  } = useResetPassword();
  // has errored with 400 status code
  const hasErrored = error && error?.response?.status === 400;

  const resetPassword = ({ token, password, pwcheck }: ResetPasswordDTO) => {
    resetPasswordMutation({
      token,
      password,
      pwcheck,
    });
  };

  useEffect(() => {
    if (wasResetPasswordSuccessful) {
      reset();
      navigate(`${APP_ROUTE.ROOT}?from=reset-password`);
    }
  }, [wasResetPasswordSuccessful, reset, navigate]);

  // handle toast error from failed reset password

  useEffect(() => {
    if (hasErrored) {
      toast.error(error?.response?.data?.error || t`Something went wrong. Please try again.`);
    }
    return () => {
      reset();
    };
  }, [hasErrored, error, reset]);

  return { resetPassword, hasResetPasswordFailed, isLoading };
};
