import { Trans } from '@lingui/macro';
import { Card, Heading } from '@rotational/beacon-core';

import OtterLookingDown from '@/components/icons/otter-looking-down';
import OnboardingStepper from '@/features/onboarding/components/OnboardingStepper';

function SuccessfullAccountCreation() {
  return (
    <div className="relative mx-auto mt-20 w-fit pt-20">
      <OtterLookingDown className="absolute -right-16 -top-[10.8rem]" />
      <Card contentClassName="border border-[#72A2C0] rounded-md p-4 md:p-8 text-sm">
        <Card.Header>
          <h1 className="text-[18px] font-bold">
            <Trans>Thank you for creating your Ensign account!</Trans>
          </h1>
        </Card.Header>
        <Card.Body>
          <Heading as="h1" className="mt-4 mb-3 ">
            <Trans>
              We value security and care deeply about protecting your information. For that reason,
              we request that you verify your email account.{' '}
            </Trans>
          </Heading>
          <div className="space-y-2">
            <Heading as="h2" className="text-md font-bold">
              <Trans>Next Steps</Trans>
            </Heading>
            <p>
              {' '}
              <Trans>
                Please check your email and click the secure link in the verification email we just
                sent you. You can then log into Ensign to start building!
              </Trans>
            </p>
          </div>
        </Card.Body>
      </Card>
      <OnboardingStepper />
    </div>
  );
}

export default SuccessfullAccountCreation;
