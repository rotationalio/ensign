/* eslint-disable prettier/prettier */
import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

export type ChangePasswordFormDTO = {
  new_password: string;
  confirm_password: string;
};

export const FORM_INITIAL_VALUES = {
  new_password: '',
  confirm_password: '',
};

export const FORM_VALIDATION_SCHEMA = object({
  new_password: string().required(t`New password is required`),
  confirm_password: string()
    .required(t`Confirm password is required`)
    .test('passwords-match', t`Passwords must match`, function (value) {
      return this.parent.new_password === value;
    }),
});

export const FORM_OPTIONS = (onSubmit: any, initialValues: any) => ({
  initialValues: {
    ...FORM_INITIAL_VALUES,
    ...initialValues,
  },
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useChangePasswordForm = (onSubmit: any, initialValues: any) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));
