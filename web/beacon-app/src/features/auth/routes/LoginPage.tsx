/* eslint-disable prettier/prettier */
import { useState } from 'react';
import toast, { Toaster } from 'react-hot-toast';
import { Link } from 'react-router-dom';
import { AriaButton as Button, Heading, Toast } from '@rotational/beacon-core';

import { decodeToken } from '@/utils/decodeToken';

import LoginForm from '../components/Login/LoginForm';
import { useLogin } from '../hooks/useLogin';
import { isAuthenticated } from '../types/LoginService';
export function Login() {
  const [, setIsOpen] = useState(false);
  //const navigate = useNavigate();
  const login = useLogin() as any;

  const onClose = () => {
    setIsOpen(false);
  };

  if (isAuthenticated(login)) {
    console.log('called');
    console.log(decodeToken(login.auth.access_token));
    toast.success('Login successful', {
      duration: 5000,
      position: 'top-right',
      className: 'w-[300px] h-[50px]',
    });
    // setTimeout(() => {
    //   navigate('/dashboard');
    // }, 5000);
  }

  return (
    <>
      {login.hasAuthFailed && (
        <Toast
          isOpen={login.hasAuthFailed}
          onClose={onClose}
          variant="danger"
          title="Something went wrong, please try again later."
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

          <LoginForm onSubmit={login.authenticate} isDisabled={login.isAuthenticating} />
        </div>
        <div className="space-y-4 rounded-md border border-[#1D65A6] bg-[#1D65A6] p-4 text-white sm:p-8 md:w-[402px]">
          <h1 className="text-center font-bold">Need an Account ?</h1>

          <ul className="ml-5 list-disc">
            <li>new prototypes without refactoring legacy database schemas</li>
            <li>real-time dashboards and analytics in days rather than months?</li>
            <li>rich, tailored experiences so your users knows how much they means to you?</li>
            <li>MLOps pipelines that bridge the gap between the training and deployment phases?</li>
          </ul>

          <div className="flex justify-center">
            <Link to="/register" className="btn btn-primary ">
              {' '}
              <Button color="secondary" className="text-gray-800 mt-4 bg-white">
                Create Account{' '}
              </Button>
            </Link>
          </div>
          <Toaster />
        </div>
      </div>
    </>
  );
}

export default Login;
