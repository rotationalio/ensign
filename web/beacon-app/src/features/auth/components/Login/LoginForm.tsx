import { Form, FormikHelpers, FormikProvider } from 'formik';

import Button from '@/components/ui/Button';
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
          <TextField type="hidden" {...getFieldProps('invite_token')} />
          <TextField
            placeholder="holly@golight.ly"
            fullWidth
            label="Email"
            data-testid="email"
            errorMessage={touched.email && errors.email}
            {...getFieldProps('email')}
          />
          <PasswordField
            placeholder={`Password`}
            label="Password (required)"
            data-testid="password"
            errorMessage={touched.password && errors.password}
            fullWidth
            {...getFieldProps('password')}
          />
        </div>
        <div className="my-3 flex justify-between">
          <div id="google-recaptcha" className="flex flex-col"></div>
          <Button
            data-testid="login-button"
            type="submit"
            size="large"
            variant="secondary"
            isLoading={isLoading}
            className="mt-2 min-w-[100px] py-2"
            isDisabled={isDisabled}
            aria-label="Log in"
          >
            Log in
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}

export default LoginForm;
