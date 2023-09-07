/* eslint-disable prettier/prettier */
import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { array, object, string } from 'yup';

export const FORM_INITIAL_VALUES = {
  developer_segment: null,
  profession_segment: '',
} as any;

export const FORM_VALIDATION_SCHEMA = object({
  developer_segment: array()
    .min(1, t`Please select at least one option.`)
    .nullable()
    .required(t`Please select at least one option.`),

  profession_segment: string().required(t`Please select one option.`),
});
export const FORM_OPTIONS = (onSubmit: any, initialValues: any) => ({
  initialValues: {
    ...FORM_INITIAL_VALUES,
    ...initialValues,
  },
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const usePreferenceForm = (onSubmit: any, initialValues: any) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));
