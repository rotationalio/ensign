/* eslint-disable prettier/prettier */
import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

import { NewProjectDTO } from './createProjectService';

export const FORM_INITIAL_VALUES = {
  name: '',
  description: '',
} satisfies Omit<NewProjectDTO, 'tenantID'>;

export const FORM_VALIDATION_SCHEMA = object({
  name: string()
    .trim()
    .required(t`Project name is required`),
  description: string().notRequired(),
});
export const FORM_OPTIONS = (onSubmit: any) => ({
  initialValues: FORM_INITIAL_VALUES,
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useNewProjectForm = (onSubmit: any) => useFormik(FORM_OPTIONS(onSubmit));
