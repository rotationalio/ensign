import { AriaButton as Button } from '@rotational/beacon-core';
import { memo, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';

import { ROUTES } from '@/application';
import Logo from '@/components/common/Logo';

function LandingHeader() {
  const location = useLocation();
  const isConfirmationPage = location.pathname == ROUTES.VERIFY_PAGE;

  return (
    <nav className="border-b border-primary-800 py-8">
      <div className="container mx-auto flex flex-wrap items-center justify-between">
        <Logo />
        <div className="space-x-8">
          {isConfirmationPage && (
            <Link to="/">
              <Button
                data-testid="login-button"
                color="secondary"
                className="mt-4 min-w-[100px] py-2"
                aria-label="Log in"
              >
                Log in
              </Button>
            </Link>
          )}
          <Link to="/" className="font-bold capitalize text-primary">
            Starter Plan
          </Link>
          {/* <Link to="/">
            <Button variant="tertiary" size="small">
              Upgrade
            </Button>
          </Link> */}
        </div>
      </div>
    </nav>
  );
}

export default memo(LandingHeader);
