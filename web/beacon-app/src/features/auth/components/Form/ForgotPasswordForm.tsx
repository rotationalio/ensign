import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';

// import { useEffect } from 'react';
import StyledTextField from '@/components/ui/TextField/TextField';

import { useForgotPasswordForm } from '../../hooks/useForgotPasswordForm';

type ForgotPasswordFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  isSubmitting?: boolean;
};

const ForgotPasswordForm = ({ onSubmit, isSubmitting }: ForgotPasswordFormProps) => {
  const formik = useForgotPasswordForm(onSubmit);
  const { getFieldProps } = formik;

  // useEffect(() => {
  //   if (isSubmitted) {
  //     resetForm();
  //   }
  // }, [isSubmitted, resetForm]);

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
          data-cy="forgot-password-email-input"
        />
        <ErrorMessage
          name="email"
          component={'p'}
          className="text-error-900 text-xs text-danger-700"
          data-cy="forgot-password-email-error"
        />
        <div className="mt-3 flex justify-between">
          <div></div>
          <Button
            type="submit"
            variant="secondary"
            disabled={isSubmitting}
            isLoading={isSubmitting}
            className="mt-2"
            data-cy="forgot-password-submit-bttn"
          >
            <Trans>Submit</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
};

export default ForgotPasswordForm;
