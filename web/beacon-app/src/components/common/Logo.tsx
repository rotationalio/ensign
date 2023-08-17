import { Link } from 'react-router-dom';

import { EXTERNAL_LINKS } from '@/application';
import RotationalLogo from '@/assets/images/rotational.svg';

function Logo() {
  return (
    <Link to={EXTERNAL_LINKS.ROTATIONAL} target="_blank" data-testid="logo">
      <div className="flex items-center space-x-2">
        <img src={RotationalLogo} alt="Rotational Labs" className="h-12 w-12" />
        <h1 className="text-2xl font-bold text-primary">Rotational Labs</h1>
      </div>
    </Link>
  );
}

export default Logo;
