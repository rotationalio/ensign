import { Toast } from '@rotational/beacon-core';
import { FormikHelpers } from 'formik';
import { Link, useNavigate } from 'react-router-dom';

import slugify from '@/utils/slugifyDomain';

import RegistrationForm from '../components/RegistrationForm';
import { useRegister } from '../hooks/useRegister';
import { NewUserAccount } from '../types/RegisterService';

export function Registration() {
  const register = useRegister();
  const navigateTo = useNavigate();
  const handleSubmitRegistration = (
    values: NewUserAccount,
    helpers: FormikHelpers<NewUserAccount>
  ) => {
    values.domain = slugify(values.domain);

    register.createNewAccount(values, {
      onSuccess: (_response) => {
        navigateTo('/auth/verify-account', { replace: true });
      },
      onSettled: (_response) => {
        helpers.setSubmitting(false);
      },
    });
  };

  return (
    <>
      <Toast
        isOpen={register.hasAccountFailed}
        title="Something went wrong, please try again later."
        description={(register.error as any)?.response?.data?.error}
      />
      <div className="flex flex-col gap-4 px-4 py-8 text-sm sm:p-8 md:flex-row md:p-16 xl:text-base">
        <div className="space-y-4 rounded-md border border-[#1D65A6] bg-[#1D65A6] p-4 text-white sm:p-8 md:w-2/6">
          <h1 className="text-center font-bold">
            Building event-driven applications can be fast, convenient and even fast!
          </h1>
          <p className="text-center font-bold">Start today on our no-cost Starter Plan</p>
          <p>
            If you have always wanted to try out eventing, but couldn&apos;t justify the hight cost
            of entry or the expertise required, Ensign is for you!
          </p>
          <p>Want to build...</p>
          <ul className="ml-5 list-disc">
            <li>new prototypes without refactoring legacy database schemas</li>
            <li>real-time dashboards and analytics in days rather than months?</li>
            <li>rich, tailored experiences so your users knows how much they means to you?</li>
            <li>MLOps pipelines that bridge the gap between the training and deployment phases?</li>
          </ul>
          <p>Let&apos;s do it hero 💪</p>
        </div>
        <div className="grow rounded-md border border-[#1D65A6] p-4 sm:p-8 md:w-5/6 md:pr-16">
          <div className="mb-4 space-y-3">
            <h2 className="text-base font-bold">Create your starter Ensign account.</h2>
            <p>
              Already have an account?{' '}
              <Link to="/signin" className="font-semibold text-[#1d65a6]">
                Skip the line and just sign in
              </Link>
              .
            </p>
          </div>
          <RegistrationForm onSubmit={handleSubmitRegistration} />
        </div>
      </div>
    </>
  );
}

export default Registration;
