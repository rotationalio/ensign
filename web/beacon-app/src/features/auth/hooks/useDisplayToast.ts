import { t } from '@lingui/macro';
import { useEffect } from 'react';
import toast from 'react-hot-toast';

import { updateQueryStringValueWithoutNavigation } from '@/utils/misc';
const useDisplayToast = (param: any) => {
  // handle toast from successfull email verification
  useEffect(() => {
    if (param?.accountVerified && param?.accountVerified === '1') {
      const isVerified = localStorage.getItem('isEmailVerified');
      if (isVerified === 'true') {
        toast.success(
          t`Thank you for verifying your email address.
          Log in now to start using Ensign.`
        );
      }
    }
    return () => {
      localStorage.removeItem('isEmailVerified');
    };
  }, [param?.accountVerified]);

  // handle toast from successfull reset password

  useEffect(() => {
    if (param?.from && param?.from === 'reset-password') {
      toast.success(
        t`Your password has been reset successfully. Please log in with your new password.`
      );
    }
    // remove to query param to avoid toast from showing up again
    return () => {
      updateQueryStringValueWithoutNavigation('from', null);
    };
  }, [param?.from]);
};

export default useDisplayToast;
