import { t } from '@lingui/macro';
import { useEffect } from 'react';
import toast from 'react-hot-toast';
const useVerfiedAccount = (param: any) => {
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
};

export default useVerfiedAccount;
