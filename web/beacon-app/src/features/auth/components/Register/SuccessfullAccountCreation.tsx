import { Trans } from '@lingui/macro';
import { Card } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import ChevronInCircle from '@/components/icons/chevron-in-circle';
import OtterLookingDown from '@/components/icons/otter-looking-down';

function SuccessfullAccountCreation() {
  return (
    <div className="relative mx-auto mt-10 w-fit pt-20">
      <OtterLookingDown className="absolute -right-16 -top-[8.8rem]" />
      <Card contentClassName="border border-primary-900 rounded-md p-4 md:p-8 text-sm">
        <Card.Header>
          <h1 className="text-base font-bold">
            <Trans>Thank you for creating your Ensign account!</Trans>
          </h1>
        </Card.Header>
        <Card.Body>
          <p className="mt-4 mb-3">
            <Trans>Please check your email to verify your account.</Trans>
          </p>
          <div className="space-y-2">
            <p className="font-semibold">
              <Trans>Next Steps</Trans>
            </p>
            <ul className="space-y-2">
              <li className="flex items-center gap-2">
                <ChevronInCircle />
                <Link to="" className="underline">
                  <Trans>Read the documentation</Trans>
                </Link>
              </li>
              <li className="flex items-center gap-2">
                <ChevronInCircle />
                <Link to="" className="underline">
                  <Trans>Checkout a tutorial</Trans>
                </Link>
              </li>
              <li className="flex items-center gap-2">
                <span>
                  <ChevronInCircle />
                </span>
                <Link to="" className="underline">
                  <Trans>Learn how sea otters are more than just cute (Wait, what?)</Trans>
                </Link>
              </li>
            </ul>
          </div>
        </Card.Body>
      </Card>
    </div>
  );
}

export default SuccessfullAccountCreation;
