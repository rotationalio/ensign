import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

export const FORM_INITIAL_VALUES = {
  name: '',
} as any;

export const FORM_VALIDATION_SCHEMA = object({
  name: string()
    .trim()
    .required(t`Name is required.`),
});

export const FORM_OPTIONS = (onSubmit: any, initialValues: any) => ({
  initialValues: {
    ...FORM_INITIAL_VALUES,
    ...initialValues,
  },
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useNameForm = (onSubmit: any, initialValues: any) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));
