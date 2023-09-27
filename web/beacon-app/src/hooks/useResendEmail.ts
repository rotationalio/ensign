import { t } from '@lingui/macro';
import { useCallback, useEffect, useState } from 'react';
import toast from 'react-hot-toast';

import { appConfig } from '@/application/config';
const useResendEmail = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);
  const [result, setResult] = useState(false);

  const request = async (email: string) => {
    try {
      const response = await fetch(`${appConfig.tenantApiUrl}/resend`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });
      // send the status code to the client
      const status = response.status;

      return status;
    } catch (error) {
      console.log(error);
    }
  };

  const resendEmail = async (email: string) => {
    setLoading(true);
    try {
      const response = await request(email);
      // console.log('[] response', response);
      // if response status is 204, it means that the email was sent
      if (response === 204) {
        setResult(true);
        setLoading(false);
      }
    } catch (error) {
      setError(true);
      setLoading(false);
    }
  };
  const reset = useCallback(() => {
    setLoading(false);
    setError(false);
    setResult(false);
  }, []);

  useEffect(() => {
    if (result) {
      // close other toast if any
      toast.dismiss();
      toast.success(
        t`Verification email sent. Please check your inbox and follow the instructions.`
      );
    }
    return () => {
      // clear resend result
      result && reset();
    };
  }, [result, reset]);

  return { loading, error, result, resendEmail, reset };
};

export default useResendEmail;
