import { useState } from 'react';

import axiosInstance from '@/application/api/ApiService';

export interface ResetPasswordMutation {
  resetPasswordRequest: (data: ResetPasswordDTO) => Promise<void>;
  isLoading: boolean;
  error: any;
  data: any;
  isSuccess: boolean;
  reset: () => void;
}

export interface ResetPasswordDTO {
  token: string;
  password: string;
  pwcheck: string;
}

const useResetPasswordMutation = (): ResetPasswordMutation => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [data, setData] = useState(null);
  const [isSuccess, setIsSuccess] = useState(false);

  const resetPasswordRequest = async (data: ResetPasswordDTO) => {
    try {
      setIsLoading(true);
      const response = await axiosInstance.post(`/reset-password`, {
        data,
      });

      if (response.status === 200 || response.status === 204) {
        setIsSuccess(true);
      }

      setData(response.data);
    } catch (error: any) {
      setError(error);
    } finally {
      setIsLoading(false);
    }
  };

  const reset = () => {
    setIsLoading(false);
    setError(null);
    setData(null);
    setIsSuccess(false);
  };

  return {
    resetPasswordRequest,
    isLoading,
    error,
    data,
    isSuccess,
    reset,
  };
};

export default useResetPasswordMutation;
