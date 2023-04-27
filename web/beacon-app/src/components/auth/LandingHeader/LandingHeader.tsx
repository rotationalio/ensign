import { Button } from '@rotational/beacon-core';
import { memo } from 'react';
import { Link, useLocation } from 'react-router-dom';

import { ROUTES } from '@/application';
import Logo from '@/components/common/Logo';

function LandingHeader() {
  const location = useLocation();
  const isConfirmationPage = location.pathname == ROUTES.VERIFY_PAGE;
  const isEmailConfirmationPage = location.pathname == ROUTES.VERIFY_EMAIL;

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
          {isEmailConfirmationPage && (
            <Link to="/register">
              <Button
                data-testid="registration-button"
                color="primary"
                className="mt-4 font-bold"
                size="medium"
                aria-label="Get started"
              >
                Get started
              </Button>
            </Link>
          )}
          <Link to="/" className="font-bold capitalize text-primary">
            Ensign Beta
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
