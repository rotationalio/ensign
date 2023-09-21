/* eslint-disable prettier/prettier */
import { useFormik } from 'formik';
import { object } from 'yup';

import FormValidation from '@/lib/validation';
export type ChangePasswordFormDTO = {
  new_password: string;
  confirm_password: string;
};

export const FORM_INITIAL_VALUES = {
  new_password: '',
  confirm_password: '',
};

export const FORM_VALIDATION_SCHEMA = object({
  new_password: FormValidation.passwordValidation,
  confirm_password: FormValidation.ConfirmPassword('new_password'),
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
