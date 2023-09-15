import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';

import StyledTextField from '@/components/ui/TextField/TextField';

import { useForgotPasswordForm } from './hooks/useForgotPasswordForm';

type ForgotPasswordFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
};

const ForgotPasswordForm = ({ onSubmit }: ForgotPasswordFormProps) => {
  const formik = useForgotPasswordForm(onSubmit);
  const { getFieldProps } = formik;
  return (
    <FormikProvider value={formik}>
      <Form>
        <StyledTextField
          fullWidth
          placeholder="Email address"
          label="Email address"
          labelClassName="sr-only"
          className="mb-4"
          {...getFieldProps('email')}
        />
        <ErrorMessage
          name="email"
          component={'p'}
          className="text-error-900 py-2 text-xs text-danger-700"
        />

        <Button type="submit" variant="secondary" data-cy="submit-forgot-password">
          <Trans>Submit</Trans>
        </Button>
      </Form>
    </FormikProvider>
  );
};

export default ForgotPasswordForm;
