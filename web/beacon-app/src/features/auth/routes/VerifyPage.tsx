/* eslint-disable prettier/prettier */
import { Card, Container } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { useCheckVerifyToken } from '../hooks/useCheckVerifyToken';

function VerifyPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') as string;

  const { data, wasVerificationChecked, error, verifyUserEmail } = useCheckVerifyToken(token);

  useEffect(() => {
    if (!data) {
      verifyUserEmail();
    }
  }, [token, verifyUserEmail, data]);

  console.log('[data]', data);

  return (
    <>
      <Container>
        {wasVerificationChecked && (
          <Card className="rounded-lg border border-solid border-primary-800 text-2xl">
            <p className="mx-auto mt-8 max-w-xl">
              Thank you for verifying your email address. Log in now to start using Ensign.
            </p>
          </Card>
        )}
        {wasVerificationChecked && !data?.verified && <div>Not Verified</div>}
        {error && <div>Something went wrong</div>}
      </Container>
    </>
  );
}

export default VerifyPage;
