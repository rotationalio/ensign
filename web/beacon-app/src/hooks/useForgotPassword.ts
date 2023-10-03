import { useState } from 'react';

import axiosInstance from '@/application/api/ApiService';

const useForGotPasswordMutation = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [data, setData] = useState(null);
  const [isSuccess, setIsSuccess] = useState(false);

  const forgetPasswordRequest = async (email: string) => {
    try {
      setIsLoading(true);
      const response = await axiosInstance.post(`/forgot-password`, {
        email,
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
    forgetPasswordRequest,
    isLoading,
    error,
    data,
    isSuccess,
    reset,
  };
};

export default useForGotPasswordMutation;
