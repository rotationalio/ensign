import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

export const FORM_INITIAL_VALUES = {
  organization: '',
} as any;

export const FORM_VALIDATION_SCHEMA = object({
  organization: string()
    .trim()
    .required(t`Team or organization name is required.`),
});

export const FORM_OPTIONS = (onSubmit: any) => ({
  initialValues: FORM_INITIAL_VALUES,
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useOrganizationForm = (onSubmit: any) => useFormik(FORM_OPTIONS(onSubmit));
