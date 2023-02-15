import { TextField } from '@rotational/beacon-core';
import { useState } from 'react';

import { HidePassword } from '@/components/icons/hidePassword';
import { ShowPassword } from '@/components/icons/showPassword';

export default function Password() {
  const [showEye, setShowEye] = useState(false);

  const toggleEyeIcon = () => {
    setShowEye(!showEye);
  };

  return (
    <div className="flex space-x-3">
      <TextField
        label="Password"
        placeholder="Password"
        type={!showEye ? 'password' : 'text'}
        fullWidth
        data-testid="password"
      />
      <button onClick={toggleEyeIcon} className="mt-4 self-center" data-testid="button">
        {showEye ? <ShowPassword /> : <HidePassword />}
        <span className="sr-only" data-testid="screenReadText">
          {showEye ? 'Hide Password' : 'Show Password'}
        </span>
      </button>
    </div>
  );
}
