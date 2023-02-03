import { Link } from 'react-router-dom';

import RotationalLogo from '@/assets/images/rotational.svg';

function Logo() {
  return (
    <Link to="/" data-testid="logo">
      <div className="flex items-center space-x-2">
        <img src={RotationalLogo} alt="Rotational Lab" className="h-12 w-12" />
        <h1 className="text-2xl font-bold text-primary">Rotational Lab</h1>
      </div>
    </Link>
  );
}

export default Logo;
