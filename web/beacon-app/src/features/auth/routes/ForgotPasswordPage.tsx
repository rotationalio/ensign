import { Trans } from '@lingui/macro';
import { Card } from '@rotational/beacon-core';

import OtterLookingDown from '@/components/icons/otter-looking-down';

// import useForGetPasswordMutation from '@/hooks/useForgetPasswordMutation';
import ForgotPasswordForm from '../components/ForgotPassword/ForgotPasswordForm';

const ForgotPasswordPage = () => {
  // const { forgetPasswordRequest, isLoading } = useForGetPasswordMutation();
  const submitFormHandler = (values: any) => {
    console.log(values);
    // forgetPasswordRequest(values.email);
  };
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
          <ForgotPasswordForm onSubmit={submitFormHandler} />
        </Card.Body>
      </Card>
    </div>
  );
};

export default ForgotPasswordPage;
