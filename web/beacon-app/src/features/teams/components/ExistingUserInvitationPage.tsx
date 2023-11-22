import { Heading } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { APP_ROUTE } from '@/constants';
import { useLogin } from '@/features/auth';
import { LoginForm } from '@/features/auth/components';
import { InviteAuthUser, isAuthenticated } from '@/features/auth/types/LoginService';
import { useOrgStore } from '@/store';
import { getCookie } from '@/utils/cookies';
import { decodeToken } from '@/utils/decodeToken';

import useFetchInviteAuthentication from '../hooks/useFetchInviteAuthentication';
import TeamInvitationCard from './TeamInvitationCard';

export default function ExistingUserInvitationPage({ data }: { data: any }) {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const Store = useOrgStore((state) => state) as any;

  const invitee_token = searchParams.get('token');

  useOrgStore.persist.clearStorage();
  const login = useLogin() as any;
  const initialValues = {
    email: data?.email || '',
    password: '',
    invite_token: invitee_token,
  } as InviteAuthUser;

  const { invitationRequest, authData, wasInvitationAuthenticated } = useFetchInviteAuthentication(
    invitee_token as string
  );

  const authenticatedUser = getCookie('authenticatedUser') === 'true';

  // If the existing invited user is authenticated, obtain the access token from the invitation
  // and redirect the user their new team's dashboard.
  useEffect(() => {
    if (authenticatedUser) {
      invitationRequest(invitee_token);
      if (wasInvitationAuthenticated && authData?.access_token) {
        const token = decodeToken(authData?.access_token) as any;
        Store.setAuthUser(token, !!authData?.access_token);

        navigate(APP_ROUTE.DASHBOARD);        
      }
    }
  }, [
    authenticatedUser,
    invitationRequest,
    invitee_token,
    navigate,
    Store,
    authData,
    wasInvitationAuthenticated,
  ]);

  // If the existing invited user is not authenticated, the user will have to submit
  // the log in form before they're redirected to their new team's dashboard.
  if (!authenticatedUser) {
    if (isAuthenticated(login)) {
      const token = decodeToken(login.auth.access_token) as any;
      Store.setAuthUser(token, !!login.authenticated);

      navigate(APP_ROUTE.DASHBOARD);
    }
  }

  return (
    <div>
      <div className="pt-8 sm:px-9 md:px-16 2xl:px-40">
        <TeamInvitationCard data={data} />
      </div>
      <div className="sm:px-auto mx-auto flex flex-col gap-10 pt-8 pb-8 text-sm sm:p-8 md:flex-row md:justify-center md:px-16 md:py-8 xl:text-base">
        <div className="space-y-4 rounded-md border border-[#1D65A6] bg-[#1D65A6] p-4 text-white sm:p-8 md:w-[402px]">
          <h1 className="text-center font-bold">Join the Team</h1>
          <p>
            Log in to your existing account to accept the invitation and start working with your
            teammates!
          </p>
        </div>
        <div className="rounded-md border border-[#1D65A6] p-4 sm:p-8 md:w-[790px] md:pr-16">
          <div className="mb-4 space-y-3">
            <Heading as="h1" className="text-base font-bold">
              Log into your Ensign Account.
            </Heading>
          </div>
          <LoginForm
            onSubmit={login.authenticate}
            initialValues={initialValues}
            isDisabled={login.isAuthenticating}
            isLoading={login.isAuthenticating}
          />
        </div>
      </div>
    </div>
  );
}
