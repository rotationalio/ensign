import { useFormik } from 'formik';
import { object, ref, string } from 'yup';

export const FORM_INITIAL_VALUES = {
  password: '',
  pwcheck: '',
};

export const FORM_VALIDATION_SCHEMA = object({
  password: string()
    .required('Password is required.')
    .matches(/^(?=.*[a-z])/, 'The password must contain at least one lowercase letter.')
    .matches(/^(?=.*[A-Z])/, 'The password must contain at least one uppercase letter.')
    .matches(/^(?=.*[0-9])/, 'The password must contain at least one number.')
    .matches(
      /^(?=.*[!/[@#$%^&*+,-./:;<=>?^_`{|}~])/,
      'The password must contain at least one special character.'
    )
    .matches(/^(?=.{12,})/, 'The password must be at least 12 characters long.'),

  pwcheck: string()
    .oneOf([ref('password')], 'The paasswords must match.')
    .required('Please re-enter your password to confirm.'),
  reset_token: string().notRequired(),
});

export const FORM_OPTIONS = (onSubmit: any) => ({
  initialValues: FORM_INITIAL_VALUES,
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const usePasswordResetForm = (onSubmit: any) => useFormik(FORM_OPTIONS(onSubmit));
