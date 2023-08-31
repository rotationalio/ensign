/* eslint-disable prettier/prettier */
import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

export const FORM_INITIAL_VALUES = {
  developer_segment: {
    value: '',
    label: '',
  },
} as any;

export const FORM_VALIDATION_SCHEMA = object({
  developer_segment: string()
    .required(t`Please select at least one option`)
    .test('developer_segment', t`Please select at least one option`, function (value: any) {
      const { developer_segment } = this.parent;
      return value?.length > 0 || developer_segment?.length > 0;
    }),
});
export const FORM_OPTIONS = (onSubmit: any, initialValues: any) => ({
  initialValues: {
    ...FORM_INITIAL_VALUES,
    developer_segment: initialValues?.developer_segment,
  },
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const usePreferenceForm = (onSubmit: any, initialValues: any) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));
