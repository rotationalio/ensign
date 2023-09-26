/* eslint-disable prettier/prettier */
import { Trans } from '@lingui/macro';
import { Button, Heading } from '@rotational/beacon-core';
import { useCallback, useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import styled from 'styled-components';

import LoginFooter from '@/features/auth/components/LoginFooter';
import useQueryParams from '@/hooks/useQueryParams';
import useResendEmail from '@/hooks/useResendEmail';
import { clearSessionStorage } from '@/utils/cookies';

import LoginForm from '../components/Login/LoginForm';
import { useLogin } from '../hooks/useLogin';
import useSubmitLogin from '../hooks/useSubmitLogin';
import useDisplayVerfiedAccountToast from '../hooks/useVerifiedAccount';
import { isAuthenticated } from '../types/LoginService';
export function Login() {
  const param = useQueryParams();

  const [currentUserEmail, setCurrentUserEmail] = useState('');
  const { authenticate, error, isAuthenticating } = useLogin() as any;
  const { resendEmail, reset } = useResendEmail();
  useDisplayVerfiedAccountToast(param);
  // console.log('[] resendResult', resendResult);
  const { onSubmitHandler } = useSubmitLogin({
    setData: setCurrentUserEmail,
    onReset: reset,
    onSetCurrentUserEmail: setCurrentUserEmail,
  });

  const resendEmailHandler = useCallback(() => {
    resendEmail(currentUserEmail);
  }, [currentUserEmail, resendEmail]);

  useEffect(() => {
    if (!isAuthenticated(authenticate)) {
      clearSessionStorage();
    }
  }, [authenticate]);

  useEffect(() => {
    if (error && error.response.status === 403) {
      toast.error(
        <div className="flex flex-col gap-5">
          <p>
            <Trans>Please verify your email address and try again!</Trans>
          </p>
          <div>
            <Button size="small" className="max-w-40 " onClick={resendEmailHandler}>
              <Trans>Resend Email</Trans>
            </Button>
          </div>
        </div>
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [error]);

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
            onSubmit={onSubmitHandler}
            isDisabled={isAuthenticating}
            isLoading={isAuthenticating}
          />
        </div>
        <LoginFooter ButtonElement={StyledButton} isAuthenticating={isAuthenticating} />
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
