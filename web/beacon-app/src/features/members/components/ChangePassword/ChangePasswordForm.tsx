import { t, Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';

import PasswordField from '@/components/ui/PasswordField';

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
        <PasswordField
          labelClassName="font-bold"
          label={t`Enter new password`}
          fullWidth
          {...getFieldProps('new_password')}
        />
        <ErrorMessage name="new_password" component="small" className="text-xs text-danger-500" />
        <PasswordField
          label={t`Confirm new password`}
          labelClassName="font-bold"
          fullWidth
          {...getFieldProps('confirm_password')}
        />
        <ErrorMessage
          name="confirm_password"
          component="small"
          className="text-xs text-danger-500"
        />
        <div className="py-5 text-center">
          <Button type="submit" isLoading={isSubmitting} disabled={isSubmitting}>
            <Trans>Save</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}

export default ChangePasswordForm;
