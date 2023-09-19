/* eslint-disable prettier/prettier */
import { t, Trans } from '@lingui/macro';
import { Button, Heading } from '@rotational/beacon-core';
import { useEffect } from 'react';
import toast from 'react-hot-toast';
import { Link, useNavigate } from 'react-router-dom';
import styled from 'styled-components';

import { APP_ROUTE } from '@/constants';
import useQueryParams from '@/hooks/useQueryParams';
import { useOrgStore } from '@/store';
import { clearCookies, getCookie } from '@/utils/cookies';
import { decodeToken } from '@/utils/decodeToken';

import LoginForm from '../components/Login/LoginForm';
import { useLogin } from '../hooks/useLogin';
import { isAuthenticated } from '../types/LoginService';

export function Login() {
  const param = useQueryParams();

  const navigate = useNavigate();
  const Store = useOrgStore((state) => state) as any;
  const login = useLogin() as any;

  useEffect(() => {
    if (param?.accountVerified && param?.accountVerified === '1') {
      const isVerified = localStorage.getItem('isEmailVerified');
      if (isVerified === 'true') {
        toast.success(
          t`Thank you for verifying your email address.
          Log in now to start using Ensign.`
        );
        localStorage.removeItem('isEmailVerified');
      }
    }
  }, [param?.accountVerified]);

  useEffect(() => {
    if (!isAuthenticated(login)) {
      clearCookies();
    }
  }, [login]);

  useEffect(() => {
    if (login.authenticated) {
      // set state with the values from the profile
      const token = decodeToken(login?.auth?.access_token) as any;
      Store.setAuthUser(token, !!login.authenticated);

      navigate(APP_ROUTE.DASHBOARD);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [login.authenticated, navigate, login?.auth?.access_token]);

  return (
    <>
      <div className="px-auto mx-auto flex flex-col gap-10 py-8 text-sm sm:p-8 md:flex-row md:justify-center md:p-16 xl:text-base">
        <div className="rounded-md border border-[#1D65A6] p-4 sm:p-8 md:w-[738px] md:pr-16">
          <div className="mb-4 space-y-3">
            <Heading as="h1" className="text-base font-bold">
              <Trans>Log into your Ensign Account.</Trans>
            </Heading>
          </div>
          <LoginForm
            onSubmit={(values: any) => {
              const payload = {
                email: values.email,
                password: values.password,
              } as any;
              if (getCookie('invitee_token')) {
                payload['invite_token'] = getCookie('invitee_token');
              }

              login.authenticate(payload);
            }}
            isDisabled={login.isAuthenticating}
            isLoading={login.isAuthenticating}
          />
        </div>
        <div className="space-y-4 rounded-md border border-[#1D65A6] bg-[#1D65A6] p-4 text-white sm:p-8 md:w-[402px]">
          <h1 className="text-center font-bold">
            <Trans>Need an Account?</Trans>
          </h1>

          <ul className="ml-5 list-disc">
            <li>
              <Trans>Set up your first event stream in minutes</Trans>
            </li>
            <li>
              <Trans>No DevOps foo needed</Trans>
            </li>
            <li>
              <Trans>Goodbye YAML!</Trans>
            </li>
            <li>
              <Trans>We 🤍 SDKs</Trans>
            </li>
            <li>
              <Trans>Learn from beginner-friendly examples</Trans>
            </li>
            <li>
              <Trans>No credit card required</Trans>
            </li>
            <li>
              <Trans>Cancel any time</Trans>
            </li>
          </ul>

          <div className="flex justify-center">
            <Link to="/register">
              <StyledButton
                variant="ghost"
                disabled={login.isAuthenticating}
                className="mt-4"
                data-testid="get__started"
              >
                <Trans>Get Started</Trans>
              </StyledButton>
            </Link>
          </div>
        </div>
      </div>
    </>
  );
}

const StyledButton = styled(Button)((props) => ({
  ...(props.variant === 'ghost' && {
    backgroundColor: 'white!important',
    color: 'rgba(52 58 64)!important',
    border: 'none!important',
    height: 'auto!important',
    width: 'auto!important',
    '&:hover': {
      background: 'rgba(255,255,255, 0.8)!important',
      borderColor: 'rgba(255,255,255, 0.8)!important',
    },
    '&:active': {
      background: 'rgba(255,255,255, 0.8)!important',
      borderColor: 'rgba(255,255,255, 0.8)!important',
    },
  }),
}));

export default Login;
