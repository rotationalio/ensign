import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import PasswordResetForm from '../components/ForgotPassword/PasswordResetForm';

const PasswordResetPage = () => {
  const submitFormHandler = (values: any) => {
    console.log(values);
  };
  return (
    <div className="px-4 py-8 text-sm sm:p-8 md:flex-row md:p-16 xl:text-base">
      <div className="mx-auto rounded-md border border-[#1D65A6] p-4 sm:p-8 md:w-5/6 md:pr-16">
        <Heading as="h1" className="mb-2 text-lg font-bold">
          <Trans>Password Reset</Trans>
        </Heading>
        <p className="mb-4">
          <Trans>Please enter a new password for your Ensign account.</Trans>
        </p>
        <PasswordResetForm onSubmit={submitFormHandler} />
      </div>
    </div>
  );
};

export default PasswordResetPage;
