import { t } from '@lingui/macro';
import { useEffect, useRef } from 'react';
import toast from 'react-hot-toast';

import { updateQueryStringValueWithoutNavigation } from '@/utils/misc';
const useDisplayToast = (param: any) => {
  // handle toast from successfull email verification
  const isVerifiedRef = useRef<boolean>(param?.accountVerified && param?.accountVerified === '1');
  const isResetRef = useRef<boolean>(param?.from && param?.from === 'reset-password');

  useEffect(() => {
    if (isVerifiedRef.current) {
      toast.success(
        t`Thank you for verifying your email address.
          Log in now to start using Ensign.`
      );
    }
    return () => {
      updateQueryStringValueWithoutNavigation('from', null);
      isVerifiedRef.current = false;
    };
  }, [isVerifiedRef]);

  // handle toast from successfull reset password

  useEffect(() => {
    if (isResetRef.current) {
      toast.success(
        t`Your password has been reset successfully. Please log in with your new password.`
      );
    }
    // remove to query param to avoid toast from showing up again
    return () => {
      updateQueryStringValueWithoutNavigation('from', null);
      isResetRef.current = false;
    };
  }, [isResetRef]);
};

export default useDisplayToast;
