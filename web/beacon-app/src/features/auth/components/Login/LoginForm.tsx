import { t, Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider } from 'formik';
import { Link } from 'react-router-dom';

import { ROUTES } from '@/application';
import PasswordField from '@/components/ui/PasswordField/PasswordField';
import TextField from '@/components/ui/TextField';

import { useLoginForm } from '../../types/LoginFormService';
import { AuthUser, InviteAuthUser } from '../../types/LoginService';

type LoginFormProps = {
  onSubmit: (
    values: AuthUser | InviteAuthUser,
    helpers: FormikHelpers<AuthUser | InviteAuthUser>
  ) => void;
  isDisabled?: boolean;
  isLoading?: boolean;
  initialValues?: AuthUser | InviteAuthUser;
};

function LoginForm({ onSubmit, isDisabled, isLoading, initialValues }: LoginFormProps) {
  const formik = useLoginForm(onSubmit, initialValues);

  const { touched, errors, getFieldProps } = formik;

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="mb-4 space-y-2">
          <TextField
            placeholder="holly@golight.ly"
            fullWidth
            label={t`Email`}
            data-testid="email"
            errorMessage={touched.email && errors.email}
            {...getFieldProps('email')}
          />
          <PasswordField
            placeholder={t`Password`}
            label={t`Password (required)`}
            data-testid="password"
            errorMessage={touched.password && errors.password}
            fullWidth
            {...getFieldProps('password')}
          />
        </div>
        <div className="my-3 flex justify-between">
          <Link to={ROUTES.FORGOT_PASSWORD} className="mt-3 text-[#1D65A6] underline">
            <Trans>Forgot password?</Trans>
          </Link>
          <Button
            data-testid="login-button"
            type="submit"
            size="medium"
            variant="secondary"
            isLoading={isLoading}
            className="mt-2"
            disabled={isDisabled}
            aria-label={t`Log in`}
          >
            <Trans>Log in</Trans>
          </Button>
        </div>
        <div id="google-recaptcha" className="flex flex-col"></div>
      </Form>
    </FormikProvider>
  );
}

export default LoginForm;
