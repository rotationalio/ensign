/* eslint-disable prettier/prettier */
import { Heading } from '@rotational/beacon-core';
import { Link, useNavigate } from 'react-router-dom';

import Button from '@/components/ui/Button';
import { Toast } from '@/components/ui/Toast';
import { APP_ROUTE } from '@/constants';
import { useOrgStore } from '@/store';
import { decodeToken } from '@/utils/decodeToken';

import LoginForm from '../components/Login/LoginForm';
import { useLogin } from '../hooks/useLogin';
import { isAuthenticated } from '../types/LoginService';

export function Login() {
  const navigate = useNavigate();
  useOrgStore.persist.clearStorage();
  const login = useLogin() as any;

  if (isAuthenticated(login)) {
    const token = decodeToken(login.auth.access_token) as any;
    //console.log('token', token);

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
      <div className="px-auto mx-auto flex flex-col gap-10 py-8 text-sm sm:p-8 md:flex-row md:justify-center md:p-16 xl:text-base">
        <div className="rounded-md border border-[#1D65A6] p-4 sm:p-8 md:w-[738px] md:pr-16">
          <div className="mb-4 space-y-3">
            <Heading as="h1" className="text-base font-bold">
              Log into your Ensign Account.
            </Heading>
          </div>
          <LoginForm
            onSubmit={login.authenticate}
            isDisabled={login.isAuthenticating}
            isLoading={login.isAuthenticating}
          />
        </div>
        <div className="space-y-4 rounded-md border border-[#1D65A6] bg-[#1D65A6] p-4 text-white sm:p-8 md:w-[402px]">
          <h1 className="text-center font-bold">Need an Account?</h1>

          <ul className="ml-5 list-disc">
            <li>Set up your first event stream in minutes</li>
            <li>No DevOps foo needed</li>
            <li>Goodbye YAML!</li>
            <li>We 🤍 SDKs </li>
            <li>Learn from beginner-friendly examples</li>
            <li>No credit card required</li>
            <li>Cancel any time</li>
          </ul>

          <div className="flex justify-center">
            <Link to="/register" className="btn btn-primary ">
              <Button
                isDisabled={login.isAuthenticating}
                className="mt-4 bg-white text-gray-800"
                data-testid="get__started"
              >
                Get Started
              </Button>
            </Link>
          </div>
        </div>
      </div>
    </>
  );
}

export default Login;
