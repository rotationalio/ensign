/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable prettier/prettier */
import { Button, Container, Loader } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { useCheckVerifyToken } from '../hooks/useCheckVerifyToken';

function VerifyPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') as string;

  const { wasVerificationChecked, error, verifyUserEmail, isCheckingToken } =
    useCheckVerifyToken(token);

  useEffect(() => {
    const handleTokenCheck = async () => {
      await verifyUserEmail();
    };
    if (token) {
      handleTokenCheck();
    }
  }, [token, verifyUserEmail]);

  const redirectToLogin = () => {
    window.location.href = '/';
  };

  return (
    <>
      <Container className="my-20">
        <div className="mx-auto max-w-xl rounded-lg border border-solid border-primary-800">
          {isCheckingToken && <Loader />}
          {wasVerificationChecked && (
            <div className="p-10 text-center">
              <div className="ml-8 text-left">
                <p className="mt-6">Thank you for verifying your email address.</p>
                <p className="mt-6">Log in now to start using Ensign.</p>
              </div>
              <Button color="secondary" size="large" className="mt-16" onClick={redirectToLogin}>
                Log in
              </Button>
            </div>
          )}
          {error && (
            <div className="p-10 text-center">
              <div className="ml-8 text-left">
                <p className="mt-6">Sorry, your email can’t be verified.</p>
                <p className="mt-6">Please retry or contact us at support@rotational.io</p>
              </div>
            </div>
          )}
          {!token && (
            <div className="p-10 text-center">
              <div className="ml-8 text-left">
                <p className="mt-6">
                  Sorry we couldn't verify your email address. It looks like you're missing the
                  verification token.
                </p>
              </div>
            </div>
          )}
        </div>
      </Container>
    </>
  );
}

export default VerifyPage;
