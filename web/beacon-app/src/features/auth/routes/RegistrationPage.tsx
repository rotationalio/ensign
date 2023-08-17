import { FormikHelpers } from 'formik';
import { Link, useNavigate } from 'react-router-dom';

// import { slugify } from '@/utils/slugifyDomain';
import RegistrationForm from '../components/Register/RegistrationForm';
import { useRegister } from '../hooks/useRegister';
import { NewUserAccount } from '../types/RegisterService';

export function Registration() {
  const register = useRegister();
  const navigateTo = useNavigate();
  const handleSubmitRegistration = (
    values: NewUserAccount,
    helpers: FormikHelpers<NewUserAccount>
  ) => {
    if (!values.terms_agreement) {
      helpers.setFieldError(
        'terms_agreement',
        'Please agree to the terms and conditions before creating your Ensign account.'
      );
      helpers.setSubmitting(false);
      return;
    }

    register.createNewAccount(values, {
      onSuccess: (_response) => {
        navigateTo('/verify-account', { replace: true });
      },
      onSettled: (_response) => {
        helpers.setSubmitting(false);
      },
    });
  };

  return (
    <>
      {/* {register.hasAccountFailed && (
        <Toast
          isOpen={register.hasAccountFailed}
          variant="danger"
          description={(register.error as any)?.response?.data?.error}
        />
      )} */}

      <div className="flex flex-col gap-4 px-4 py-8 text-sm sm:p-8 md:flex-row md:p-16 xl:text-base">
        <div className="grow rounded-md border border-[#1D65A6] p-4 sm:p-8 md:w-5/6 md:pr-16">
          <div className="mb-4 space-y-3">
            <h2 className="text-base font-bold">Create your starter Ensign account.</h2>
            <p>
              Already have an account?{' '}
              <Link to="/" className="font-semibold text-[#1d65a6]">
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
