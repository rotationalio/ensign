import { Toast } from '@rotational/beacon-core';
import { FormikHelpers } from 'formik';
import { useMemo, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { NewInvitedUserAccount, useRegister } from '@/features/auth';

import NewInviteRegistrationForm from './RegisterNewUser/NewInviteRegistrationForm';
import TeamInvitationCard from './TeamInvitationCard';

export function NewUserInvitationPage({ data }: { data: any }) {
  const [, setIsOpen] = useState(false);
  const register = useRegister();
  const navigateTo = useNavigate();
  const [searchParams] = useSearchParams();

  const invitee_token = searchParams.get('token');
  const handleSubmitRegistration = (
    values: NewInvitedUserAccount,
    helpers: FormikHelpers<NewInvitedUserAccount>
  ) => {
    if (!values.terms_agreement) {
      helpers.setFieldError(
        'terms_agreement',
        'Please agree to terms and conditions before creating Ensign account'
      );
      helpers.setSubmitting(false);
      return;
    }

    const payload = {
      ...values,
      organization: data.org_name,
      // TODO: should be optional on the backend side so we will remove it
      domain: 'N/A',
    };

    register.createNewAccount(payload as any, {
      onSuccess: (_response) => {
        navigateTo('/verify-account', { replace: true });
      },
      onSettled: (_response) => {
        helpers.setSubmitting(false);
      },
    });
  };

  const initialValues: NewInvitedUserAccount = useMemo(
    () => ({
      name: '',
      email: data.email,
      password: '',
      pwcheck: '',
      terms_agreement: false,
      privacy_agreement: false,
      invite_token: invitee_token as string,
    }),
    [data.email, invitee_token]
  );

  const onClose = () => {
    setIsOpen(false);
  };

  return (
    <div>
      <Toast
        isOpen={register.hasAccountFailed}
        onClose={onClose}
        variant="danger"
        description={(register.error as any)?.response?.data?.error}
      />
      <div>
        <div className="mx-auto px-4 pt-8 sm:px-8 md:px-16">
          <TeamInvitationCard data={data} />
        </div>
        <div className="flex flex-col gap-4 px-4 py-8 text-sm sm:p-8 md:flex-row md:p-16 xl:text-base">
          <div className="space-y-4 rounded-md border border-[#1D65A6] bg-[#1D65A6] p-4 text-white sm:p-8 md:w-2/6">
            <h1 className="text-center font-bold">Join the Team</h1>
            <p>
              We designed Ensign to make building event-driven applications fast, convenient, and
              fun! That means working together.
            </p>
            <p>Ensign is great for...</p>
            <ul className="ml-5 list-disc">
              <li>rapid prototying</li>
              <li>real-time analytics</li>
              <li>personalized user experiences</li>
              <li>streaming MLOps pipelines</li>
            </ul>
            <p>Let&apos;s do it team. 💪</p>
          </div>
          <div className="grow rounded-md border border-[#1D65A6] p-4 sm:p-8 md:w-5/6 md:pr-16">
            <div className="mb-4 space-y-3">
              <h2 className="text-base font-bold">Create your Ensign account.</h2>
            </div>
            <NewInviteRegistrationForm
              onSubmit={handleSubmitRegistration}
              initialValues={initialValues}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default NewUserInvitationPage;
