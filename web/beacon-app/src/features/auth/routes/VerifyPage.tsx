/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable prettier/prettier */
import { Container, Loader } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { useCheckVerifyToken } from '../hooks/useCheckVerifyToken';

function VerifyPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') as string;
  const navigate = useNavigate();
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

  useEffect(() => {
    if (wasVerificationChecked && !error) {
      localStorage.setItem('isEmailVerified', 'true');
      navigate('/?accountVerified=1');
    }
  }, [wasVerificationChecked, error, navigate]);

  return (
    <>
      <Container className="my-20">
        <div className="mx-auto min-h-min max-w-xl rounded-lg border border-solid border-primary-800">
          {isCheckingToken && (
            <div className="items-center justify-center p-10 text-center">
              <Loader />
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
