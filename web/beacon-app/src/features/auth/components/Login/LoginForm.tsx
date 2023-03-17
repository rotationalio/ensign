import { Label } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider } from 'formik';
import { useState } from 'react';

import { CloseEyeIcon } from '@/components/icons/closeEyeIcon';
import { OpenEyeIcon } from '@/components/icons/openEyeIcon';
import Button from '@/components/ui/Button';
import TextField from '@/components/ui/TextField';

import { useLoginForm } from '../../types/LoginFormService';
import { AuthUser } from '../../types/LoginService';

type LoginFormProps = {
  onSubmit: (values: AuthUser, helpers: FormikHelpers<AuthUser>) => void;
  isDisabled?: boolean;
};

function LoginForm({ onSubmit, isDisabled }: LoginFormProps) {
  const formik = useLoginForm(onSubmit);

  const { touched, errors, getFieldProps } = formik;

  const [openEyeIcon, setOpenEyeIcon] = useState(false);

  const toggleEyeIcon = () => {
    setOpenEyeIcon(!openEyeIcon);
  };

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="mb-4 space-y-2">
          <fieldset>
            <Label htmlFor="email">Email</Label>
            <TextField
              placeholder="holly@golight.ly"
              fullWidth
              className="border-none pb-2"
              data-testid="email"
              errorMessage={touched.email && errors.email}
              {...getFieldProps('email')}
            />
          </fieldset>
          <fieldset>
            <Label htmlFor="Password">Password</Label>
            <div className="relative">
              <TextField
                placeholder={`Password (required)`}
                type={!openEyeIcon ? 'password' : 'text'}
                className="border-none"
                data-testid="password"
                errorMessage={touched.password && errors.password}
                fullWidth
                {...getFieldProps('password')}
              />
              <button
                type="button"
                onClick={toggleEyeIcon}
                className="absolute right-2 top-3 h-8 pb-2"
                data-testid="togglePassword"
              >
                {openEyeIcon ? <OpenEyeIcon /> : <CloseEyeIcon />}
                <span className="sr-only" data-testid="screenReadText">
                  {openEyeIcon ? 'Hide Password' : 'Show Password'}
                </span>
              </button>
            </div>
          </fieldset>
        </div>
        <div className="my-3 flex justify-between">
          <div id="google-recaptcha" className="flex flex-col"></div>
          <Button
            data-testid="login-button"
            type="submit"
            size="large"
            variant="secondary"
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
