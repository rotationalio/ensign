import { t, Trans } from '@lingui/macro';
import { Card } from '@rotational/beacon-core';
import { useEffect } from 'react';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';

import { ROUTES } from '@/application';
import OtterLookingDown from '@/components/icons/otter-looking-down';

// import useForGetPasswordMutation from '@/hooks/useForgetPasswordMutation';
import ForgotPasswordForm from '../components/Form/ForgotPasswordForm';
import { useForgotPassword } from '../hooks/useForgotPassword';

const ForgotPasswordPage = () => {
  const navigate = useNavigate();
  const {
    forgotPassword,
    wasForgotPasswordSuccessful,
    hasForgotPasswordFailed,
    isLoading,
    error,
    reset,
  } = useForgotPassword();
  const submitFormHandler = (values: any) => {
    const payload = {
      email: values.email,
    };
    forgotPassword(payload);
  };

  useEffect(() => {
    if (wasForgotPasswordSuccessful) {
      navigate(ROUTES.RESET_VERIFICATION);
      reset();
    }
  }, [wasForgotPasswordSuccessful, navigate, reset]);

  useEffect(() => {
    if (hasForgotPasswordFailed) {
      toast.error(
        error?.response?.data?.error ||
          t`Unable to submit forgot password request. Please try again or contact support, if the problem continues.`
      );
      reset();
    }
  }, [hasForgotPasswordFailed, error, reset]);

  return (
    <div className="relative mx-auto mt-20 w-fit pt-20">
      <OtterLookingDown className="absolute -right-16 -top-[10.8rem]" />
      <Card contentClassName="border border-[#72A2C0] rounded-md p-4 md:p-8 text-sm">
        <Card.Body>
          <p className="mb-4">
            <Trans>
              Forgot your password? No problem. Enter your email address to recover your login
              credentials.
            </Trans>
          </p>
          <ForgotPasswordForm onSubmit={submitFormHandler} isSubmitting={isLoading} />
        </Card.Body>
      </Card>
    </div>
  );
};

export default ForgotPasswordPage;
