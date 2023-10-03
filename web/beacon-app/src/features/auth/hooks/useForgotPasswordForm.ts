import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

export const FORM_INITIAL_VALUES = {
  email: '',
} as any;

export const FORM_VALIDATION_SCHEMA = object({
  email: string()
    .trim()
    .email(t`Please enter a valid email address.`)
    .required(t`Email is required.`),
});

export const FORM_OPTIONS = (onSubmit: any) => ({
  initialValues: FORM_INITIAL_VALUES,
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useForgotPasswordForm = (onSubmit: any) => useFormik(FORM_OPTIONS(onSubmit));
