import { AriaButton as Button } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import { Logo } from '../ui';

function LandingHeader() {
  return (
    <nav className="flex flex-wrap items-center justify-between border-b border-primary-800 pb-4">
      <Logo />
      <div className="space-x-8">
        <Link to="/" className="font-bold capitalize text-primary">
          Starter Plan
        </Link>
        <Link to="/">
          <Button variant="tertiary" size="small">
            Upgrade
          </Button>
        </Link>
      </div>
    </nav>
  );
}

export default LandingHeader;
