import { Card } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import ChevronInCircle from '@/components/icons/chevron-in-circle';
import OtterLookingDown from '@/components/icons/otter-looking-down';

function SuccessfullAccountCreation() {
  return (
    <div className="relative w-fit pt-20">
      <OtterLookingDown className="absolute -right-16 -top-[8.8rem]" />
      <Card contentClassName="border border-primary-900 rounded-md p-4 md:p-8 text-sm">
        <Card.Header>
          <h1 className="text-base font-bold">Thank you for creating your Ensign account!</h1>
        </Card.Header>
        <Card.Body>
          <p className="mt-4 mb-3">Please check your email to verify your account.</p>
          <div className="space-y-2">
            <p className="font-semibold">Next Steps</p>
            <ul className="space-y-2">
              <li className="flex items-center gap-2">
                <ChevronInCircle />
                <Link to="" className="underline">
                  Read the documentation
                </Link>
              </li>
              <li className="flex items-center gap-2">
                <ChevronInCircle />
                <Link to="" className="underline">
                  Checkout a tutorial
                </Link>
              </li>
              <li className="flex items-center gap-2">
                <span>
                  <ChevronInCircle />
                </span>
                <Link to="" className="underline">
                  Learn how sea otters are more than just cute (Wait, what?)
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
