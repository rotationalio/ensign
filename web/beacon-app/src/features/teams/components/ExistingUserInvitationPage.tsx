import { Heading, Toast } from '@rotational/beacon-core';
import { Link, useNavigate } from 'react-router-dom';

import { APP_ROUTE } from '@/constants';
import { isAuthenticated, useLogin } from '@/features/auth';
import { LoginForm } from '@/features/auth/components';
import { useOrgStore } from '@/store';
import { decodeToken } from '@/utils/decodeToken';

import TeamInvitationCard from './TeamInvitationCard';

export default function ExistingUserInvitationPage() {
  const navigate = useNavigate();
  useOrgStore.persist.clearStorage();
  const login = useLogin() as any;

  if (isAuthenticated(login)) {
    const token = decodeToken(login.auth.access_token) as any;
    // console.log('token', token)

    useOrgStore.setState({
      org: token?.org,
      user: token?.sub,
      isAuthenticated: !!login.authenticated,
      name: token?.name,
      email: token?.email,
      picture: token?.picture,
      permissions: token?.permissions,
    });

    // if(!login.auth?.last_login){
    //   navigate(APP_ROUTE.GETTING_STARTED);
    // }
    // else{
    navigate(APP_ROUTE.DASHBOARD);
    //}
  }
  return (
    <>
      {login.hasAuthFailed && (
        <Toast
          isOpen={login.hasAuthFailed}
          variant="danger"
          description={(login.error as any)?.response?.data?.error}
        />
      )}
      <div className="mx-auto pt-8 sm:px-9 md:px-16 2xl:px-40">
        <TeamInvitationCard />
      </div>
      <div className="px-auto mx-auto flex flex-col gap-10 py-8 text-sm sm:p-8 md:flex-row md:justify-center md:p-16 xl:text-base">
        <div className="space-y-4 rounded-md border border-[#1D65A6] bg-[#1D65A6] p-4 text-white sm:p-8 md:w-[402px]">
          <h1 className="text-center font-bold">Join the Team</h1>
          <p>
            Log in to your existing account to accept the invitation and start working with your
            teammates!
          </p>
        </div>
        <div className="rounded-md border border-[#1D65A6] p-4 sm:p-8 md:w-[738px] md:pr-16">
          <div className="mb-4 space-y-3">
            <Heading as="h1" className="text-base font-bold">
              Log into your Ensign Account.
            </Heading>
            <p>
              Log in to accept. Don't have an account?{' '}
              <Link to="/new-invitation" className="font-semibold text-[#1D65A6]">
                Create one now.
              </Link>
            </p>
          </div>
          <LoginForm
            onSubmit={login.authenticate}
            /* TODO: Make button disabled until form is filled */
            isDisabled={login.isAuthenticating}
            isLoading={login.isAuthenticating}
          />
        </div>
      </div>
    </>
  );
}