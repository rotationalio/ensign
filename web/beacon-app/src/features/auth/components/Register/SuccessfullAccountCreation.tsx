import { Trans } from '@lingui/macro';
import { Button, Card, Heading } from '@rotational/beacon-core';
import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { EXTERNAL_LINKS, ROUTES } from '@/application';
import OtterLookingDown from '@/components/icons/otter-looking-down';
import useResendEmail from '@/hooks/useResendEmail';

function SuccessfullAccountCreation() {
  const navigateTo = useNavigate();

  const storage = localStorage.getItem('esg.new.user');
  const [userEmail, setUserEmail] = useState<string | null>(storage);
  const { resendEmail } = useResendEmail();

  // console.log('userEmail', userEmail);

  useEffect(() => {
    if (userEmail) {
      setUserEmail(userEmail);
    } else {
      navigateTo(ROUTES.LOGIN);
    }

    return () => {
      localStorage.removeItem('esg.new.user');
    };
  }, [userEmail, navigateTo, storage]);

  const resendEmailHandler = useCallback(() => {
    resendEmail(userEmail as string);
  }, [userEmail, resendEmail]);

  return (
    <div className="relative mx-auto mt-20 w-fit pt-20">
      {userEmail && (
        <>
          <OtterLookingDown className="absolute -right-16 -top-[10.8rem]" />
          <Card contentClassName="border border-[#72A2C0] rounded-md p-4 md:p-8 text-sm">
            <Card.Header>
              <h1 className="text-[18px] font-bold">
                <Trans>Let's Go!</Trans>
              </h1>
            </Card.Header>
            <Card.Body>
              <Heading as="h1" className="mb-3 mt-4 ">
                <Trans>
                  To keep your account safe, we sent a verification email to{' '}
                  {userEmail ? (
                    <span className="font-bold" data-cy="registration-email">
                      {userEmail}
                    </span>
                  ) : (
                    'the email address provided during sign up'
                  )}
                  . Click the secure link in that email to continue.
                </Trans>
              </Heading>
              <p>
                {' '}
                <Trans>
                  Didn't receive an email?{' '}
                  <span className="font-bold text-[#1F4CED]">Resend the verification email</span> or
                  email{' '}
                  <a href={`mailto:${EXTERNAL_LINKS.SUPPORT}`} className="font-bold text-[#1F4CED]">
                    support@rotational.io
                  </a>
                  .
                </Trans>
              </p>
              <Button
                variant="primary"
                onClick={resendEmailHandler}
                className="mt-4 font-bold text-white"
              >
                Resend verification email
              </Button>{' '}
            </Card.Body>
          </Card>
        </>
      )}
    </div>
  );
}

export default SuccessfullAccountCreation;
