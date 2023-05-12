/* eslint-disable prettier/prettier */
import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

import type { Project } from '../types/Project';

export type UpdateProjectOwnerFormDTO = {
  current_owner: Project['owner'];
  new_owner: string;
};

export const FORM_INITIAL_VALUES = {
  current_owner: { name: '', id: '' },
  new_owner: '',
} satisfies UpdateProjectOwnerFormDTO;

export const FORM_VALIDATION_SCHEMA = object({
  new_owner: string().required(t`Please select a new owner.`),
});

export const FORM_OPTIONS = (onSubmit: any, initialValues: Project) => ({
  initialValues: {
    ...FORM_INITIAL_VALUES,
    current_owner: {
      name: initialValues.owner.name,
      id: initialValues.owner.id,
    },
  },
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useUpdateProjectOwnerForm = (onSubmit: any, initialValues: Project) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));
