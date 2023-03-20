import { useState } from 'react';

import { CloseEyeIcon } from '@/components/icons/closeEyeIcon';
import { OpenEyeIcon } from '@/components/icons/openEyeIcon';

import { PasswordFieldProps } from './PasswordField.type';

export const usePasswordField = (props: PasswordFieldProps): PasswordFieldProps => {
  const [showPassword, setShowPassword] = useState(false);

  const handleShowPassword = () => setShowPassword(!showPassword);

  return {
    label: 'Password',
    ...props,
    type: showPassword ? 'text' : 'password',
    rightIcon: (
      <button onClick={handleShowPassword} data-testid="button">
        {showPassword ? <OpenEyeIcon /> : <CloseEyeIcon />}
        <span className="sr-only" data-testid="screenReadText">
          {showPassword ? 'Hide Password' : 'Show Password'}
        </span>
      </button>
    ),
  };
};
