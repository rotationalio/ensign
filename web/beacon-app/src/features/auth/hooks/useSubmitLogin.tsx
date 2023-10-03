import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { useEffect } from 'react';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';

import { APP_ROUTE } from '@/constants';
import { useOrgStore } from '@/store';
import { getCookie, removeCookie } from '@/utils/cookies';
import { decodeToken } from '@/utils/decodeToken';

import { useLogin } from '../hooks/useLogin';
type Props = {
  setData: any;
  onReset: any;
  onSetCurrentUserEmail: any;
  resendEmailHandler: any;
};

const useSubmitLogin = ({ setData, onReset, onSetCurrentUserEmail, resendEmailHandler }: Props) => {
  const { authenticate, authenticated, auth, error, isAuthenticating } = useLogin() as any;
  const Store = useOrgStore((state) => state) as any;
  const navigate = useNavigate();
  const hasUnverifiedEmailError = error && error.response.status === 403;
  const onSubmitHandler = (values: any) => {
    onReset();
    const payload = {
      email: values.email,
      password: values.password,
    } as any;
    if (getCookie('invitee_token')) {
      payload['invite_token'] = getCookie('invitee_token');
    }

    setData(values.email);

    authenticate(payload);
  };

  useEffect(() => {
    if (authenticated) {
      onSetCurrentUserEmail('');
      const token = decodeToken(auth?.access_token) as any;
      Store.setAuthUser(token, !!authenticated);
      removeCookie('invitee_token');
      navigate(APP_ROUTE.DASHBOARD);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [authenticated, navigate, auth?.access_token]);

  useEffect(() => {
    if (hasUnverifiedEmailError) {
      toast.error(
        <div className="flex flex-col gap-5">
          <p>
            <Trans>Please verify your email address and try again!</Trans>
          </p>
          <div>
            <Button size="small" className="max-w-40 " onClick={resendEmailHandler}>
              <Trans>Resend Email</Trans>
            </Button>
          </div>
        </div>
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [hasUnverifiedEmailError]);

  return { onSubmitHandler, hasUnverifiedEmailError, authenticate, isAuthenticating };
};

export default useSubmitLogin;
