import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';

import { useChangePasswordForm } from '../../types/changePasswordFormService';

type ChangePasswordForm = {
  handleSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  initialValues: any;
};

function ChangePasswordForm({ handleSubmit, initialValues }: ChangePasswordForm) {
  const formik = useChangePasswordForm(handleSubmit, initialValues);

  const { getFieldProps, isSubmitting } = formik;

  return (
    <FormikProvider value={formik}>
      <Form className="space-y-3">
        <TextField label={t`Enter new password`} fullWidth {...getFieldProps('new_password')} />
        <ErrorMessage name="new_password" component="small" className="text-xs text-danger-500" />
        <TextField
          label={t`confirm new password`}
          fullWidth
          {...getFieldProps('confirm_password')}
          iv
          className="pt-3 text-center"
        />
        <ErrorMessage
          name="confirm_password"
          component="small"
          className="text-xs text-danger-500"
        />
        <Button type="submit" isLoading={isSubmitting} disabled={isSubmitting}>
          <Trans>Save</Trans>
        </Button>
      </Form>
    </FormikProvider>
  );
}

export default ChangePasswordForm;
