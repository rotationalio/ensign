/* eslint-disable prettier/prettier */
import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

export const FORM_INITIAL_VALUES = {
  workspace: '',
} as any;

export const FORM_VALIDATION_SCHEMA = object({
  workspace: string()
    .trim()
    .required(t`Workspace name is required`),
});
export const FORM_OPTIONS = (onSubmit: any) => ({
  initialValues: FORM_INITIAL_VALUES,
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useNewWorkspaceForm = (onSubmit: any) => useFormik(FORM_OPTIONS(onSubmit));
