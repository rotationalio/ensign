/* eslint-disable prettier/prettier */
import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

import type { Project } from '../types/Project';
import { UpdateProjectDTO } from './updateProjectService';

export type UpdateProjectFormDTO = {
  current_project: string;
} & Omit<UpdateProjectDTO['projectPayload'], 'id' | 'description'>;

export const FORM_INITIAL_VALUES = {
  name: '',
} satisfies Omit<UpdateProjectDTO['projectPayload'], 'id' | 'description'> & {
  project?: string;
};

export const FORM_VALIDATION_SCHEMA = object({
  name: string()
    .trim()
    .required(t`Project name is required.`)
    .max(512, t`Project name must be less than 512 characters.`),
});
export const FORM_OPTIONS = (onSubmit: any, initialValues: Partial<Project>) => ({
  initialValues: { ...FORM_INITIAL_VALUES, project: initialValues?.name },
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useUpdateProjectForm = (onSubmit: any, initialValues: Partial<Project>) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));
