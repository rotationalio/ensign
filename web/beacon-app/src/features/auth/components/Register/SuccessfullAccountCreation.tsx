import { Trans } from '@lingui/macro';
import { Button, Card, Heading } from '@rotational/beacon-core';

import { ROUTES } from '@/application';
import OtterLookingDown from '@/components/icons/otter-looking-down';

function SuccessfullAccountCreation() {
  const registrationEmail = sessionStorage.getItem('newRegistrationEmail');

  return (
    <div className="relative mx-auto mt-20 w-fit pt-20">
      <OtterLookingDown className="absolute -right-16 -top-[10.8rem]" />
      <Card contentClassName="border border-[#72A2C0] rounded-md p-4 md:p-8 text-sm">
        <Card.Header>
          <h1 className="text-[18px] font-bold">
            <Trans>Let's Go!</Trans>
          </h1>
        </Card.Header>
        <Card.Body>
          <Heading as="h1" className="mt-4 mb-3 ">
            <Trans>
              To keep your account safe, we sent a verification email to{' '}
              {registrationEmail ? (
                <span className="font-bold">{registrationEmail}</span>
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
              <a href={`mailto:${ROUTES.SUPPORT}`} className="font-bold text-[#1F4CED]">
                support@rotational.io
              </a>
              .
            </Trans>
          </p>
          <Button
            variant="primary"
            onClick={() => console.log('Resend verification message!')}
            className="mt-4 font-bold text-white"
          >
            Resend verification email
          </Button>{' '}
        </Card.Body>
      </Card>
    </div>
  );
}

export default SuccessfullAccountCreation;
