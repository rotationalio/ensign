import { Trans } from '@lingui/macro';
import { memo } from 'react';
import { Link } from 'react-router-dom';

import Logo from '@/components/common/Logo';

function LandingHeader() {
  return (
    <nav className="border-b border-primary-800 py-8">
      <div className="container mx-auto flex flex-wrap items-center justify-between">
        <Logo />
        <div className="space-x-8">
          <Link to="/" className="font-bold capitalize text-primary">
            <Trans>Starter Plan</Trans>
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
