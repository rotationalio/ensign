import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';
import { useEffect } from 'react';

import StyledTextField from '@/components/ui/TextField/TextField';

import { useForgotPasswordForm } from './hooks/useForgotPasswordForm';

type ForgotPasswordFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  isLoading?: boolean;
  isSubmitted?: boolean;
};

const ForgotPasswordForm = ({ onSubmit, isLoading, isSubmitted }: ForgotPasswordFormProps) => {
  const formik = useForgotPasswordForm(onSubmit);
  const { getFieldProps, resetForm } = formik;

  useEffect(() => {
    if (isSubmitted) {
      resetForm();
    }
  }, [isSubmitted, resetForm]);

  return (
    <FormikProvider value={formik}>
      <Form>
        <StyledTextField
          fullWidth
          placeholder="Email address"
          label="Email address"
          labelClassName="sr-only"
          className="mb-2"
          {...getFieldProps('email')}
        />
        <ErrorMessage
          name="email"
          component={'p'}
          className="text-error-900 text-xs text-danger-700"
        />
        <div className="mt-3 flex justify-between">
          <div></div>
          <Button
            type="submit"
            variant="secondary"
            isLoading={isLoading}
            className="mt-2"
            data-cy="submit-forgot-password"
          >
            <Trans>Submit</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
};

export default ForgotPasswordForm;
